package httpadapter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/asset-maintenance-service/internal/application"
	"github.com/example/asset-maintenance-service/internal/infrastructure/memory"
	"github.com/example/asset-maintenance-service/internal/platform/observability"
)

func TestMissingAssetHTTPStatus(t *testing.T) {
	svc := &application.Service{Assets: memory.NewAssetRepository(), Now: time.Now}
	handler := NewHandler(svc, &observability.Metrics{})
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/assets/missing", nil))
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}
