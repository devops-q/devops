package logger

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

type Logger struct {
	*slog.Logger
	service string
}

var instance *Logger

func Init() *Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   a.Key,
					Value: slog.StringValue(a.Value.Time().Format(time.RFC3339)),
				}
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	slogger := slog.New(handler)
	slog.SetDefault(slogger)

	instance = &Logger{
		Logger:  slogger,
		service: "default",
	}
	return instance
}

func GetLogger() *Logger {
	if instance == nil {
		instance = Init()
	}
	return instance
}

func (l *Logger) WithService(service string) *Logger {
	return &Logger{
		Logger:  l.Logger,
		service: service,
	}
}

func (l *Logger) logWithService(level slog.Level, msg string, args ...interface{}) {
	attrs := []any{"context", l.service}
	attrs = append(attrs, args...)
	l.Logger.Log(nil, level, msg, attrs...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logWithService(slog.LevelDebug, msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.logWithService(slog.LevelInfo, msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logWithService(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.logWithService(slog.LevelError, msg, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.logWithService(slog.LevelError, msg)
	os.Exit(1)
}
