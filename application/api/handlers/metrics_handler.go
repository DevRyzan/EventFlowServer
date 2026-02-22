package handlers

import (
	"log/slog"
	"net/http"

	"eventflow/application/metrics"
	"eventflow/infra/domain"

	"github.com/labstack/echo/v4"
)

// MetricsHandler handles the metrics endpoint.
type MetricsHandler struct {
	service *metrics.Service
}

// NewMetricsHandler creates a new MetricsHandler.
func NewMetricsHandler(service *metrics.Service) *MetricsHandler {
	return &MetricsHandler{service: service}
}

// Get handles GET /metrics.
func (h *MetricsHandler) Get(c echo.Context) error {
	var req domain.MetricsRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("metrics get: invalid query params", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid query params"})
	}

	if err := validateMetricsRequest(&req); err != nil {
		slog.Warn("metrics get: validation failed", "event_name", req.EventName, "err", err)
		return err
	}

	resp, err := h.service.Get(c.Request().Context(), &req)
	if err != nil {
		slog.Error("metrics get: internal error", "event_name", req.EventName, "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	slog.Info("metrics get: success", "event_name", req.EventName, "total_count", resp.TotalCount)
	return c.JSON(http.StatusOK, resp)
}

func validateMetricsRequest(req *domain.MetricsRequest) error {
	if req.EventName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "event_name is required")
	}
	if req.From <= 0 || req.To <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "from and to are required (Unix timestamps)")
	}
	if req.From > req.To {
		return echo.NewHTTPError(http.StatusBadRequest, "from must be less than or equal to to")
	}
	return nil
}
