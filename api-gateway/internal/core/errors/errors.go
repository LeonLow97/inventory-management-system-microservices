package coreerrors

type Kind string

const (
	KindInvalidArgument    Kind = "invalid_argument"
	KindUnauthenticated    Kind = "unauthenticated"
	KindPermissionDenied   Kind = "permission_denied"
	KindNotFound           Kind = "not_found"
	KindAlreadyExists      Kind = "already_exists"
	KindFailedPrecondition Kind = "failed_precondition"
	KindOutOfRange         Kind = "out_of_range"
	KindUnavailable        Kind = "unavailable"
	KindDeadlineExceeded   Kind = "deadline_exceeded"
	KindResourceExhausted  Kind = "resource_exhausted"
	KindConflict           Kind = "conflict"
	KindInternal           Kind = "internal"
)

type Error struct {
	Kind    Kind
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return string(e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func New(kind Kind, message string, cause error) *Error {
	return &Error{Kind: kind, Message: message, Cause: cause}
}
