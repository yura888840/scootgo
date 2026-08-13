package logger

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yura888840/scootgo/pkg/config"
	"github.com/yura888840/scootgo/pkg/tracing"
)

func Get(configs map[string]string) zerolog.Logger {
	logLevel := getLogLevel(configs)

	appEnv := configs[config.ConfigAppEnv]
	if appEnv == "dev" {
		return setupHumanReadableLogs(logLevel)
	}

	return setupMachineReadableLogs(logLevel)
}

func getLogLevel(configs map[string]string) int {
	logLevel, err := strconv.Atoi(configs[config.ConfigLogLevel])
	if err != nil {
		return int(zerolog.InfoLevel)
	}

	return logLevel
}

func setupHumanReadableLogs(logLevel int) zerolog.Logger {
	return getBaseLogger().
		Level(zerolog.Level(logLevel)).
		Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func setupMachineReadableLogs(logLevel int) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.LevelFieldName = "severity"
	return getBaseLogger().
		Level(zerolog.Level(logLevel))
}

func getBaseLogger() zerolog.Logger {
	return log.With().
		Caller().
		Logger().Hook(tracing.TracingHook{})
}
