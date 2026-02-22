package handlers

import (
	"log/slog"
	"net/http"

	"eventflow/application/queries/get_metrics"

	"github.com/labstack/echo/v4"
)

// MetricsHandler handles GET /metrics.
type MetricsHandler struct {
	service *get_metrics.Service
}

// NewMetricsHandler creates a new MetricsHandler.
func NewMetricsHandler(service *get_metrics.Service) *MetricsHandler {
	return &MetricsHandler{service: service}
}

// Handle processes the get metrics request.
func (h *MetricsHandler) Handle(c echo.Context) error {
	var req get_metrics.GetMetricsRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn(
			InvalidRequestBodyWarn,
			"err", err,
		)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": InvalidRequestBodyMsg})
	}

	if err := validateMetricsRequest(&req); err != nil {
		slog.Warn(
			ValidationFailedWarn,
			"event_name", req.EventName,
			"err", err,
		)
		return err
	}

	resp, err := h.service.Handle(c.Request().Context(), &req)
	if err != nil {
		slog.Error(
			InternalErrorWarn,
			"event_name", req.EventName,
			"err", err,
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": InternalErrorError})
	}

	slog.Info(
		SuccessMsg,
		"event_name", req.EventName,
		"total_count", resp.TotalCount,
	)
	return c.JSON(http.StatusOK, resp)
}

func validateMetricsRequest(req *get_metrics.GetMetricsRequest) error {
	if req.EventName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, EventNameRequiredMsg)
	}
	if req.From <= 0 || req.To <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, FromRequiredMsg)
	}
	if req.From > req.To {
		return echo.NewHTTPError(http.StatusBadRequest, FromToInvalidMsg)
	}
	return nil
}
