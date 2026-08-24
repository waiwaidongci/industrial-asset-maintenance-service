package config

import (
	"os"
	"testing"
	"time"
)

func TestShutdownTimeoutFromYAML(t *testing.T) {
	c := Config{ShutdownTimeout: 10 * time.Second, RequestTimeout: 15 * time.Second}
	parseYAML(&c, "shutdown_timeout: 3s")
	if c.ShutdownTimeout != 3*time.Second {
		t.Fatalf("ShutdownTimeout = %s", c.ShutdownTimeout)
	}
}

func TestShutdownTimeoutFromEnv(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT", "4s")
	c := Config{ShutdownTimeout: 10 * time.Second, RequestTimeout: 15 * time.Second}
	applyEnv(&c)
	if c.ShutdownTimeout != 4*time.Second {
		t.Fatalf("ShutdownTimeout = %s", c.ShutdownTimeout)
	}
	_ = os.Getenv("SHUTDOWN_TIMEOUT")
}

func TestRejectsZeroShutdownTimeout(t *testing.T) {
	c := Config{ShutdownTimeout: 10 * time.Second, RequestTimeout: 15 * time.Second}
	parseYAML(&c, "shutdown_timeout: 0s")
	if c.ShutdownTimeout != 10*time.Second || c.RequestTimeout != 15*time.Second {
		t.Fatalf("zero shutdown timeout changed configuration: %+v", c)
	}
}

func TestCapsShutdownTimeout(t *testing.T) {
	c := Config{ShutdownTimeout: 10 * time.Second}
	parseYAML(&c, "shutdown_timeout: 48h")
	if c.ShutdownTimeout != 24*time.Hour {
		t.Fatalf("ShutdownTimeout = %s", c.ShutdownTimeout)
	}
}
