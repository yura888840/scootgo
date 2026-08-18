//go:build mysql

package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	apphandler "github.com/yura888840/scootgo/pkg/app/handler"
	"github.com/yura888840/scootgo/pkg/config"
	"github.com/yura888840/scootgo/pkg/healthcheck"
	"github.com/yura888840/scootgo/pkg/ride"
	"github.com/yura888840/scootgo/pkg/tracing"
)

func setupHTTPServer(
	ctx context.Context,
	errChan chan<- error,
	log *zerolog.Logger,
	configs map[string]string,
	rootCmd *cobra.Command,
) {
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Serve HTTP request",
		Run: func(cmd *cobra.Command, args []string) {
			serve(ctx, errChan, log, configs)
		},
	}

	rootCmd.AddCommand(serveCmd)
}

func serve(
	ctx context.Context,
	errChan chan<- error,
	log *zerolog.Logger,
	configs map[string]string,
) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	serverMode := configs[config.ConfigServerMode]
	gin.SetMode(serverMode)

	router := gin.New()
	router.Use(gin.Recovery())

	healthcheck.Routes(router, log, func(context.Context) error { return nil }, func() bool { return true }, errChan)

	router.Use(tracing.SetCorrelationID())
	router.Use(tracing.LogRequest(log))
	router.Use(gzip.Gzip(gzip.DefaultCompression))

	apphandler.NewRouter().Routes(router, configs)
	ride.CreateNewRouter(log).SetupRoutes(router)

	port := configs[config.ConfigAppPort]
	address := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 60 * time.Second,
	}

	go func() {
		log.Info().Str("address", address).Msg("HTTP server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	cancel()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel = context.WithTimeout(context.Background(), gracefulShutdownTimout)
	defer cancel()
	shutdownHTTPServer(ctx, log, srv)
}

func shutdownHTTPServer(ctx context.Context, log *zerolog.Logger, srv *http.Server) {
	if err := srv.Shutdown(ctx); err != nil {
		log.Info().Msgf("Server forced to shutdown: %v", err)
		return
	}

	log.Info().Msg("Server exiting")
}
