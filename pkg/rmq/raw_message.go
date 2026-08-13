package rmq

import amqp "github.com/rabbitmq/amqp091-go"

type RawMessage struct {
	AMQPMessage amqp.Delivery
}
