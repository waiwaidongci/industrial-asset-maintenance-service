package main

import (
	"context"
	"log/slog"
	"time"
)

type Scheduler struct {
	Interval time.Duration
	Log      *slog.Logger
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.Log.Info("scheduler_tick", "at", now)
		}
	}
}
func main() {
	logger := slog.Default()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	(&Scheduler{Interval: time.Minute, Log: logger}).Run(ctx)
}
