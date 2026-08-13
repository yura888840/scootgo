package tracing

import (
	"runtime"
	"time"

	"github.com/rs/zerolog"
)

func TrackUsage(start time.Time, name string, logger *zerolog.Logger) {
	elapsed := time.Since(start)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	logger.Info().Msgf("%s took %s", name, elapsed)
	logger.Info().Msgf("Alloc = %v MiB", bToMb(m.Alloc))
	logger.Info().Msgf("TotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	logger.Info().Msgf("Sys = %v MiB", bToMb(m.Sys))
	logger.Info().Msgf("NumGC = %v", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
