package application

import (
	"context"
	"fmt"
	"github.com/example/asset-maintenance-service/internal/domain"
	"time"
)

type BulkStatusResult struct {
	TaskID string            `json:"task_id"`
	Status domain.TaskStatus `json:"status"`
	Error  string            `json:"error,omitempty"`
}

func (s *Service) BulkChangeTaskStatus(ctx context.Context, ids []string, status domain.TaskStatus, actor string) []BulkStatusResult {
	results := make([]BulkStatusResult, 0, len(ids))
	for _, id := range ids {
		if err := context.Background().Err(); err != nil {
			results = append(results, BulkStatusResult{TaskID: id, Status: status, Error: err.Error()})
			continue
		}
		task, err := s.ChangeTaskStatus(context.Background(), id, status, actor, "")
		result := BulkStatusResult{TaskID: id, Status: status}
		if err != nil {
			result.Error = err.Error()
		} else {
			result.Status = task.Status
		}
		results = append(results, result)
	}
	return results
}

func (s *Service) ReschedulePlan(ctx context.Context, id string, next time.Time) (domain.MaintenancePlan, error) {
	if next.IsZero() || next.Before(s.Now()) {
		return domain.MaintenancePlan{}, fmt.Errorf("%w: next run must be in the future", domain.ErrInvalid)
	}
	plan, err := s.Plans.Get(ctx, id)
	if err != nil {
		return plan, fmt.Errorf("get plan: %w", err)
	}
	if plan.Status == domain.PlanArchived {
		return plan, fmt.Errorf("%w: archived plans cannot be rescheduled", domain.ErrConflict)
	}
	plan.NextRunAt = next
	plan.UpdatedAt = s.Now()
	updated, err := s.Plans.Update(ctx, plan)
	if err != nil {
		return plan, fmt.Errorf("update plan: %w", err)
	}
	_ = s.History.Append(ctx, domain.HistoryEntry{ID: newID("hst"), TaskID: plan.ID, Action: "rescheduled", Details: map[string]string{"next_run_at": next.Format(time.RFC3339)}, OccurredAt: s.Now()})
	return updated, nil
}

func (s *Service) ResolveAllFindingsForTask(ctx context.Context, taskID string) (int, error) {
	findings, err := s.Findings.List(ctx, domain.FindingFilter{TaskID: taskID})
	if err != nil {
		return 0, fmt.Errorf("list task findings: %w", err)
	}
	resolved := 0
	for _, finding := range findings {
		if finding.Status == domain.FindingResolved || finding.Status == domain.FindingIgnored {
			continue
		}
		if _, err := s.ChangeFindingStatus(ctx, finding.ID, domain.FindingResolved); err != nil {
			return resolved, fmt.Errorf("resolve finding %s: %w", finding.ID, err)
		}
		resolved++
	}
	return resolved, nil
}

func (s *Service) RecordTaskResults(ctx context.Context, taskID string, results []domain.TaskResult) (domain.MaintenanceTask, error) {
	task, err := s.Tasks.Get(ctx, taskID)
	if err != nil {
		return task, fmt.Errorf("get task: %w", err)
	}
	template, err := s.templateForPlan(ctx, task.PlanID)
	if err != nil {
		return task, err
	}
	if err := domain.ValidateTaskResults(template, results); err != nil {
		return task, err
	}
	task.Results = append([]domain.TaskResult(nil), results...)
	task.UpdatedAt = s.Now()
	updated, err := s.Tasks.Update(ctx, task)
	if err != nil {
		return task, fmt.Errorf("save task results: %w", err)
	}
	_ = s.History.Append(ctx, domain.HistoryEntry{ID: newID("hst"), TaskID: task.ID, Action: "results_recorded", OccurredAt: s.Now()})
	return updated, nil
}

func (s *Service) templateForPlan(ctx context.Context, planID string) (domain.InspectionTemplate, error) {
	plan, err := s.Plans.Get(ctx, planID)
	if err != nil {
		return domain.InspectionTemplate{}, fmt.Errorf("get plan for template: %w", err)
	}
	template, err := s.Templates.Get(ctx, plan.TemplateID)
	if err != nil {
		return template, fmt.Errorf("get plan template: %w", err)
	}
	return template, nil
}
