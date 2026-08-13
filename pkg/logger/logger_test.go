package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rs/zerolog"
	"github.com/yura888840/scootgo/pkg/config"
)

var tests = []struct {
	name        string
	input       map[string]string
	outputLevel zerolog.Level
}{
	{
		"default level",
		map[string]string{
			config.ConfigAppEnv: "dev",
		},
		zerolog.InfoLevel,
	},
	{
		"passed level",
		map[string]string{
			config.ConfigAppEnv:   "prod",
			config.ConfigLogLevel: "0",
		},
		zerolog.DebugLevel,
	},
}

func TestGet(t *testing.T) {
	for _, test := range tests {
		logger := Get(test.input)
		assert.Equal(t, test.outputLevel, logger.GetLevel())
	}
}
