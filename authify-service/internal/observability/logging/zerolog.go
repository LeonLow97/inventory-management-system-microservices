package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/LeonLow97/internal/core/ports"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
)

type zerologLogger struct {
	logger zerolog.Logger
	ctx    context.Context
}

// New creates a timestamped logger with colored console output on terminals
// and JSON output for other writers.
func New(output io.Writer) ports.Logger {
	if output == nil {
		output = io.Discard
	}

	logger := zerolog.New(outputWriter(output)).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Logger()

	return &zerologLogger{logger: logger}
}

func (l *zerologLogger) WithContext(ctx context.Context) ports.Logger {
	if ctx == nil {
		return l
	}
	return &zerologLogger{logger: l.logger, ctx: ctx}
}

func (l *zerologLogger) WithError(err error) ports.Logger {
	if err == nil {
		return l
	}
	return &zerologLogger{
		logger: l.logger.With().Err(err).Logger(),
		ctx:    l.ctx,
	}
}

func (l *zerologLogger) Debug(message string, fields ...ports.LogField) {
	write(l.logger.Debug(), l.ctx, message, fields)
}

func (l *zerologLogger) Info(message string, fields ...ports.LogField) {
	write(l.logger.Info(), l.ctx, message, fields)
}

func (l *zerologLogger) Warn(message string, fields ...ports.LogField) {
	write(l.logger.Warn(), l.ctx, message, fields)
}

func (l *zerologLogger) Error(message string, fields ...ports.LogField) {
	write(l.logger.Error(), l.ctx, message, fields)
}

func write(event *zerolog.Event, ctx context.Context, message string, fields []ports.LogField) {
	if ctx != nil {
		event.Ctx(ctx)
	}
	for _, field := range fields {
		switch value := field.Value.(type) {
		case error:
			event.AnErr(field.Key, value)
		case string:
			event.Str(field.Key, value)
		case bool:
			event.Bool(field.Key, value)
		case int:
			event.Int(field.Key, value)
		case int64:
			event.Int64(field.Key, value)
		case float64:
			event.Float64(field.Key, value)
		default:
			event.Interface(field.Key, value)
		}
	}
	event.Msg(message)
}

func outputWriter(output io.Writer) io.Writer {
	fileDescriptor, ok := output.(interface{ Fd() uintptr })
	if !ok || !isatty.IsTerminal(fileDescriptor.Fd()) {
		return output
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        output,
		TimeFormat: time.RFC3339,
	}
	consoleWriter.FormatLevel = func(level interface{}) string {
		return formatLevel(level)
	}
	return consoleWriter
}

func formatLevel(level interface{}) string {
	name := strings.ToLower(fmt.Sprint(level))
	label := strings.ToUpper(name)
	colorCode := ""

	switch name {
	case "trace", "debug":
		colorCode = "36"
	case "info":
		colorCode = "94"
	case "warn":
		colorCode = "33"
	case "error":
		colorCode = "31"
	case "fatal", "panic":
		colorCode = "1;31"
	}

	label = fmt.Sprintf("%-5s", label)
	if colorCode == "" || os.Getenv("NO_COLOR") != "" {
		return label
	}
	return fmt.Sprintf("\x1b[%sm%s\x1b[0m", colorCode, label)
}
