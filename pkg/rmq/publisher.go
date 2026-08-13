package rmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

// AMQPPublisher interface
type AMQPPublisher interface {
	PublishJSON(context.Context, string, any, time.Time) error
	Close()
}

// Publisher struct implementing AMQPPublisher
type Publisher struct {
	channel  *amqp.Channel
	exchange string
	log      *zerolog.Logger
}

func NewPublisher(conn *AMQPConnector, log *zerolog.Logger) (AMQPPublisher, error) {
	ch, err := conn.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return &Publisher{
		channel:  ch,
		exchange: conn.exchange,
		log:      log,
	}, nil
}

func (p *Publisher) PublishJSON(
	ctx context.Context,
	routingKey string,
	message any,
	currentTimestamp time.Time,
) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	jsonMessage, err := prepJSONMessage(message, currentTimestamp)
	if err != nil {
		return err
	}

	err = p.channel.PublishWithContext(ctx, p.exchange, routingKey, false, false, jsonMessage)
	if err != nil {
		return fmt.Errorf("failed to publish message %+v: %w", message, err)
	}

	return nil
}

func (p *Publisher) Close() {
	if !p.channel.IsClosed() {
		err := p.channel.Close()
		if err != nil && !errors.Is(err, amqp.ErrClosed) {
			p.log.Error().Err(err).Msg("failed to close publisher channel")
		}
	}
}

func prepJSONMessage(
	message any,
	messageTimestamp time.Time,
) (amqp.Publishing, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return amqp.Publishing{}, fmt.Errorf("failed to marshal message to JSON: %w", err)
	}

	return amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   messageTimestamp,
	}, nil
}
