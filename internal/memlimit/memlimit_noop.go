//go:build !memsim

package memlimit

import (
	"context"
	"time"
)

func Start(_ context.Context, _ uint64, _ time.Duration) { /* no-op */ }
