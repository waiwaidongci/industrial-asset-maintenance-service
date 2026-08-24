package domain

import "testing"

func TestBlockedTaskRemainsVisible(t *testing.T) {
	task := MaintenanceTask{ID: "task-1", Status: TaskBlocked}
	if !task.Matches(TaskFilter{Status: TaskBlocked}) {
		t.Fatal("blocked task disappeared from its status filter")
	}
}

func TestInvalidFindingDoesNotMutateState(t *testing.T) {
	finding := Finding{TaskID: "task-1", AssetID: "asset-1", Severity: SeverityHigh, Title: "overheat", Status: FindingStatus("corrupted")}
	if err := finding.Validate(); err != ErrInvalid {
		t.Fatalf("Validate error = %v, want ErrInvalid", err)
	}
}
