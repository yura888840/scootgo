package rmq

import (
	"errors"
	"fmt"
	"net/url"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"github.com/yura888840/scootgo/pkg/config"
)

type AMQPConnector struct {
	conn            *amqp.Connection
	exchange        string
	dlExchange      string
	connClosureChan <-chan *amqp.Error
	log             *zerolog.Logger
}

func (a *AMQPConnector) IsConnected() bool {
	return a.conn != nil && !a.conn.IsClosed()
}

func NewAMQPConnector(
	envVars map[string]string,
	log *zerolog.Logger,
) (*AMQPConnector, error) {
	connString := createConnString(envVars)
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}

	return &AMQPConnector{
		conn:            conn,
		exchange:        "amq.topic",
		dlExchange:      "dlx.dead-letter",
		connClosureChan: conn.NotifyClose(make(chan *amqp.Error)),
		log:             log,
	}, nil
}

func (a *AMQPConnector) Channel() (*amqp.Channel, error) {
	return a.conn.Channel()
}

func (a *AMQPConnector) Close() {
	err := a.conn.Close()
	if err != nil && !errors.Is(err, amqp.ErrClosed) {
		a.log.Warn().Err(err).Msg("failed to close amqp connection")
	}
}

func createConnString(envVars map[string]string) string {
	userName := envVars[config.ConfigRMQUser]
	password := envVars[config.ConfigRMQPassword]
	host := envVars[config.ConfigRMQHost]
	virtHost := envVars[config.ConfigRMQVirtualHost]
	port := envVars[config.ConfigRMQPort]

	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s", userName, url.PathEscape(password), host, port, virtHost)
}
