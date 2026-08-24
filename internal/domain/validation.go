package domain

import "strings"

func (a Asset) Validate() error {
	if strings.TrimSpace(a.Name) == "" || strings.TrimSpace(a.AssetType) == "" {
		return ErrInvalid
	}
	if a.Status != "" && a.Status != AssetActive && a.Status != AssetInactive && a.Status != AssetRetired {
		return ErrInvalid
	}
	return nil
}
func (s MaintenanceStrategy) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return ErrInvalid
	}
	if s.TriggerType != TriggerCalendar && s.TriggerType != TriggerRuntime && s.TriggerType != TriggerEvent {
		return ErrInvalid
	}
	if s.TriggerType == TriggerCalendar && s.IntervalHours <= 0 {
		return ErrInvalid
	}
	return nil
}
func (t InspectionTemplate) Validate() error {
	if strings.TrimSpace(t.Name) == "" || len(t.Items) == 0 {
		return ErrInvalid
	}
	if t.Version <= 0 {
		return ErrInvalid
	}
	seen := map[string]bool{}
	for _, i := range t.Items {
		if strings.TrimSpace(i.ID) == "" || strings.TrimSpace(i.Name) == "" || seen[i.ID] {
			return ErrInvalid
		}
		seen[i.ID] = true
		if i.Min != nil && i.Max != nil && *i.Min > *i.Max {
			return ErrInvalid
		}
	}
	return nil
}
func (p MaintenancePlan) Validate() error {
	if p.AssetID == "" || p.TemplateID == "" || p.NextRunAt.IsZero() {
		return ErrInvalid
	}
	if p.Status != "" && p.Status != PlanActive && p.Status != PlanPaused && p.Status != PlanArchived {
		return ErrInvalid
	}
	return nil
}
func (f Finding) Validate() error {
	if f.TaskID == "" || f.AssetID == "" || strings.TrimSpace(f.Title) == "" {
		return ErrInvalid
	}
	if f.Severity != SeverityLow && f.Severity != SeverityMedium && f.Severity != SeverityHigh && f.Severity != SeverityCritical {
		return ErrInvalid
	}
	return nil
}
