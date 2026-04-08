package modules

import "github.com/streadway/amqp"

type Rabbit struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
	Q    amqp.Queue
}
