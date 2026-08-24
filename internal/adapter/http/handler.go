package httpadapter

import (
	"github.com/example/asset-maintenance-service/internal/application"
	"github.com/example/asset-maintenance-service/internal/domain"
	"github.com/example/asset-maintenance-service/internal/platform/httpx"
	"github.com/example/asset-maintenance-service/internal/platform/observability"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Svc     *application.Service
	Metrics *observability.Metrics
}

func NewHandler(s *application.Service, m *observability.Metrics) *Handler {
	return &Handler{Svc: s, Metrics: m}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(p, "/")
	if p == "healthz" {
		httpx.JSON(w, 200, map[string]string{"status": "ok"})
		return
	}
	if p == "readyz" {
		httpx.JSON(w, 200, map[string]string{"status": "ready"})
		return
	}
	if p == "metrics" {
		h.Metrics.Render(w)
		return
	}
	if p == "summary" && r.Method == http.MethodGet {
		v, err := h.Svc.Summary(r.Context())
		if err != nil {
			httpx.Error(w, err)
			return
		}
		httpx.JSON(w, http.StatusOK, v)
		return
	}
	if len(parts) == 0 {
		return
	}
	switch parts[0] {
	case "assets":
		h.assets(w, r, parts)
	case "strategies":
		h.strategies(w, r, parts)
	case "templates":
		h.templates(w, r, parts)
	case "plans":
		h.plans(w, r, parts)
	case "tasks":
		h.tasks(w, r, parts)
	case "findings":
		h.findings(w, r, parts)
	default:
		http.NotFound(w, r)
	}
}
func (h *Handler) assets(w http.ResponseWriter, r *http.Request, p []string) {
	if r.Method == http.MethodPost && len(p) == 1 {
		var v domain.Asset
		if e := httpx.Decode(r, &v); e != nil {
			httpx.Error(w, e)
			return
		}
		v, e := h.Svc.CreateAsset(r.Context(), v)
		if e != nil {
			httpx.Error(w, e)
			return
		}
		h.Metrics.Assets.Add(1)
		httpx.JSON(w, 201, v)
		return
	}
	if r.Method == http.MethodGet && len(p) == 1 {
		v, e := h.Svc.ListAssets(r.Context(), domain.AssetFilter{Status: domain.AssetStatus(r.URL.Query().Get("status")), AssetType: r.URL.Query().Get("asset_type"), Location: r.URL.Query().Get("location"), Search: r.URL.Query().Get("search")})
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	if len(p) >= 2 {
		v, e := h.Svc.GetAsset(r.Context(), p[1])
		if e != nil {
			httpx.Error(w, e)
			return
		}
		if r.Method == http.MethodGet {
			httpx.JSON(w, 200, v)
			return
		}
		if r.Method == http.MethodPatch {
			var q struct {
				Name     string             `json:"name"`
				Location string             `json:"location"`
				Status   domain.AssetStatus `json:"status"`
				Metadata map[string]string  `json:"metadata"`
			}
			if e = httpx.Decode(r, &q); e != nil {
				httpx.Error(w, e)
				return
			}
			if q.Name != "" {
				v.Name = q.Name
			}
			if q.Location != "" {
				v.Location = q.Location
			}
			if q.Status != "" {
				v.Status = q.Status
			}
			if q.Metadata != nil {
				v.Metadata = q.Metadata
			}
			v, e = h.Svc.UpdateAsset(r.Context(), v)
			if e != nil {
				httpx.Error(w, e)
				return
			}
			httpx.JSON(w, 200, v)
			return
		}
	}
	http.NotFound(w, r)
}
func (h *Handler) strategies(w http.ResponseWriter, r *http.Request, p []string) {
	if r.Method == http.MethodGet {
		v, e := h.Svc.ListStrategies(r.Context())
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	if r.Method == http.MethodPost {
		var v domain.MaintenanceStrategy
		if e := httpx.Decode(r, &v); e != nil {
			httpx.Error(w, e)
			return
		}
		v, e := h.Svc.CreateStrategy(r.Context(), v)
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 201, v)
		return
	}
	http.NotFound(w, r)
}
func (h *Handler) templates(w http.ResponseWriter, r *http.Request, p []string) {
	if r.Method == http.MethodGet {
		v, e := h.Svc.ListTemplates(r.Context(), r.URL.Query().Get("active") == "true")
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	if r.Method == http.MethodPost {
		var v domain.InspectionTemplate
		if e := httpx.Decode(r, &v); e != nil {
			httpx.Error(w, e)
			return
		}
		v, e := h.Svc.CreateTemplate(r.Context(), v)
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 201, v)
		return
	}
	if len(p) >= 2 && r.Method == http.MethodGet {
		v, e := h.Svc.GetTemplate(r.Context(), p[1])
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	http.NotFound(w, r)
}
func (h *Handler) plans(w http.ResponseWriter, r *http.Request, p []string) {
	if r.Method == http.MethodGet && len(p) == 1 {
		v, e := h.Svc.ListPlans(r.Context(), domain.PlanFilter{AssetID: r.URL.Query().Get("asset_id"), Status: domain.PlanStatus(r.URL.Query().Get("status"))})
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	if r.Method == http.MethodPost && len(p) == 1 {
		var v domain.MaintenancePlan
		if e := httpx.Decode(r, &v); e != nil {
			httpx.Error(w, e)
			return
		}
		v, e := h.Svc.CreatePlan(r.Context(), v)
		if e != nil {
			httpx.Error(w, e)
			return
		}
		h.Metrics.Plans.Add(1)
		httpx.JSON(w, 201, v)
		return
	}
	if len(p) >= 2 {
		v, e := h.Svc.GetPlan(r.Context(), p[1])
		if e != nil {
			httpx.Error(w, e)
			return
		}
		if r.Method == http.MethodGet && len(p) == 2 {
			httpx.JSON(w, 200, v)
			return
		}
		if r.Method == http.MethodPost && len(p) == 3 && p[2] == "status" {
			var q struct {
				Status domain.PlanStatus `json:"status"`
			}
			if e = httpx.Decode(r, &q); e != nil {
				httpx.Error(w, e)
				return
			}
			v, e = h.Svc.ChangePlanStatus(r.Context(), v.ID, q.Status)
			if e != nil {
				httpx.Error(w, e)
				return
			}
			httpx.JSON(w, 200, v)
			return
		}
		if r.Method == http.MethodPost && len(p) == 3 && p[2] == "tasks" {
			scheduled := time.Now()
			if s := r.URL.Query().Get("scheduled_at"); s != "" {
				scheduled, _ = time.Parse(time.RFC3339, s)
			}
			t, e := h.Svc.CreateTask(r.Context(), v.ID, scheduled)
			if e != nil {
				httpx.Error(w, e)
				return
			}
			h.Metrics.Tasks.Add(1)
			httpx.JSON(w, 201, t)
			return
		}
	}
	http.NotFound(w, r)
}
func (h *Handler) tasks(w http.ResponseWriter, r *http.Request, p []string) {
	if r.Method == http.MethodGet && len(p) == 1 {
		v, e := h.Svc.ListTasks(r.Context(), domain.TaskFilter{AssetID: r.URL.Query().Get("asset_id"), PlanID: r.URL.Query().Get("plan_id"), Status: domain.TaskStatus(r.URL.Query().Get("status")), Assignee: r.URL.Query().Get("assignee")})
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	if len(p) >= 2 {
		v, e := h.Svc.GetTask(r.Context(), p[1])
		if e != nil {
			httpx.Error(w, e)
			return
		}
		if r.Method == http.MethodGet && len(p) == 2 {
			httpx.JSON(w, 200, v)
			return
		}
		if r.Method == http.MethodPost && len(p) == 3 && p[2] == "status" {
			var q struct {
				Status domain.TaskStatus `json:"status"`
				Actor  string            `json:"actor"`
				Notes  string            `json:"notes"`
			}
			if e = httpx.Decode(r, &q); e != nil {
				httpx.Error(w, e)
				return
			}
			v, e = h.Svc.ChangeTaskStatus(r.Context(), v.ID, q.Status, q.Actor, q.Notes)
			if e != nil {
				httpx.Error(w, e)
				return
			}
			httpx.JSON(w, 200, v)
			return
		}
		if r.Method == http.MethodGet && len(p) == 3 && p[2] == "history" {
			v, e := h.Svc.TaskHistory(r.Context(), v.ID)
			if e != nil {
				httpx.Error(w, e)
				return
			}
			httpx.JSON(w, 200, v)
			return
		}
	}
	http.NotFound(w, r)
}
func (h *Handler) findings(w http.ResponseWriter, r *http.Request, p []string) {
	if r.Method == http.MethodGet && len(p) == 1 {
		v, e := h.Svc.ListFindings(r.Context(), domain.FindingFilter{AssetID: r.URL.Query().Get("asset_id"), TaskID: r.URL.Query().Get("task_id"), Status: domain.FindingStatus(r.URL.Query().Get("status")), Severity: domain.FindingSeverity(r.URL.Query().Get("severity"))})
		if e != nil {
			httpx.Error(w, e)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	if r.Method == http.MethodPost && len(p) == 1 {
		var v domain.Finding
		if e := httpx.Decode(r, &v); e != nil {
			httpx.Error(w, e)
			return
		}
		v, e := h.Svc.ReportFinding(r.Context(), v)
		if e != nil {
			httpx.Error(w, e)
			return
		}
		h.Metrics.Findings.Add(1)
		httpx.JSON(w, 201, v)
		return
	}
	if len(p) >= 2 {
		v, e := h.Svc.Findings.Get(r.Context(), p[1])
		if e != nil {
			httpx.Error(w, e)
			return
		}
		if r.Method == http.MethodPost && len(p) == 3 && p[2] == "status" {
			var q struct {
				Status domain.FindingStatus `json:"status"`
			}
			if e = httpx.Decode(r, &q); e != nil {
				httpx.Error(w, e)
				return
			}
			v, e = h.Svc.ChangeFindingStatus(r.Context(), v.ID, q.Status)
			if e != nil {
				httpx.Error(w, e)
				return
			}
			httpx.JSON(w, 200, v)
			return
		}
		httpx.JSON(w, 200, v)
		return
	}
	http.NotFound(w, r)
}
func _unused(s string) int { v, _ := strconv.Atoi(s); return v }
