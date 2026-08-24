package application

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/example/asset-maintenance-service/internal/domain"
	"github.com/example/asset-maintenance-service/internal/infrastructure/memory"
)

func schedulerFixture(t *testing.T, count int) (*Scheduler, *memory.PlanRepository, *memory.TaskRepository, time.Time) {
	t.Helper()
	now := time.Date(2026, 8, 24, 9, 0, 0, 0, time.UTC)
	plans := memory.NewPlanRepository()
	tasks := memory.NewTaskRepository()
	for i := 0; i < count; i++ {
		plan := domain.MaintenancePlan{ID: string(rune('a' + i)), AssetID: "asset", TemplateID: "template", Status: domain.PlanActive, NextRunAt: now.Add(-time.Hour)}
		if _, err := plans.Create(context.Background(), plan); err != nil {
			t.Fatal(err)
		}
	}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewScheduler(plans, tasks, time.Minute, log), plans, tasks, now
}

func TestDispatchDueConcurrentPlans(t *testing.T) {
	scheduler, _, tasks, now := schedulerFixture(t, 4)
	start := make(chan struct{})
	ready := make(chan struct{}, 1)
	done := make(chan int, 1)
	go func() {
		ready <- struct{}{}
		<-start
		done <- scheduler.DispatchDue(context.Background(), now)
	}()
	<-ready
	close(start)
	select {
	case got := <-done:
		if got != 4 {
			t.Fatalf("expected four created tasks, got %d", got)
		}
	case <-time.After(time.Second):
		t.Fatal("concurrent dispatch did not finish")
	}
	all, err := tasks.List(context.Background(), domain.TaskFilter{})
	if err != nil || len(all) != 4 {
		t.Fatalf("expected four persisted tasks, got %d, err=%v", len(all), err)
	}
}

func TestDispatchDueClosesResults(t *testing.T) {
	scheduler, _, _, now := schedulerFixture(t, 2)
	done := make(chan int, 1)
	go func() { done <- scheduler.DispatchDue(context.Background(), now) }()
	select {
	case got := <-done:
		if got != 2 {
			t.Fatalf("expected two results, got %d", got)
		}
	case <-time.After(time.Second):
		t.Fatal("dispatch did not close and drain its result channel")
	}
}

func TestDispatchDueCancellation(t *testing.T) {
	scheduler, _, _, now := schedulerFixture(t, 3)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := scheduler.DuePlans(ctx, now); err == nil {
		t.Fatal("expected canceled due-plan lookup to return an error")
	}
}
