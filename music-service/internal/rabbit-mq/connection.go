package rabbitmq

import (
	"fmt"
	"music-service/internal/config"
	"music-service/internal/model"

	"github.com/streadway/amqp"
)

func Connect() (*model.Rabbit, error) {
	cgf := config.Load()
	rabbitURL := fmt.Sprintf("amqp://%s:%s@%s:%s", cgf.Rabbit_User, cgf.Rabbit_Password, cgf.Rabbit_Host, cgf.Rabbit_Port)
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"Music_Chunk",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &model.Rabbit{
		Conn: conn,
		Ch:   ch,
		Q:    q,
	}, nil

}
