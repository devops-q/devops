package database

import (
	"context"
	"time"

	gormLogger "gorm.io/gorm/logger"
	"itu-minitwit/pkg/logger"
)

type GormLogger struct {
	logger *logger.Logger
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		logger: logger.GetLogger().WithService("gorm"),
	}
}

func (l *GormLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	l.logger.WithGroup("gorm").Info(msg, "args", args)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	l.logger.WithGroup("gorm").Warn(msg, "args", args)
}

func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	l.logger.WithGroup("gorm").Error(msg, "args", args)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	logAttrs := []any{
		"elapsed", elapsed,
		"rows", rows,
		"sql", sql,
	}
	if err != nil {
		logAttrs = append(logAttrs, "error", err)
		l.logger.WithGroup("gorm").Error("query", logAttrs...)
		return
	}
}
