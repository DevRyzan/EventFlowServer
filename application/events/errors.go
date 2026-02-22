package events

import "errors"

// ErrDuplicate is returned when an event with the same idempotency key already exists.
var ErrDuplicate = errors.New("duplicate event")
