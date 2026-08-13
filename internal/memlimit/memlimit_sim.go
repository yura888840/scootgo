//go:build memsim

package memlimit

import (
	"context"
	"log"
	"os"
	"runtime"
	"time"
)

func Start(ctx context.Context, limitMB uint64, every time.Duration) {
	t := time.NewTicker(every)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				usedMB := m.Alloc / 1024 / 1024
				if usedMB > limitMB {
					log.Printf("[memsim] Memory limit exceeded: %d MB > %d MB", usedMB, limitMB)
					os.Exit(1)
				}
			}
		}
	}()
}
