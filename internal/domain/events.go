package domain

import (
	"context"
	"time"
)

type EventType string

const (
	EventTaskCreated       EventType = "task.created"
	EventTaskStatusChanged EventType = "task.status_changed"
	EventFindingReported   EventType = "finding.reported"
	EventPlanDue           EventType = "plan.due"
)

type Event struct {
	ID          string            `json:"id"`
	Type        EventType         `json:"type"`
	AggregateID string            `json:"aggregate_id"`
	Payload     map[string]string `json:"payload,omitempty"`
	At          time.Time         `json:"at"`
}
type EventPublisher interface {
	Publish(context.Context, Event) error
}
