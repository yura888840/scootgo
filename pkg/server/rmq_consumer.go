//go:build mysql

package server

import (
	"context"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	appconsumer "github.com/yura888840/scootgo/pkg/app/consumer"
	"github.com/yura888840/scootgo/pkg/rmq"
)

func setupRMQConsumer(
	ctx context.Context,
	errChannel chan<- error,
	log *zerolog.Logger,
	configs map[string]string,
	rootCmd *cobra.Command,
) {
	var consumeCmd = &cobra.Command{
		Use:   "consume",
		Short: "Start the AMQP consumer",
		Run: func(cmd *cobra.Command, args []string) {
			ConsumeMessages(ctx, errChannel, log, configs, args)
		},
	}

	rootCmd.AddCommand(consumeCmd)
}

func ConsumeMessages(
	ctx context.Context,
	errChannel chan<- error,
	log *zerolog.Logger,
	configs map[string]string,
	args []string,
) {
	if len(args) < 3 {
		log.Warn().Msg("usage: consume <queue> <routingKey> <batchSize>")
		return
	}

	queueName := args[0]
	routingKey := args[1]
	messageBatchSizeInput := args[2]

	messageBatchSize, err := strconv.Atoi(messageBatchSizeInput)
	if err != nil {
		log.Warn().Msgf("invalid consumer batch size %s", messageBatchSizeInput)

		messageBatchSize = 1
		log.Info().Msg("batch size is defaulted to 1")
	}

	amqpConnector, err := rmq.NewAMQPConnector(configs, log)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to RabbitMQ")
		errChannel <- err
		return
	}
	defer amqpConnector.Close()

	consumer, err := rmq.NewConsumer(amqpConnector, queueName, routingKey, messageBatchSize, log)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create consumer for RabbitMQ queue %s and routingKey %s", queueName, routingKey)
		errChannel <- err
		return
	}

	defer consumer.Close()

	messageHandler := getMessageHandler(routingKey, log)

	if messageHandler != nil {
		go consumer.Consume(ctx, errChannel, messageHandler)
		log.Info().Msgf("consumer setup successfully for queue %s and routingKey %s", queueName, routingKey)

		// Listen for the interrupt signal.
		<-ctx.Done()
		return
	}

	log.Warn().Msgf("consumer not configured for routingKey %s", routingKey)
}

func getMessageHandler(
	routingKey string,
	log *zerolog.Logger,
) func(
	context.Context,
	[]rmq.RawMessage,
	chan<- []rmq.RawMessage,
	chan<- []rmq.RawMessage,
	chan<- []rmq.RawMessage,
) {

	switch routingKey {
	case "sample.created":
		consumer := appconsumer.NewSampleConsumer(log)
		return consumer.Consume

	default:
		return nil
	}
}
