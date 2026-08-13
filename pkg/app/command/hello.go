package command

import (
	"context"

	"github.com/rs/zerolog"
)

type Hello struct {
	log *zerolog.Logger
}

func NewHello(log *zerolog.Logger) Hello {
	return Hello{log: log}
}

func (c Hello) Execute(ctx context.Context, args []string) int {
	if ctx.Err() != nil {
		c.log.Error().Err(ctx.Err()).Msg("hello command cancelled")
		return 1
	}

	message := "hello from the boilerplate command"
	if len(args) > 0 {
		message = args[0]
	}

	c.log.Info().Str("message", message).Msg("sample command executed")

	return 0
}
