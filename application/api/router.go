package api

import (
	"eventflow/application/api/handlers"
	"eventflow/infra/contracts"
	"eventflow/infra/domain"

	"github.com/labstack/echo/v4"
)

// MetricsService providing the metrics endpoint.
type MetricsService interface {
	Get(
		ctx echo.Context, 
		req *domain.MetricsRequest,
	) (*domain.MetricsResponse, error)
}

// Router sets up Echo routes.
func Router(
	e *echo.Echo, 
	eventStore contracts.EventStore, 
	metricsService MetricsService,
) {
	eventsHandler := handlers.NewEventsHandler(eventStore)
	metricsHandler := handlers.NewMetricsHandler(metricsService)

	e.POST("/events", eventsHandler.Create)
	e.GET("/metrics", metricsHandler.Get)
}
