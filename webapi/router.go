package webapi

import (
	"eventflow/application/commands/create_event"
	"eventflow/application/queries/get_metrics"
	"eventflow/middleware"
	"eventflow/webapi/handlers"

	"github.com/labstack/echo/v4"
)

// Router sets up Echo routes.
func Router(
	e *echo.Echo,
	createEventService *create_event.Service,
	getMetricsService *get_metrics.Service,
) {
	e.Use(middleware.RequestLogger())

	eventsHandler := handlers.NewEventsHandler(createEventService)
	metricsHandler := handlers.NewMetricsHandler(getMetricsService)

	e.POST("/events", eventsHandler.Handle)
	e.GET("/metrics", metricsHandler.Handle)
}
