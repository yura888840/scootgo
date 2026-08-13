package main

import (
	"context"
	"time"

	"github.com/yura888840/scootgo/internal/memlimit"
	"github.com/yura888840/scootgo/pkg/config"
	"github.com/yura888840/scootgo/pkg/logger"
	"github.com/yura888840/scootgo/pkg/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	memlimit.Start(ctx, 128, 5*time.Second) // included only with -tags memsim

	configs := config.GetAll()
	log := logger.Get(configs)

	server.Run(&log, configs)
}
