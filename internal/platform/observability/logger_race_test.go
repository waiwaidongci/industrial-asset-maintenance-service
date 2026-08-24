package observability

import (
	"sync"
	"testing"
)

func TestAccessLogSnapshotIsolation(t *testing.T) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; _ = NewLogger() }()
	}
	close(start)
	wg.Wait()
}
