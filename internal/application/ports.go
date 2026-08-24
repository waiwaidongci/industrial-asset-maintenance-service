package application

import (
	"context"
	"github.com/example/asset-maintenance-service/internal/domain"
)

type AssetRepository interface {
	Create(context.Context, domain.Asset) (domain.Asset, error)
	Get(context.Context, string) (domain.Asset, error)
	Update(context.Context, domain.Asset) (domain.Asset, error)
	List(context.Context, domain.AssetFilter) ([]domain.Asset, error)
}
type StrategyRepository interface {
	Create(context.Context, domain.MaintenanceStrategy) (domain.MaintenanceStrategy, error)
	Get(context.Context, string) (domain.MaintenanceStrategy, error)
	List(context.Context) ([]domain.MaintenanceStrategy, error)
}
type TemplateRepository interface {
	Create(context.Context, domain.InspectionTemplate) (domain.InspectionTemplate, error)
	Get(context.Context, string) (domain.InspectionTemplate, error)
	Update(context.Context, domain.InspectionTemplate) (domain.InspectionTemplate, error)
	List(context.Context, bool) ([]domain.InspectionTemplate, error)
}
type PlanRepository interface {
	Create(context.Context, domain.MaintenancePlan) (domain.MaintenancePlan, error)
	Get(context.Context, string) (domain.MaintenancePlan, error)
	Update(context.Context, domain.MaintenancePlan) (domain.MaintenancePlan, error)
	List(context.Context, domain.PlanFilter) ([]domain.MaintenancePlan, error)
}
type TaskRepository interface {
	Create(context.Context, domain.MaintenanceTask) (domain.MaintenanceTask, error)
	Get(context.Context, string) (domain.MaintenanceTask, error)
	Update(context.Context, domain.MaintenanceTask) (domain.MaintenanceTask, error)
	List(context.Context, domain.TaskFilter) ([]domain.MaintenanceTask, error)
}
type FindingRepository interface {
	Create(context.Context, domain.Finding) (domain.Finding, error)
	Get(context.Context, string) (domain.Finding, error)
	Update(context.Context, domain.Finding) (domain.Finding, error)
	List(context.Context, domain.FindingFilter) ([]domain.Finding, error)
}
type HistoryRepository interface {
	Append(context.Context, domain.HistoryEntry) error
	List(context.Context, string) ([]domain.HistoryEntry, error)
}
type Publisher interface {
	Publish(context.Context, domain.Event) error
}
type Clock func() interface{}
