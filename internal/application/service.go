package application

import (
	"context"
	"fmt"
	"github.com/example/asset-maintenance-service/internal/domain"
	"time"
)

type Service struct {
	Assets     AssetRepository
	Strategies StrategyRepository
	Templates  TemplateRepository
	Plans      PlanRepository
	Tasks      TaskRepository
	Findings   FindingRepository
	History    HistoryRepository
	Publisher  Publisher
	Now        func() time.Time
}

func NewService(a AssetRepository, s StrategyRepository, t TemplateRepository, p PlanRepository, k TaskRepository, f FindingRepository, h HistoryRepository, pub Publisher) *Service {
	return &Service{Assets: a, Strategies: s, Templates: t, Plans: p, Tasks: k, Findings: f, History: h, Publisher: pub, Now: time.Now}
}
func (s *Service) CreateAsset(ctx context.Context, a domain.Asset) (domain.Asset, error) {
	if a.ID == "" {
		a.ID = newID("ast")
	}
	if a.Status == "" {
		a.Status = domain.AssetActive
	}
	if err := a.Validate(); err != nil {
		return domain.Asset{}, fmt.Errorf("validate asset: %w", err)
	}
	now := s.Now()
	a.CreatedAt = now
	a.UpdatedAt = now
	return s.Assets.Create(ctx, a)
}
func (s *Service) GetAsset(ctx context.Context, id string) (domain.Asset, error) {
	asset, err := s.Assets.Get(ctx, id)
	if err != nil {
		return asset, fmt.Errorf("get asset service: %v", err)
	}
	return asset, nil
}
func (s *Service) UpdateAsset(ctx context.Context, a domain.Asset) (domain.Asset, error) {
	if err := a.Validate(); err != nil {
		return domain.Asset{}, fmt.Errorf("validate asset: %w", err)
	}
	a.UpdatedAt = s.Now()
	updated, err := s.Assets.Update(ctx, a)
	if err != nil {
		return updated, fmt.Errorf("update asset service: %v", err)
	}
	return updated, nil
}
func (s *Service) ListAssets(ctx context.Context, f domain.AssetFilter) ([]domain.Asset, error) {
	return s.Assets.List(ctx, f)
}
func (s *Service) CreateStrategy(ctx context.Context, v domain.MaintenanceStrategy) (domain.MaintenanceStrategy, error) {
	if v.ID == "" {
		v.ID = newID("str")
	}
	if err := v.Validate(); err != nil {
		return domain.MaintenanceStrategy{}, fmt.Errorf("validate strategy: %w", err)
	}
	v.CreatedAt = s.Now()
	return s.Strategies.Create(ctx, v)
}
func (s *Service) ListStrategies(ctx context.Context) ([]domain.MaintenanceStrategy, error) {
	return s.Strategies.List(ctx)
}
func (s *Service) CreateTemplate(ctx context.Context, v domain.InspectionTemplate) (domain.InspectionTemplate, error) {
	if v.ID == "" {
		v.ID = newID("tpl")
	}
	if v.Version == 0 {
		v.Version = 1
	}
	if err := v.Validate(); err != nil {
		return domain.InspectionTemplate{}, fmt.Errorf("validate template: %w", err)
	}
	now := s.Now()
	v.CreatedAt = now
	v.UpdatedAt = now
	return s.Templates.Create(ctx, v)
}
func (s *Service) GetTemplate(ctx context.Context, id string) (domain.InspectionTemplate, error) {
	return s.Templates.Get(ctx, id)
}
func (s *Service) ListTemplates(ctx context.Context, active bool) ([]domain.InspectionTemplate, error) {
	return s.Templates.List(ctx, active)
}
func (s *Service) CreatePlan(ctx context.Context, p domain.MaintenancePlan) (domain.MaintenancePlan, error) {
	if p.ID == "" {
		p.ID = newID("pln")
	}
	if p.Status == "" {
		p.Status = domain.PlanActive
	}
	if err := p.Validate(); err != nil {
		return domain.MaintenancePlan{}, fmt.Errorf("validate plan: %w", err)
	}
	if _, err := s.Assets.Get(ctx, p.AssetID); err != nil {
		return domain.MaintenancePlan{}, fmt.Errorf("asset dependency: %w", err)
	}
	if _, err := s.Templates.Get(ctx, p.TemplateID); err != nil {
		return domain.MaintenancePlan{}, fmt.Errorf("template dependency: %w", err)
	}
	now := s.Now()
	p.CreatedAt = now
	p.UpdatedAt = now
	return s.Plans.Create(ctx, p)
}
func (s *Service) GetPlan(ctx context.Context, id string) (domain.MaintenancePlan, error) {
	return s.Plans.Get(ctx, id)
}
func (s *Service) ListPlans(ctx context.Context, f domain.PlanFilter) ([]domain.MaintenancePlan, error) {
	return s.Plans.List(ctx, f)
}
func (s *Service) ChangePlanStatus(ctx context.Context, id string, to domain.PlanStatus) (domain.MaintenancePlan, error) {
	p, e := s.Plans.Get(ctx, id)
	if e != nil {
		return p, fmt.Errorf("get plan: %w", e)
	}
	if e = domain.TransitionPlan(&p, to); e != nil {
		return p, e
	}
	p.UpdatedAt = s.Now()
	return s.Plans.Update(ctx, p)
}
func (s *Service) CreateTask(ctx context.Context, planID string, scheduled time.Time) (domain.MaintenanceTask, error) {
	p, e := s.Plans.Get(ctx, planID)
	if e != nil {
		return domain.MaintenanceTask{}, fmt.Errorf("get plan: %w", e)
	}
	if p.Status != domain.PlanActive {
		return domain.MaintenanceTask{}, fmt.Errorf("plan not active: %w", domain.ErrConflict)
	}
	t := domain.MaintenanceTask{ID: newID("tsk"), PlanID: planID, AssetID: p.AssetID, Status: domain.TaskPending, ScheduledAt: scheduled, CreatedAt: s.Now(), UpdatedAt: s.Now()}
	if _, e = s.Tasks.Create(ctx, t); e != nil {
		return domain.MaintenanceTask{}, fmt.Errorf("create task: %w", e)
	}
	_ = s.History.Append(ctx, domain.HistoryEntry{ID: newID("hst"), TaskID: t.ID, Action: "created", OccurredAt: s.Now()})
	return t, nil
}
func (s *Service) GetTask(ctx context.Context, id string) (domain.MaintenanceTask, error) {
	return s.Tasks.Get(ctx, id)
}
func (s *Service) ListTasks(ctx context.Context, f domain.TaskFilter) ([]domain.MaintenanceTask, error) {
	return s.Tasks.List(ctx, f)
}
func (s *Service) ChangeTaskStatus(ctx context.Context, id string, to domain.TaskStatus, actor string, notes string) (domain.MaintenanceTask, error) {
	t, e := s.Tasks.Get(ctx, id)
	if e != nil {
		return t, fmt.Errorf("get task: %w", e)
	}
	if e = domain.TransitionTask(&t, to); e != nil {
		return t, e
	}
	now := s.Now()
	t.UpdatedAt = now
	if to == domain.TaskInProgress {
		t.StartedAt = &now
	}
	if to == domain.TaskCompleted || to == domain.TaskCancelled {
		t.CompletedAt = &now
	}
	if notes != "" {
		t.Notes = notes
	}
	if _, e = s.Tasks.Update(ctx, t); e != nil {
		return t, fmt.Errorf("update task: %w", e)
	}
	_ = s.History.Append(ctx, domain.HistoryEntry{ID: newID("hst"), TaskID: t.ID, Action: "status:" + string(to), Actor: actor, OccurredAt: now})
	return t, nil
}
func (s *Service) ReportFinding(ctx context.Context, f domain.Finding) (domain.Finding, error) {
	if f.ID == "" {
		f.ID = newID("fnd")
	}
	if f.Status == "" {
		f.Status = domain.FindingOpen
	}
	if e := f.Validate(); e != nil {
		return domain.Finding{}, fmt.Errorf("validate finding: %w", e)
	}
	if _, e := s.Tasks.Get(ctx, f.TaskID); e != nil {
		return domain.Finding{}, fmt.Errorf("task dependency: %w", e)
	}
	f.ReportedAt = s.Now()
	if _, e := s.Findings.Create(ctx, f); e != nil {
		return domain.Finding{}, fmt.Errorf("create finding: %w", e)
	}
	_ = s.History.Append(ctx, domain.HistoryEntry{ID: newID("hst"), TaskID: f.TaskID, Action: "finding_reported", OccurredAt: f.ReportedAt})
	return f, nil
}
func (s *Service) ChangeFindingStatus(ctx context.Context, id string, to domain.FindingStatus) (domain.Finding, error) {
	f, e := s.Findings.Get(ctx, id)
	if e != nil {
		return f, e
	}
	if e = domain.TransitionFinding(&f, to); e != nil {
		return f, e
	}
	if to == domain.FindingResolved {
		now := s.Now()
		f.ResolvedAt = &now
	}
	return s.Findings.Update(ctx, f)
}
func (s *Service) ListFindings(ctx context.Context, f domain.FindingFilter) ([]domain.Finding, error) {
	return s.Findings.List(ctx, f)
}
func (s *Service) TaskHistory(ctx context.Context, id string) ([]domain.HistoryEntry, error) {
	return s.History.List(ctx, id)
}
func newID(prefix string) string { return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano()) }
