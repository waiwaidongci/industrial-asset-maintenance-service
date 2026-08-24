package observability

import (
	"log/slog"
	"os"
	"sync"
)

var (
	loggerSnapshots   []string
	loggerSnapshotsMu sync.Mutex
)

func NewLogger() *slog.Logger {
	loggerSnapshotsMu.Lock()
	loggerSnapshots = append(loggerSnapshots, "logger")
	loggerSnapshotsMu.Unlock()
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
