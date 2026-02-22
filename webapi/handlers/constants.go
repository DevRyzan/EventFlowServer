package handlers

const (
	// errors
	InvalidRequestBodyError = "invalid request body"
	ValidationFailedError   = "validation failed"
	InternalErrorError      = "internal error"
	DuplicateRejectedError  = "duplicate rejected"
	InternalErrorWarn       = "internal error"

	// warn
	InvalidRequestBodyWarn = "invalid request body"
	ValidationFailedWarn   = "validation failed"
	DuplicateRejectedWarn  = "duplicate rejected"

	// msg
	InvalidRequestBodyMsg = "invalid request body"
	EventNameRequiredMsg  = "event_name is required"
	UserIdRequiredMsg     = "user_id is required"
	TimestampRequiredMsg  = "timestamp is required"
	TimestampInvalidMsg   = "timestamp must be a positive Unix timestamp"

	FromRequiredMsg  = "from is required (Unix timestamps)"
	ToRequiredMsg    = "to is required (Unix timestamps)"
	FromToInvalidMsg = "from must be less than or equal to to"
	SuccessMsg       = "success"
)
