package handlers

import (
	"log/slog"
	"net/http"

	commands "eventflow/application/dtos/commands"
	"eventflow/application/events"

	"github.com/labstack/echo/v4"
)

// EventsHandler handles the events endpoint.
type EventsHandler struct {
	service *events.Service
}

// NewEventsHandler creates a new EventsHandler.
func NewEventsHandler(service *events.Service) *EventsHandler {
	return &EventsHandler{service: service}
}

// Create handles POST /events.
func (h *EventsHandler) Create(c echo.Context) error {
	var req commands.CreateEventRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn("event create: invalid request body", "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := validateEventRequest(&req); err != nil {
		slog.Warn("event create: validation failed", "event_name", req.EventName, "user_id", req.UserId, "err", err)
		return err
	}

	ctx := c.Request().Context()
	resp, err := h.service.Create(ctx, &req)
	if err != nil {
		if err == events.ErrDuplicate {
			slog.Info("event create: duplicate rejected", "event_name", req.EventName, "user_id", req.UserId)
			return c.JSON(http.StatusConflict, map[string]string{"error": "duplicate event"})
		}
		slog.Error("event create: internal error", "event_name", req.EventName, "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	slog.Info("event create: accepted", "event_name", req.EventName, "user_id", req.UserId)
	return c.JSON(http.StatusAccepted, resp)
}

func validateEventRequest(req *commands.CreateEventRequest) error {
	if req.EventName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "event_name is required")
	}
	if req.UserId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id is required")
	}
	if req.Timestamp <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "timestamp must be a positive Unix timestamp")
	}
	return nil
}
