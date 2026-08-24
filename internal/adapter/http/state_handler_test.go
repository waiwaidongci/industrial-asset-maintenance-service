package httpadapter

import (
	"bytes"
	"context"
	"github.com/example/asset-maintenance-service/internal/application"
	"github.com/example/asset-maintenance-service/internal/domain"
	"github.com/example/asset-maintenance-service/internal/infrastructure/memory"
	"github.com/example/asset-maintenance-service/internal/platform/observability"
	"net/http"
	"net/http/httptest"
	"testing"
)

func seededFindingHandler(t *testing.T) (*Handler, *memory.FindingRepository) {
	t.Helper()
	repo := memory.NewFindingRepository()
	_, err := repo.Create(context.Background(), domain.Finding{ID: "finding-1", TaskID: "task-1", AssetID: "asset-1", Severity: domain.SeverityHigh, Title: "overheat", Status: domain.FindingOpen})
	if err != nil {
		t.Fatal(err)
	}
	svc := &application.Service{Findings: repo}
	return NewHandler(svc, &observability.Metrics{}), repo
}

func TestTaskStatusHandlerUsesDomainTransition(t *testing.T) {
	repo := memory.NewTaskRepository()
	_, err := repo.Create(context.Background(), domain.MaintenanceTask{ID: "task-1", PlanID: "plan-1", AssetID: "asset-1", Status: domain.TaskPending})
	if err != nil {
		t.Fatal(err)
	}
	h := NewHandler(&application.Service{Tasks: repo}, &observability.Metrics{})
	req := httptest.NewRequest(http.MethodPost, "/tasks/task-1/status", bytes.NewBufferString(`{"status":"corrupted"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	if res.Code == http.StatusOK {
		t.Fatalf("illegal task transition returned HTTP %d", res.Code)
	}
}

func TestTaskStateHTTPRejectsIllegalTransition(t *testing.T) {
	h, _ := seededFindingHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/findings/finding-1/status", bytes.NewBufferString(`{"status":"corrupted"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	if res.Code == http.StatusOK {
		t.Fatalf("illegal transition returned HTTP %d", res.Code)
	}
}

func TestFindingStateHTTPConsistency(t *testing.T) {
	h, repo := seededFindingHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/findings/finding-1/status", bytes.NewBufferString(`{"status":"corrupted"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()
	h.ServeHTTP(res, req)
	f, err := repo.Get(context.Background(), "finding-1")
	if err != nil {
		t.Fatal(err)
	}
	if f.Status != domain.FindingOpen {
		t.Fatalf("invalid request changed status to %q", f.Status)
	}
}
