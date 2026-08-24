package memory

import (
	"context"
	"testing"
	"time"

	"github.com/example/asset-maintenance-service/internal/domain"
)

func TestHistoryListIsolatesEntries(t *testing.T) {
	r := NewHistoryRepository()
	entry := domain.HistoryEntry{ID: "h-1", TaskID: "task-1", Action: "created", Details: map[string]string{"actor": "ops"}, OccurredAt: time.Unix(1, 0)}
	if err := r.Append(context.Background(), entry); err != nil {
		t.Fatal(err)
	}
	list, err := r.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	list[0].Action = "tampered"
	list[0].Details["actor"] = "attacker"
	again, err := r.List(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if again[0].Action != "created" || again[0].Details["actor"] != "ops" {
		t.Fatalf("history changed through returned snapshot: %#v", again[0])
	}
}
