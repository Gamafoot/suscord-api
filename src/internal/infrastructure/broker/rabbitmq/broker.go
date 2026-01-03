package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"suscord/internal/domain/broker/event"

	pkgErrors "github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type broker struct {
	conn        *amqp.Connection
	channelPool chan *amqp.Channel
}

func NewBroker(url string, channelPool int) (*broker, error) {

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	pool := make(chan *amqp.Channel, channelPool)

	for i := 0; i < channelPool; i++ {
		ch, err := conn.Channel()
		if err != nil {
			return nil, err
		}
		pool <- ch
	}

	return &broker{
		conn:        conn,
		channelPool: pool,
	}, nil
}

func (p *broker) Publish(
	ctx context.Context,
	event event.Event,
) error {
	select {
	case ch := <-p.channelPool:
		defer func() {
			p.channelPool <- ch
		}()

		q, err := ch.QueueDeclare(
			getExchange(event.EventName()),
			true,
			true,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Fatal(err)
		}

		body, err := json.Marshal(event)
		if err != nil {
			return err
		}

		return ch.PublishWithContext(
			ctx,
			q.Name,
			event.EventName(),
			true,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)

	case <-ctx.Done():
		return ctx.Err()
	}
}

func getExchange(eventName string) string {
	items := strings.Split(eventName, ".")
	return strings.Join(items[:len(items)-1], ".") + "*"
}
