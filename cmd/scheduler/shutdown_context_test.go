package main

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestSchedulerStopsOnSignal(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	done := make(chan struct{})
	go func() { (&Scheduler{Interval: time.Millisecond, Log: slog.Default()}).Run(ctx); close(done) }()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("scheduler ignored shutdown cancellation")
	}
}
