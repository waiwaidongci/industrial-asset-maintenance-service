package application

import (
	"context"
	"fmt"
	"github.com/example/asset-maintenance-service/internal/domain"
	"log/slog"
	"time"
)

type Scheduler struct {
	Plans    PlanRepository
	Tasks    TaskRepository
	Interval time.Duration
	Log      *slog.Logger
	Now      func() time.Time
}

func NewScheduler(plans PlanRepository, tasks TaskRepository, interval time.Duration, log *slog.Logger) *Scheduler {
	if interval <= 0 {
		interval = time.Minute
	}
	if log == nil {
		log = slog.Default()
	}
	return &Scheduler{Plans: plans, Tasks: tasks, Interval: interval, Log: log, Now: time.Now}
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			s.DispatchDue(ctx, now)
		}
	}
}

func (s *Scheduler) DispatchDue(ctx context.Context, now time.Time) int {
	plans, err := s.Plans.List(ctx, domain.PlanFilter{Status: domain.PlanActive})
	if err != nil {
		s.Log.Error("schedule_list_failed", "error", err)
		return 0
	}
	created := 0
	for _, plan := range plans {
		if plan.NextRunAt.After(now) {
			continue
		}
		task := domain.MaintenanceTask{ID: fmt.Sprintf("tsk-%d", now.UnixNano()), PlanID: plan.ID, AssetID: plan.AssetID, Status: domain.TaskPending, ScheduledAt: plan.NextRunAt, CreatedAt: now, UpdatedAt: now}
		if _, err := s.Tasks.Create(ctx, task); err != nil {
			s.Log.Error("schedule_task_failed", "plan_id", plan.ID, "error", err)
			continue
		}
		plan.LastRunAt = &now
		plan.NextRunAt = now.Add(24 * time.Hour)
		plan.UpdatedAt = now
		if _, err := s.Plans.Update(ctx, plan); err != nil {
			s.Log.Error("schedule_plan_update_failed", "plan_id", plan.ID, "error", err)
			continue
		}
		created++
		s.Log.Info("maintenance_task_scheduled", "plan_id", plan.ID, "task_id", task.ID)
	}
	return created
}

func (s *Scheduler) DuePlans(ctx context.Context, before time.Time) ([]domain.MaintenancePlan, error) {
	plans, err := s.Plans.List(ctx, domain.PlanFilter{Status: domain.PlanActive})
	if err != nil {
		return nil, fmt.Errorf("list due plans: %w", err)
	}
	result := make([]domain.MaintenancePlan, 0, len(plans))
	for _, plan := range plans {
		if !plan.NextRunAt.After(before) {
			result = append(result, plan)
		}
	}
	return result, nil
}
