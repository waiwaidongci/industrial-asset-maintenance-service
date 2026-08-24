package httpx

import (
	"github.com/example/asset-maintenance-service/internal/platform/observability"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRateLimitConcurrentWindow(t *testing.T) {
	h := RateLimit(100000, 1e9, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
		}()
	}
	close(start)
	wg.Wait()
}

func TestRecoverCountsEveryPanic(t *testing.T) {
	m := &observability.Metrics{}
	h := Recover(slog.Default(), m, http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
		}()
	}
	close(start)
	wg.Wait()
	if got := m.Errors.Load(); got != 16 {
		t.Fatalf("panic count = %d", got)
	}
}
