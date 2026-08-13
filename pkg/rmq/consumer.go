package rmq

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type consumer struct {
	channel         *amqp.Channel
	queue           *amqp.Queue
	dlQueue         *amqp.Queue
	batchSize       int
	connClosureChan <-chan *amqp.Error
	log             *zerolog.Logger
}

type deadletter struct {
	exchange   string
	routingkey string
	ttl        *int
}

func NewConsumer(
	c *AMQPConnector,
	queueName string,
	routingKey string,
	batchSize int,
	log *zerolog.Logger,
) (*consumer, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.Qos(batchSize, 0, false)
	if err != nil {
		ch.Close()
		return nil, err
	}

	dlq, err := createDLExchangeAndQueue(ch, c.dlExchange, c.exchange, queueName, routingKey)
	if err != nil {
		ch.Close()
		return nil, err
	}

	q, err := createAndBindQueue(ch, queueName, routingKey, c.exchange, deadletter{
		exchange:   c.dlExchange,
		routingkey: routingKey,
	})
	if err != nil {
		ch.Close()
		return nil, err
	}

	return &consumer{
		channel:         ch,
		queue:           q,
		dlQueue:         dlq,
		batchSize:       batchSize,
		connClosureChan: c.connClosureChan,
		log:             log,
	}, nil
}

func createDLExchangeAndQueue(ch *amqp.Channel, dlExchangeName, exchangeName, queueName, routingKey string) (*amqp.Queue, error) {
	err := ch.ExchangeDeclare(
		dlExchangeName, // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	delayQueueName := "delay." + queueName

	deadletterMessageTTL := 300000 // 5 minutes
	deadletterSetting := deadletter{
		exchange:   exchangeName, // set original exchange as deadletter exchange
		routingkey: routingKey,
		ttl:        &deadletterMessageTTL,
	}

	dlq, err := createAndBindQueue(ch, delayQueueName, routingKey, dlExchangeName, deadletterSetting)
	if err != nil {
		return nil, err
	}

	return dlq, nil
}

func createAndBindQueue(ch *amqp.Channel, queueName, routingKey, exchange string, dlSetting deadletter) (*amqp.Queue, error) {
	args := amqp.Table{amqp.QueueTypeArg: amqp.QueueTypeQuorum}

	args["x-dead-letter-exchange"] = dlSetting.exchange
	args["x-dead-letter-routing-key"] = dlSetting.routingkey
	if dlSetting.ttl != nil {
		args["x-message-ttl"] = *dlSetting.ttl
	}

	q, err := ch.QueueDeclare(queueName, true, false, false, false, args)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (c *consumer) Close() {
	if !c.channel.IsClosed() {
		err := c.channel.Close()
		if err != nil && !errors.Is(err, amqp.ErrClosed) {
			c.log.Error().Err(err).Msg("failed to close consumer channel")
		}
	}
}

func (c *consumer) Consume(
	ctx context.Context,
	errChan chan<- error,
	msgProcessor func(
		context.Context,
		[]RawMessage,
		chan<- []RawMessage,
		chan<- []RawMessage,
		chan<- []RawMessage,
	),
) {
	go c.monitorClosure(ctx, errChan)

	msgChannel, err := c.channel.ConsumeWithContext(
		ctx,
		c.queue.Name, // queue
		"",           // consumer
		false,        // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		c.log.Err(err).Msg("error while consuming from channel")
		errChan <- err
	}

	amqpMsgBatches := c.readMessageInBatch(msgChannel)
	msgBatches := c.createMessageBatches(amqpMsgBatches)
	c.processMessages(ctx, msgBatches, msgProcessor)
}

func (c *consumer) monitorClosure(ctx context.Context, errChan chan<- error) {
	chCloseChan := c.channel.NotifyClose(make(chan *amqp.Error))

	select {
	case <-ctx.Done():
		return
	case err := <-c.connClosureChan:
		c.log.Warn().Err(err).Msg("amqp connection closed")
		select {
		case errChan <- err:
		case <-ctx.Done():
		}
	case err := <-chCloseChan:
		c.log.Warn().Err(err).Msgf("amqp channel closed")
		select {
		case errChan <- err:
		case <-ctx.Done():
		}
	}
}

func (c *consumer) processMessages(
	ctx context.Context,
	msgBatches <-chan []RawMessage,
	msgProcessor func(
		context.Context,
		[]RawMessage,
		chan<- []RawMessage,
		chan<- []RawMessage,
		chan<- []RawMessage,
	),
) {
	successfulMsgs := make(chan []RawMessage, c.batchSize)
	failedMsgs := make(chan []RawMessage, c.batchSize)
	rejectedMsgs := make(chan []RawMessage, c.batchSize)

	errors := make(chan error, c.batchSize*3)

	go ackMessages(successfulMsgs, errors)
	go nackMessages(failedMsgs, errors)
	go rejectMessages(rejectedMsgs, errors)
	go c.logErrors(errors)

	go func() {
		for msgBatch := range msgBatches {
			c.log.Debug().Msgf("asynchronously processing %d messages", len(msgBatch))
			go msgProcessor(ctx, msgBatch, successfulMsgs, failedMsgs, rejectedMsgs)
		}
	}()
}

func ackMessages(groupedMsgs <-chan []RawMessage, errors chan<- error) {
	for msgs := range groupedMsgs {
		for index := range msgs {
			msg := msgs[index]
			err := msg.AMQPMessage.Ack(false)
			if err != nil {
				errors <- err
			}
		}
	}
}

func nackMessages(groupedMsgs <-chan []RawMessage, errors chan<- error) {
	for msgs := range groupedMsgs {
		for index := range msgs {
			msg := msgs[index]
			err := msg.AMQPMessage.Nack(false, false)
			if err != nil {
				errors <- err
			}
		}
	}
}

func rejectMessages(groupedMsgs <-chan []RawMessage, errors chan<- error) {
	for msgs := range groupedMsgs {
		for index := range msgs {
			msg := msgs[index]
			err := msg.AMQPMessage.Reject(false)
			if err != nil {
				errors <- err
			}
		}
	}
}

func (c *consumer) logErrors(errors <-chan error) {
	for err := range errors {
		c.log.Err(err).Msg("error while (n)acking message")
	}
}

func (c *consumer) createMessageBatches(amqpMsgBatches <-chan []amqp.Delivery) <-chan []RawMessage {
	msgBatchBuf := make(chan []RawMessage, c.batchSize)

	flushInterval := time.Second * 2
	flushTicker := time.NewTicker(flushInterval)

	go func() {
		for {
			select {
			case amqpMsgBatch := <-amqpMsgBatches:
				msgBatchBuf <- c.createMessages(amqpMsgBatch)
			case <-flushTicker.C:
				continue
			}
		}
	}()

	return msgBatchBuf
}

func (c *consumer) createMessages(amqpMsgs []amqp.Delivery) []RawMessage {
	messages := make([]RawMessage, len(amqpMsgs))
	for i := range amqpMsgs {
		amqpMsg := amqpMsgs[i]
		messages[i] = RawMessage{
			AMQPMessage: amqpMsg,
		}
	}

	return messages
}

func (c *consumer) readMessageInBatch(ch <-chan amqp.Delivery) <-chan []amqp.Delivery {
	amqpMsgBatch := make([]amqp.Delivery, 0, c.batchSize)
	amqpMsgBatchBuf := make(chan []amqp.Delivery, c.batchSize)

	flushInterval := time.Second * 30
	flushTicker := time.NewTicker(flushInterval)

	go func() {
		for {
			select {
			case amqpMsg, ok := <-ch:
				if !ok {
					return
				}
				amqpMsgBatch = append(amqpMsgBatch, amqpMsg)
				if len(amqpMsgBatch) == c.batchSize {
					amqpMsgBatchBuf <- amqpMsgBatch
					amqpMsgBatch = nil
				}
			case <-flushTicker.C:
				if len(amqpMsgBatch) > 0 {
					amqpMsgBatchBuf <- amqpMsgBatch
					amqpMsgBatch = nil
				}
			}
		}
	}()

	return amqpMsgBatchBuf
}
