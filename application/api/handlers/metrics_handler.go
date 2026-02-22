package handlers

import (
	"net/http" 
	"eventflow/infra/domain" 
	"github.com/labstack/echo/v4"
)

// MetricsHandler will hadnle the metrics endpoint
type MetricsHandler struct {
	service interface {
		Get(
			ctx echo.Context, 
			req *domain.MetricsRequest,
		) (*domain.MetricsResponse, error)
	}
}

// NewMetricsHandler will create a new MetricsHandler
func NewMetricsHandler(service interface {
	Get(
		ctx echo.Context, 
		req *domain.MetricsRequest,
	) (*domain.MetricsResponse, error)
}) *MetricsHandler {
	return &MetricsHandler{service: service}
}

// GET /metrics.
func (h *MetricsHandler) Get(
	c echo.Context,
) (err error) {
	var req domain.MetricsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid query params"})
	}

	if err := validateMetricsRequest(&req); err != nil {
		return err
	}

	resp, err := h.service.Get(c, &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	return c.JSON(http.StatusOK, resp)
}

// validation for the metrics request
func validateMetricsRequest(
	req *domain.MetricsRequest,
) error {
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
