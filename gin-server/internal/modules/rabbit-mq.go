package modules

import "github.com/streadway/amqp"

type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	Q    amqp.Queue
}
