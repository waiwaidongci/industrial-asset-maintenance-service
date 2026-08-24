package domain

import (
	"errors"
	"testing"
)

func TestOpenFindingRequiresAcknowledgementBeforeResolve(t *testing.T) {
	finding := &Finding{Status: FindingOpen}
	if err := TransitionFinding(finding, FindingResolved); !errors.Is(err, ErrTransition) {
		t.Fatalf("expected direct resolution to fail, got %v", err)
	}
	if finding.Status != FindingOpen {
		t.Fatalf("failed transition changed status to %s", finding.Status)
	}
}

func TestAcknowledgedFindingCanResolve(t *testing.T) {
	finding := &Finding{Status: FindingAcknowledged}
	if err := TransitionFinding(finding, FindingResolved); err != nil {
		t.Fatalf("resolve acknowledged finding: %v", err)
	}
	if finding.Status != FindingResolved {
		t.Fatalf("expected resolved, got %s", finding.Status)
	}
}

func TestFindingStatusActivityClassification(t *testing.T) {
	for _, status := range []FindingStatus{FindingOpen, FindingAcknowledged} {
		if !IsActiveFindingStatus(status) {
			t.Fatalf("expected %s to be active", status)
		}
	}
	for _, status := range []FindingStatus{FindingResolved, FindingIgnored, FindingStatus("unknown")} {
		if IsActiveFindingStatus(status) {
			t.Fatalf("expected %s to be inactive", status)
		}
	}
}

func TestFindingStatusRecognition(t *testing.T) {
	for _, status := range []FindingStatus{FindingOpen, FindingAcknowledged, FindingResolved, FindingIgnored} {
		if !IsKnownFindingStatus(status) {
			t.Fatalf("expected %s to be known", status)
		}
	}
	if IsKnownFindingStatus(FindingStatus("unknown")) {
		t.Fatal("unknown status was accepted")
	}
}
