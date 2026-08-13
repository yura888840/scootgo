package consumer

import (
	"bytes"
	"context"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yura888840/scootgo/pkg/rmq"
)

func TestSampleConsumerConsume(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		successCount  int
		rejectedCount int
		failedCount   int
		logContains   string
	}{
		{
			name:         "valid message acknowledged",
			body:         `{"id":"123","message":"hello"}`,
			successCount: 1,
			logContains:  "sample consumer processed message",
		},
		{
			name:          "invalid json rejected",
			body:          `{`,
			rejectedCount: 1,
			logContains:   "rejecting invalid sample message",
		},
		{
			name:          "missing id rejected",
			body:          `{"message":"hello"}`,
			rejectedCount: 1,
			logContains:   "rejecting sample message without id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := zerolog.New(&buf)
			consumer := NewSampleConsumer(&logger)

			successful := make(chan []rmq.RawMessage, 1)
			failed := make(chan []rmq.RawMessage, 1)
			rejected := make(chan []rmq.RawMessage, 1)

			consumer.Consume(
				context.Background(),
				[]rmq.RawMessage{{AMQPMessage: amqp.Delivery{Body: []byte(tt.body)}}},
				successful,
				failed,
				rejected,
			)

			require.Len(t, successful, tt.successCount)
			require.Len(t, failed, tt.failedCount)
			require.Len(t, rejected, tt.rejectedCount)
			assert.Contains(t, buf.String(), tt.logContains)
		})
	}
}
