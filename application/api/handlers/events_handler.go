package handlers

import (
	"net/http"

	"eventflow/infra/contracts"
	"eventflow/infra/dto"

	"github.com/labstack/echo/v4"
)

// EventsHandler will handle the events endpoint.
type EventsHandler struct {
	store contracts.EventStore
}

// NewEventsHandler will create a new EventsHandler
func NewEventsHandler(store contracts.EventStore) *EventsHandler {
	return &EventsHandler{store: store}
}

// POST /events.
func (h *EventsHandler) Create(c echo.Context) error {
	var req dto.CreateEventRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := validateEventRequest(&req); err != nil {
		return err
	}

	event := req.ToDomain()
	ctx := c.Request().Context()
	exists, err := h.store.Exists(ctx, event.IdempotencyKey())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "duplicate event"})
	}

	if err := h.store.Insert(ctx, event); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}

	return c.JSON(http.StatusAccepted, dto.CreateEventResponse{Status: "accepted"})
}

// validateEventRequest validates the DTO request.
func validateEventRequest(req *dto.CreateEventRequest) error {
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
