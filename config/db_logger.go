package config

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/S4F4Y4T/goWebService/pkg/correlation"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type slogLogger struct {
	SlowThreshold time.Duration
}

func NewDBLogger(threshold time.Duration) logger.Interface {
	return &slogLogger{
		SlowThreshold: threshold,
	}
}

func (l *slogLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *slogLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	slog.Info(fmt.Sprintf(msg, data...))
}

func (l *slogLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	slog.Warn(fmt.Sprintf(msg, data...))
}

func (l *slogLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	slog.Error(fmt.Sprintf(msg, data...))
}

func (l *slogLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	correlationID := correlation.GetCorrelationID(ctx)

	fields := []any{
		"correlation_id", correlationID,
		"duration_ms", elapsed.Milliseconds(),
		"rows", rows,
		"sql", sql,
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fields = append(fields, "error", err)
		slog.Error("Database Error", fields...)
		return
	}

	if l.SlowThreshold > 0 && elapsed > l.SlowThreshold {
		slog.Warn("Slow Database Query", fields...)
	}
}
