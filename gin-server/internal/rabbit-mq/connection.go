package rabbitmq

import (
	"fmt"
	"gin-server/internal/config"
	"gin-server/internal/modules"

	"github.com/streadway/amqp"
)

func Connect() (*modules.RabbitMQ, error) {
	cgf := config.Load()

	// connect to rabbitMQ
	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%s", cgf.Rabbit_User, cgf.Rabbit_Password, cgf.Rabbit_Host, cgf.Rabbit_Port)
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	// create service client

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// create queue
	q, err := ch.QueueDeclare(
		"music_upload",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &modules.RabbitMQ{
		Conn: conn,
		Ch:   ch,
		Q:    q,
	}, nil
}
