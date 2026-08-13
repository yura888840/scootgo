package consumer

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/yura888840/scootgo/pkg/rmq"
)

type SampleConsumer struct {
	log *zerolog.Logger
}

func NewSampleConsumer(log *zerolog.Logger) SampleConsumer {
	return SampleConsumer{log: log}
}

func (c SampleConsumer) Consume(
	ctx context.Context,
	rawMessages []rmq.RawMessage,
	successful chan<- []rmq.RawMessage,
	failed chan<- []rmq.RawMessage,
	rejected chan<- []rmq.RawMessage,
) {
	if ctx.Err() != nil {
		failed <- rawMessages
		return
	}

	for index := range rawMessages {
		rawMessage := rawMessages[index]

		var message Message
		if err := json.Unmarshal(rawMessage.AMQPMessage.Body, &message); err != nil {
			c.log.Warn().Err(err).Msg("rejecting invalid sample message")
			rejected <- []rmq.RawMessage{rawMessage}
			continue
		}

		if message.ID == "" {
			c.log.Warn().Msg("rejecting sample message without id")
			rejected <- []rmq.RawMessage{rawMessage}
			continue
		}

		c.log.Info().Str("id", message.ID).Str("message", message.Message).Msg("sample consumer processed message")
		successful <- []rmq.RawMessage{rawMessage}
	}
}
