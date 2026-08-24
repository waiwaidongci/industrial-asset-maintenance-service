package main

import (
	"github.com/example/asset-maintenance-service/internal/platform/config"
	"testing"
	"time"
)

func TestServerUsesConfiguredShutdownTimeout(t *testing.T) {
	c := config.Config{ShutdownTimeout: 2 * time.Second, RequestTimeout: 15 * time.Second}
	if got := shutdownTimeout(c); got != 2*time.Second {
		t.Fatalf("shutdown timeout = %s", got)
	}
}
