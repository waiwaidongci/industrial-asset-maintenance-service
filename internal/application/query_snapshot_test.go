package application

import (
	"testing"
	"time"

	"github.com/example/asset-maintenance-service/internal/domain"
)

func TestSortTasksByPriorityDoesNotMutateInput(t *testing.T) {
	tasks := []domain.MaintenanceTask{
		{ID: "low", ScheduledAt: time.Unix(1, 0)},
		{ID: "high", ScheduledAt: time.Unix(2, 0)},
	}
	findings := []domain.Finding{{TaskID: "high", Severity: domain.SeverityCritical, Status: domain.FindingOpen}}
	_ = SortTasksByPriority(tasks, findings)
	if tasks[0].ID != "low" || tasks[1].ID != "high" {
		t.Fatalf("input task order mutated: %#v", tasks)
	}
}
