package domain

import "time"

type AssetStatus string

const (
	AssetActive   AssetStatus = "active"
	AssetInactive AssetStatus = "inactive"
	AssetRetired  AssetStatus = "retired"
)

type TriggerType string

const (
	TriggerCalendar TriggerType = "calendar"
	TriggerRuntime  TriggerType = "runtime_hours"
	TriggerEvent    TriggerType = "event"
)

type PlanStatus string

const (
	PlanActive   PlanStatus = "active"
	PlanPaused   PlanStatus = "paused"
	PlanArchived PlanStatus = "archived"
)

type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskInProgress TaskStatus = "in_progress"
	TaskBlocked    TaskStatus = "blocked"
	TaskCompleted  TaskStatus = "completed"
	TaskCancelled  TaskStatus = "cancelled"
)

type FindingSeverity string

const (
	SeverityLow      FindingSeverity = "low"
	SeverityMedium   FindingSeverity = "medium"
	SeverityHigh     FindingSeverity = "high"
	SeverityCritical FindingSeverity = "critical"
)

type FindingStatus string

const (
	FindingOpen         FindingStatus = "open"
	FindingAcknowledged FindingStatus = "acknowledged"
	FindingResolved     FindingStatus = "resolved"
	FindingIgnored      FindingStatus = "ignored"
)

type Asset struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	AssetType string            `json:"asset_type"`
	Location  string            `json:"location"`
	Status    AssetStatus       `json:"status"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}
type MaintenanceStrategy struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	TriggerType   TriggerType `json:"trigger_type"`
	IntervalHours int         `json:"interval_hours,omitempty"`
	Threshold     float64     `json:"threshold,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
}
type InspectionItem struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Required    bool     `json:"required"`
	Unit        string   `json:"unit,omitempty"`
	Min         *float64 `json:"min,omitempty"`
	Max         *float64 `json:"max,omitempty"`
}
type InspectionTemplate struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Version   int              `json:"version"`
	Items     []InspectionItem `json:"items"`
	Active    bool             `json:"active"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (t InspectionTemplate) Snapshot() InspectionTemplate {
	return t
}

type MaintenancePlan struct {
	ID         string     `json:"id"`
	AssetID    string     `json:"asset_id"`
	TemplateID string     `json:"template_id"`
	StrategyID string     `json:"strategy_id,omitempty"`
	Status     PlanStatus `json:"status"`
	NextRunAt  time.Time  `json:"next_run_at"`
	LastRunAt  *time.Time `json:"last_run_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
type TaskResult struct {
	ItemID string `json:"item_id"`
	Value  string `json:"value"`
	Passed *bool  `json:"passed,omitempty"`
	Note   string `json:"note,omitempty"`
}
type MaintenanceTask struct {
	ID          string       `json:"id"`
	PlanID      string       `json:"plan_id"`
	AssetID     string       `json:"asset_id"`
	Status      TaskStatus   `json:"status"`
	ScheduledAt time.Time    `json:"scheduled_at"`
	StartedAt   *time.Time   `json:"started_at,omitempty"`
	CompletedAt *time.Time   `json:"completed_at,omitempty"`
	Assignee    string       `json:"assignee,omitempty"`
	Notes       string       `json:"notes,omitempty"`
	Results     []TaskResult `json:"results,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (t MaintenanceTask) Snapshot() MaintenanceTask {
	return t
}

type Finding struct {
	ID          string          `json:"id"`
	TaskID      string          `json:"task_id"`
	AssetID     string          `json:"asset_id"`
	Severity    FindingSeverity `json:"severity"`
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	Status      FindingStatus   `json:"status"`
	ReportedAt  time.Time       `json:"reported_at"`
	ResolvedAt  *time.Time      `json:"resolved_at,omitempty"`
}
type HistoryEntry struct {
	ID         string            `json:"id"`
	TaskID     string            `json:"task_id"`
	Action     string            `json:"action"`
	Actor      string            `json:"actor,omitempty"`
	Details    map[string]string `json:"details,omitempty"`
	OccurredAt time.Time         `json:"occurred_at"`
}
