//go:build mysql

package server

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const gracefulShutdownTimout = 5 * time.Second

func Run(log *zerolog.Logger, configs map[string]string) {
	signalCtx, stopSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stopSignal()

	// Cancellable context for error handling
	ctx, cancel := context.WithCancel(signalCtx)
	defer cancel()

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A boilerplate service with HTTP, command, and consumer entry points",
	}

	errChan := make(chan error, 2)
	defer close(errChan)

	go func() {
		for err := range errChan {
			log.Warn().Err(err).Msg("system error received")
			cancel()
		}
	}()

	setupRMQConsumer(ctx, errChan, log, configs, rootCmd)
	setupHTTPServer(ctx, errChan, log, configs, rootCmd)
	triggerCommand(ctx, log, rootCmd)

	go func() {
		var err error
		err = rootCmd.Execute()
		if err != nil {
			log.Fatal().Err(err).Msg("error while executing command")
			cancel()
		}
	}()

	// unblock on either signal OR error
	<-ctx.Done()
	log.Info().Msg("shutting down gracefully, press Ctrl+C again to force")
}
