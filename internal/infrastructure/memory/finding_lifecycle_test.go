package memory

import (
	"context"
	"errors"
	"testing"

	"github.com/example/asset-maintenance-service/internal/domain"
)

func createFinding(t *testing.T, repo *FindingRepository, finding domain.Finding) {
	t.Helper()
	if _, err := repo.Create(context.Background(), finding); err != nil {
		t.Fatalf("create finding: %v", err)
	}
}

func TestDuplicateActiveFindingReturnsConflict(t *testing.T) {
	for _, status := range []domain.FindingStatus{domain.FindingOpen, domain.FindingAcknowledged} {
		repo := NewFindingRepository()
		createFinding(t, repo, domain.Finding{ID: "old", TaskID: "task-1", Title: "leak", Status: status})
		_, err := repo.Create(context.Background(), domain.Finding{ID: "new", TaskID: "task-1", Title: "leak", Status: domain.FindingOpen})
		if !errors.Is(err, domain.ErrConflict) {
			t.Fatalf("status %s: expected conflict, got %v", status, err)
		}
	}
}

func TestTerminalFindingAllowsReplacement(t *testing.T) {
	for _, status := range []domain.FindingStatus{domain.FindingResolved, domain.FindingIgnored} {
		repo := NewFindingRepository()
		createFinding(t, repo, domain.Finding{ID: "old", TaskID: "task-1", Title: "leak", Status: status})
		created, err := repo.Create(context.Background(), domain.Finding{ID: "new", TaskID: "task-1", Title: "leak", Status: domain.FindingOpen})
		if err != nil {
			t.Fatalf("status %s: create replacement: %v", status, err)
		}
		if created.ID != "new" {
			t.Fatalf("status %s: expected new finding, got %s", status, created.ID)
		}
	}
}

func TestCreateRejectsUnknownFindingStatus(t *testing.T) {
	repo := NewFindingRepository()
	_, err := repo.Create(context.Background(), domain.Finding{ID: "bad", Status: domain.FindingStatus("unknown")})
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected invalid status error, got %v", err)
	}
}

func TestUpdateRejectsUnknownFindingStatus(t *testing.T) {
	repo := NewFindingRepository()
	createFinding(t, repo, domain.Finding{ID: "finding-1", Status: domain.FindingOpen})
	_, err := repo.Update(context.Background(), domain.Finding{ID: "finding-1", Status: domain.FindingStatus("unknown")})
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected invalid status error, got %v", err)
	}
}

func TestListRejectsStoredUnknownFindingStatus(t *testing.T) {
	repo := NewFindingRepository()
	repo.values["legacy"] = domain.Finding{ID: "legacy", Status: domain.FindingStatus("unknown")}
	_, err := repo.List(context.Background(), domain.FindingFilter{})
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("expected invalid status error, got %v", err)
	}
}
