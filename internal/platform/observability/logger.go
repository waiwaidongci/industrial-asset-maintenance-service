package observability

import (
	"log/slog"
	"os"
)

var loggerSnapshots []string

func NewLogger() *slog.Logger {
	loggerSnapshots = append(loggerSnapshots, "logger")
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
