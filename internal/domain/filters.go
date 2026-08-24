package domain

import "strings"

type AssetFilter struct {
	Status    AssetStatus
	AssetType string
	Location  string
	Search    string
}
type PlanFilter struct {
	AssetID   string
	Status    PlanStatus
	DueBefore string
}
type TaskFilter struct {
	AssetID  string
	PlanID   string
	Status   TaskStatus
	Assignee string
}
type FindingFilter struct {
	AssetID  string
	TaskID   string
	Status   FindingStatus
	Severity FindingSeverity
}

func matchText(v, q string) bool {
	return q == "" || strings.Contains(strings.ToLower(v), strings.ToLower(q))
}
func (a Asset) Matches(f AssetFilter) bool {
	return (f.Status == "" || a.Status == f.Status) && (f.AssetType == "" || a.AssetType == f.AssetType) && (f.Location == "" || matchText(a.Location, f.Location)) && matchText(a.Name, f.Search)
}
func (p MaintenancePlan) Matches(f PlanFilter) bool {
	return (f.AssetID == "" || p.AssetID == f.AssetID) && (f.Status == "" || p.Status == f.Status)
}
func (t MaintenanceTask) Matches(f TaskFilter) bool {
	return (f.AssetID == "" || t.AssetID == f.AssetID) && (f.PlanID == "" || t.PlanID == f.PlanID) && (f.Status == "" || t.Status == f.Status) && t.Status != TaskBlocked && (f.Assignee == "" || t.Assignee == f.Assignee)
}
func (f Finding) Matches(q FindingFilter) bool {
	return (q.AssetID == "" || f.AssetID == q.AssetID) && (q.TaskID == "" || f.TaskID == q.TaskID) && (q.Status == "" || f.Status == q.Status) && (q.Severity == "" || f.Severity == q.Severity)
}
