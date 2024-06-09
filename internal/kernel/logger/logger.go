package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Notice(msg string, args ...any) {
	l.NoticeContext(context.TODO(), msg, args...)
}

func (l *Logger) NoticeContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, LevelNotice, msg, args...)
}

// Trace logs TBD
func (l *Logger) Trace(msg string, args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	args = append(args, slog.String("trace", fmt.Sprintf("%s:%d (%v)", file, line, details.Name())))

	l.Log(context.TODO(), LevelTrace, msg, args...)
}

// Trace logs TBD
func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	pc, file, line, _ := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	args = append(args, slog.String("trace", fmt.Sprintf("%s:%d (%v)", file, line, details.Name())))

	l.Log(ctx, LevelTrace, msg, args...)
}

// Builder object
type loggerBuilder struct {
	lvl        slog.Level
	withSource bool
	writers    []io.Writer
}

// NewBuilder return a new logger builder object
func NewBuilder() *loggerBuilder {
	return new(loggerBuilder)
}

// WithWriter sets a specific writer
func (b *loggerBuilder) WithWriter(w io.Writer) *loggerBuilder {
	b.writers = append(b.writers, w)

	return b
}

// WithLevel sets log level
func (b *loggerBuilder) WithLevel(l slog.Level) *loggerBuilder {
	b.lvl = l

	return b
}

// WithSource adds a cource file into an output record
func (b *loggerBuilder) WithSource() *loggerBuilder {
	b.withSource = true

	return b
}

// Build returns the logger
func (b *loggerBuilder) Build() *Logger {
	if len(b.writers) == 0 {
		b.writers = append(b.writers, os.Stdout)
	}

	w := io.MultiWriter(b.writers...)

	l := newLogger(b.lvl, w)

	return &Logger{l}
}

const (
	// Trace log level
	LevelTrace slog.Level = -100
	// Notice level
	LevelNotice slog.Level = 2
	// Emergency log level
	LevelEmergency slog.Level = 15
)

func newLogger(lvl slog.Level, w io.Writer) *slog.Logger {
	return slog.New(
		slog.NewJSONHandler(w, &slog.HandlerOptions{
			Level: lvl,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)

					switch {
					case level < slog.LevelDebug:
						a.Value = slog.StringValue("TRACE")
					case level < slog.LevelInfo:
						a.Value = slog.StringValue("DEBUG")
					case level < LevelNotice:
						a.Value = slog.StringValue("INFO")
					case level < slog.LevelWarn:
						a.Value = slog.StringValue("NOTICE")
					case level < slog.LevelError:
						a.Value = slog.StringValue("WARNING")
					case level < LevelEmergency:
						a.Value = slog.StringValue("ERROR")
					default:
						a.Value = slog.StringValue("EMERGENCY")
					}
				}

				return a

			},
		}),
	)
}

// Error logging attribute
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// Maps level from a string. By default returns slog.LevelInfo
func MapLevel(lvl string) slog.Level {
	switch lvl {
	case "trace":
		return LevelTrace
	case "dev", "local", "debug":
		return slog.LevelDebug
	case "notice":
		return LevelNotice
	case "warn":
		return slog.LevelWarn
	case "info":
		return slog.LevelInfo
	case "error":
		return slog.LevelError
	case "emergency":
		return LevelEmergency
	default:
		return slog.LevelInfo
	}
}
