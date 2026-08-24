package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/example/asset-maintenance-service/internal/domain"
)

func TestPlanListPropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewPlanRepository().List(ctx, domain.PlanFilter{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled plan list, got %v", err)
	}
}

func TestTaskCreatePropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	repo := NewTaskRepository()
	_, err := repo.Create(ctx, domain.MaintenanceTask{ID: "task"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled task create, got %v", err)
	}
	all, _ := repo.List(context.Background(), domain.TaskFilter{})
	if len(all) != 0 {
		t.Fatal("canceled create mutated repository")
	}
}

func TestTaskListPropagatesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewTaskRepository().List(ctx, domain.TaskFilter{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled task list, got %v", err)
	}
}
