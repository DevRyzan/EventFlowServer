package load_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"sort"
	"sync"
	"testing"
	"time"
)

const (
	defaultBaseURL  = "http://localhost:8080"
	defaultWorkers  = 10
	defaultDuration = 10 * time.Second
)

// TestLoadEvents tests the load test for the events endpoint
func TestLoadEvents(t *testing.T) {
	if testing.Short() {
		slog.Info("skipping load test in short mode")
		t.Skip("skipping load test in short mode")
	}
	runLoadTest(t, "events", defaultWorkers, defaultDuration)
}

// TestLoadMetrics tests the load test for the metrics endpoint
func TestLoadMetrics(t *testing.T) {
	if testing.Short() {
		slog.Info("skipping load test in short mode")
		t.Skip("skipping load test in short mode")
	}
	runLoadTest(t, "metrics", defaultWorkers, defaultDuration)
}

func TestLoadBoth(t *testing.T) {
	if testing.Short() {
		slog.Info("skipping load test in short mode")
		t.Skip("skipping load test in short mode")
	}
	runLoadTest(t, "both", defaultWorkers, defaultDuration)
}

// runLoadTest runs the load test for the given mode, workers and duration
func runLoadTest(t *testing.T, mode string, workers int, duration time.Duration) {
	baseURL := os.Getenv("LOADTEST_URL")
	if baseURL == "" {
		slog.Info("base URL is not set, using default base URL")
		baseURL = defaultBaseURL
	}

	t.Logf("Load test: %s | workers=%d | duration=%s | mode=%s", baseURL, workers, duration, mode)

	var wg sync.WaitGroup
	results := make(chan result, 100000)
	stop := time.After(duration)
	client := &http.Client{Timeout: 5 * time.Second}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runWorker(id, baseURL, mode, client, results, stop)
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var latencies []float64
	var success, failed int64
	for r := range results {
		if r.err != nil {
			failed++
			continue
		}
		success++
		latencies = append(latencies, r.durationMs)
	}

	total := success + failed
	durSec := duration.Seconds()
	throughput := float64(success) / durSec

	t.Logf("--- Results--")
	t.Logf(
		"Total: %d | Success: %d | Failed: %d | Throughput: %.1f req/s",
		total, success, failed, throughput,
	)

	if len(latencies) > 0 {
		sort.Float64s(latencies)
		avg := mean(latencies)
		p50 := percentile(latencies, 50)
		p95 := percentile(latencies, 95)
		p99 := percentile(latencies, 99)
		t.Logf("--- Latency (ms)   Avg: %.1f | P50: %.1f | P95: %.1f | P99: %.1f", avg, p50, p95, p99)
	}

	if failed > 0 && float64(failed)/float64(total) > 0.1 {
		slog.Error("failure rate too high", "failed", failed, "total", total, "rate", 100*float64(failed)/float64(total))
		t.Errorf("failure rate too high: %d/%d (%.1f%%)", failed, total, 100*float64(failed)/float64(total))
	}
}

type result struct {
	durationMs float64
	err        error
}

func runWorker(id int, baseURL, mode string, client *http.Client, results chan<- result, stop <-chan time.Time) {
	timestamp := time.Now().Unix()

	for {
		select {
		case <-stop:
			return
		default:
			start := time.Now()
			var err error

			switch mode {
			case "events":
				err = postEvent(client, baseURL, id, timestamp)
			case "metrics":
				err = getMetrics(client, baseURL)
			case "both":
				if e := postEvent(client, baseURL, id, timestamp); e != nil {
					err = e
				}
				if e := getMetrics(client, baseURL); e != nil && err == nil {
					err = e
				}
			}

			results <- result{
				durationMs: float64(time.Since(start).Milliseconds()),
				err:        err,
			}
			timestamp++
		}
	}
}

func postEvent(client *http.Client, baseURL string, workerID int, ts int64) error {
	body := map[string]interface{}{
		"user_id":    fmt.Sprintf("user-%d", workerID),
		"event_name": "click",
		"timestamp":  ts,
		"channel":    "web",
	}
	b, _ := json.Marshal(body)
	resp, err := client.Post(baseURL+"/events", "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusConflict {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func getMetrics(client *http.Client, baseURL string) error {
	from := time.Now().Add(-24 * time.Hour).Unix()
	to := time.Now().Unix()
	url := fmt.Sprintf("%s/metrics?event_name=click&from=%d&to=%d", baseURL, from, to)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func percentile(sorted []float64, p int) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil(float64(len(sorted)*p)/100)) - 1
	if idx < 0 {
		idx = 0
	}
	return sorted[idx]
}

func mean(v []float64) float64 {
	if len(v) == 0 {
		return 0
	}
	var sum float64
	for _, x := range v {
		sum += x
	}
	return sum / float64(len(v))
}
