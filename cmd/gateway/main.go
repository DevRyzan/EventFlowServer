// Gateway is a load balancer that round-robins requests across backend instances.
// Configure via BACKEND_URLS (comma-separated) and GATEWAY_PORT (default 8080).
package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"

	"eventflow/infra"
)

func main() {
	infra.LoadEnv(".env")
	cfg := infra.LoadGateway()

	urls := parseBackendURLs(cfg.BackendURLs)
	if len(urls) == 0 {
		slog.Error("no valid backend URLs")
		os.Exit(1)
	}

	lb := newRoundRobinBalancer(urls)
	proxy := &httputil.ReverseProxy{
		Director: lb.director,
	}

	http.Handle("/", proxy)
	http.HandleFunc("/gateway/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	slog.Info("gateway listening", "port", cfg.Port, "backends", len(urls))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
		slog.Error("gateway failed", "err", err)
		os.Exit(1)
	}
}

func parseBackendURLs(s string) []*url.URL {
	parts := strings.Split(s, ",")
	var urls []*url.URL
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		u, err := url.Parse(p)
		if err != nil {
			slog.Warn("invalid backend URL", "url", p, "err", err)
			continue
		}
		urls = append(urls, u)
	}
	return urls
}

type roundRobinBalancer struct {
	urls    []*url.URL
	counter uint64
}

func newRoundRobinBalancer(urls []*url.URL) *roundRobinBalancer {
	return &roundRobinBalancer{urls: urls}
}

func (lb *roundRobinBalancer) director(req *http.Request) {
	n := atomic.AddUint64(&lb.counter, 1)
	idx := (n - 1) % uint64(len(lb.urls))
	target := lb.urls[idx]

	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
	if target.RawQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
	}
	if _, ok := req.Header["User-Agent"]; !ok {
		req.Header.Set("User-Agent", "")
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
