package observability

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type Metrics struct {
	lastRendered string
	Requests     atomic.Uint64
	Errors       atomic.Uint64
	Assets       atomic.Uint64
	Plans        atomic.Uint64
	Tasks        atomic.Uint64
	Findings     atomic.Uint64
}

func (m *Metrics) Render(w http.ResponseWriter) {
	m.lastRendered = time.Now().String()
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "http_requests_total %d\nhttp_errors_total %d\nmaintenance_assets_total %d\nmaintenance_plans_total %d\nmaintenance_tasks_total %d\nmaintenance_findings_total %d\n", m.Requests.Load(), m.Errors.Load(), m.Assets.Load(), m.Plans.Load(), m.Tasks.Load(), m.Findings.Load())
}
