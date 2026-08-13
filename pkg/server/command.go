//go:build mysql

package server

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	appcommand "github.com/yura888840/scootgo/pkg/app/command"
)

// Command executor
func triggerCommand(
	ctx context.Context,
	log *zerolog.Logger,
	rootCmd *cobra.Command,
) {
	var triggerCmd = &cobra.Command{
		Use:   "exec",
		Short: "Trigger given command",
		Run: func(cmd *cobra.Command, args []string) {
			trigger(ctx, log, args)
		},
	}

	rootCmd.AddCommand(triggerCmd)
}

func trigger(ctx context.Context, log *zerolog.Logger, args []string) {
	if len(args) == 0 {
		log.Warn().Msg("no command specified")
		os.Exit(1)
	}

	cmdName := args[0]
	cmdArgs := args[1:]

	switch cmdName {
	case "hello":
		command := appcommand.NewHello(log)
		os.Exit(command.Execute(ctx, cmdArgs))
	default:
		log.Warn().Msgf("command %s is not supported", cmdName)
		os.Exit(1)
	}
}
