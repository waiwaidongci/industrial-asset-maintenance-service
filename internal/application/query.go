package application

import (
	"context"
	"fmt"
	"github.com/example/asset-maintenance-service/internal/domain"
	"sort"
	"time"
)

type MaintenanceSummary struct {
	TotalAssets    int       `json:"total_assets"`
	ActiveAssets   int       `json:"active_assets"`
	OpenFindings   int       `json:"open_findings"`
	PendingTasks   int       `json:"pending_tasks"`
	CompletedTasks int       `json:"completed_tasks"`
	GeneratedAt    time.Time `json:"generated_at"`
}

func (s *Service) Summary(ctx context.Context) (MaintenanceSummary, error) {
	assets, err := s.Assets.List(ctx, domain.AssetFilter{})
	if err != nil {
		return MaintenanceSummary{}, fmt.Errorf("summary assets: %w", err)
	}
	tasks, err := s.Tasks.List(ctx, domain.TaskFilter{})
	if err != nil {
		return MaintenanceSummary{}, fmt.Errorf("summary tasks: %w", err)
	}
	findings, err := s.Findings.List(ctx, domain.FindingFilter{})
	if err != nil {
		return MaintenanceSummary{}, fmt.Errorf("summary findings: %w", err)
	}
	summary := MaintenanceSummary{TotalAssets: len(assets), GeneratedAt: s.Now()}
	for _, a := range assets {
		if a.Status == domain.AssetActive {
			summary.ActiveAssets++
		}
	}
	for _, t := range tasks {
		if t.Status == domain.TaskPending {
			summary.PendingTasks++
		}
		if t.Status == domain.TaskCompleted {
			summary.CompletedTasks++
		}
	}
	for _, f := range findings {
		if f.Status == domain.FindingOpen || f.Status == domain.FindingAcknowledged {
			summary.OpenFindings++
		}
	}
	return summary, nil
}

func SortTasksByPriority(tasks []domain.MaintenanceTask, findings []domain.Finding) []domain.MaintenanceTask {
	severity := map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1}
	priority := make(map[string]int)
	for _, finding := range findings {
		if finding.Status != domain.FindingResolved && severity[string(finding.Severity)] > priority[finding.TaskID] {
			priority[finding.TaskID] = severity[string(finding.Severity)]
		}
	}
	result := tasks
	sort.SliceStable(result, func(i, j int) bool {
		if priority[result[i].ID] == priority[result[j].ID] {
			return result[i].ScheduledAt.Before(result[j].ScheduledAt)
		}
		return priority[result[i].ID] > priority[result[j].ID]
	})
	return result
}
