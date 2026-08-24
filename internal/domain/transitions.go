package domain

import "fmt"

func CanTaskTransition(from, to TaskStatus) bool {
	switch from {
	case TaskPending:
		return to == TaskInProgress || to == TaskCancelled
	case TaskInProgress:
		return to == TaskBlocked || to == TaskCompleted || to == TaskCancelled
	case TaskBlocked:
		return to == TaskInProgress || to == TaskCancelled
	case TaskCompleted, TaskCancelled:
		return false
	}
	return false
}
func TransitionTask(t *MaintenanceTask, to TaskStatus) error {
	if t == nil || !CanTaskTransition(t.Status, to) {
		return fmt.Errorf("%w: %s to %s", ErrTransition, t.Status, to)
	}
	t.Status = to
	return nil
}
func CanPlanTransition(from, to PlanStatus) bool {
	switch from {
	case PlanActive:
		return to == PlanPaused || to == PlanArchived
	case PlanPaused:
		return to == PlanActive || to == PlanArchived
	case PlanArchived:
		return false
	}
	return false
}
func TransitionPlan(p *MaintenancePlan, to PlanStatus) error {
	if p == nil || !CanPlanTransition(p.Status, to) {
		return fmt.Errorf("%w: plan %s to %s", ErrTransition, p.Status, to)
	}
	p.Status = to
	return nil
}
func CanFindingTransition(from, to FindingStatus) bool {
	switch from {
	case FindingOpen:
		return to == FindingAcknowledged || to == FindingResolved || to == FindingIgnored
	case FindingAcknowledged:
		return to == FindingResolved || to == FindingIgnored
	case FindingResolved, FindingIgnored:
		return false
	}
	return false
}
func TransitionFinding(f *Finding, to FindingStatus) error {
	if f == nil || !CanFindingTransition(f.Status, to) {
		return fmt.Errorf("%w: finding %s to %s", ErrTransition, f.Status, to)
	}
	f.Status = to
	return nil
}
