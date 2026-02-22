package api

import (
	"eventflow/application/api/handlers"
	"eventflow/application/api/middleware"
	"eventflow/application/events"
	"eventflow/application/metrics"

	"github.com/labstack/echo/v4"
)

// Router sets up Echo routes.
func Router(
	e *echo.Echo,
	eventService *events.Service,
	metricsService *metrics.Service,
) {
	e.Use(middleware.RequestLogger())

	eventsHandler := handlers.NewEventsHandler(eventService)
	metricsHandler := handlers.NewMetricsHandler(metricsService)

	e.POST("/events", eventsHandler.Create)
	e.GET("/metrics", metricsHandler.Get)
}
