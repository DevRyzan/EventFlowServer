package handlers

import (
	"log/slog"
	"net/http"

	"eventflow/application/commands/create_event"

	"github.com/labstack/echo/v4"
)

// EventsHandler handles POST /events.
type EventsHandler struct {
	service *create_event.Service
}

// NewEventsHandler creates a new EventsHandler.
func NewEventsHandler(service *create_event.Service) *EventsHandler {
	return &EventsHandler{service: service}
}

// Handle processes the create event request.
func (h *EventsHandler) Handle(c echo.Context) error {
	var req create_event.CreateEventRequest
	if err := c.Bind(&req); err != nil {
		slog.Warn(InvalidRequestBodyError, "err", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": InvalidRequestBodyMsg})
	}

	if err := validateEventRequest(&req); err != nil {
		slog.Warn(ValidationFailedError, "event_name", req.EventName, "user_id", req.UserId, "err", err)
		return err
	}

	ctx := c.Request().Context()
	resp, err := h.service.Execute(ctx, &req)
	if err != nil {
		if err == create_event.ErrDuplicate {
			slog.Info(
				DuplicateRejectedError,
				"event_name", req.EventName,
				"user_id", req.UserId,
			)
			return c.JSON(http.StatusConflict, map[string]string{"error": DuplicateRejectedError})
		}
		slog.Error(
			InternalErrorError,
			"event_name", req.EventName,
			"err", err,
		)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": InternalErrorError})
	}

	slog.Info(
		SuccessMsg,
		"event_name", req.EventName,
		"user_id", req.UserId,
	)
	return c.JSON(http.StatusAccepted, resp)
}

func validateEventRequest(req *create_event.CreateEventRequest) error {
	if req.EventName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, EventNameRequiredMsg)
	}
	if req.UserId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, UserIdRequiredMsg)
	}
	if req.Timestamp <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, TimestampInvalidMsg)
	}
	return nil
}
