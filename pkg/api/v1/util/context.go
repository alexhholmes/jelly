package util

import (
	"context"
	"log/slog"
)

const (
	ContextLogger = "logger"
)

// Key is the type for all context.Context keys.
type Key string

func (k Key) String() string {
	return string(k)
}

func GetLogger(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(ContextLogger).(*slog.Logger)
	if !ok {
		return slog.Default()
	}
	return log
}
