package ports

import "context"

// LogField is a structured field attached to a log message.
type LogField struct {
	Key   string
	Value any
}

// Logger is the application's package-neutral logging contract.
type Logger interface {
	// WithContext returns a logger that makes ctx available to log-event hooks.
	// Context values are not serialized automatically.
	WithContext(ctx context.Context) Logger
	// WithError returns a logger that includes err as the structured "error" field.
	WithError(err error) Logger
	// Debug records a diagnostic message.
	Debug(message string, fields ...LogField)
	// Info records an informational message.
	Info(message string, fields ...LogField)
	// Warn records a recoverable problem.
	Warn(message string, fields ...LogField)
	// Error records an operation failure.
	Error(message string, fields ...LogField)
}
