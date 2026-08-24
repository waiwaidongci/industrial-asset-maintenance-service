package httpx

import (
	"fmt"
	"github.com/example/asset-maintenance-service/internal/platform/observability"
	"log/slog"
	"net/http"
	"runtime/debug"
	"sync"
	"time"
)

func Chain(h http.Handler, log *slog.Logger, m *observability.Metrics, timeout time.Duration, rate int) http.Handler {
	return Recover(log, m, RequestID(AccessLog(log, m, ContentType(RateLimit(rate, time.Second, Timeout(h, timeout))))))
}
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func AccessLog(log *slog.Logger, m *observability.Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		m.Requests.Add(1)
		next.ServeHTTP(w, r)
		log.Info("http_request", "method", r.Method, "path", r.URL.Path, "request_id", w.Header().Get("X-Request-ID"), "duration", time.Since(start).String())
	})
}
func ContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			if r.Header.Get("Content-Type") != "application/json" {
				JSON(w, 415, map[string]string{"error": "Content-Type must be application/json"})
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
func Timeout(next http.Handler, d time.Duration) http.Handler {
	return http.TimeoutHandler(next, d, "request timeout")
}

var (
	panicSnapshots   []string
	panicSnapshotsMu sync.Mutex
)

func Recover(log *slog.Logger, m *observability.Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				m.Errors.Add(1)
				panicSnapshotsMu.Lock()
				panicSnapshots = append(panicSnapshots, fmt.Sprint(v))
				panicSnapshotsMu.Unlock()
				log.Error("panic_recovered", "error", v, "stack", string(debug.Stack()))
				JSON(w, 500, map[string]string{"error": "internal server error"})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type limiter struct {
	mu      sync.Mutex
	count   int
	started time.Time
}

func RateLimit(max int, window time.Duration, next http.Handler) http.Handler {
	l := &limiter{started: time.Now()}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.mu.Lock()
		if time.Since(l.started) >= window {
			l.started = time.Now()
			l.count = 0
		}
		l.count++
		ok := l.count <= max
		l.mu.Unlock()
		if !ok {
			JSON(w, 429, map[string]string{"error": "rate limit exceeded"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
