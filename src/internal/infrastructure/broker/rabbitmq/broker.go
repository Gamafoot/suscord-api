package rabbitmq

import (
	"context"
	"encoding/json"
	"suscord/internal/domain/broker/message"
	"suscord/internal/domain/logger"

	pkgErrors "github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type pooledChannel struct {
	channel  *amqp.Channel
	confirms chan amqp.Confirmation
}

type broker struct {
	conn        *amqp.Connection
	channelPool chan *pooledChannel
	logger      logger.Logger
}

func NewBroker(url string, channelPool int, logger logger.Logger) (*broker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	pool := make(chan *pooledChannel, channelPool)

	for i := 0; i < channelPool; i++ {
		ch, err := newChannel(conn)
		if err != nil {
			return nil, err
		}
		pool <- ch
	}

	ch := <-pool
	if err = exchangeDeclare(ch.channel); err != nil {
		return nil, err
	}
	pool <- ch

	return &broker{
		conn:        conn,
		channelPool: pool,
		logger:      logger,
	}, nil
}

type messageBody struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func newChannel(conn *amqp.Connection) (*pooledChannel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err = ch.Confirm(false); err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	return &pooledChannel{
		channel:  ch,
		confirms: confirms,
	}, nil
}

func (p *broker) Publish(
	ctx context.Context,
	message message.BrokerMessage,
) error {
	select {
	case ch := <-p.channelPool:
		defer func() {
			if !ch.channel.IsClosed() {
				p.channelPool <- ch
			} else {
				ch, err := newChannel(p.conn)
				if err != nil {
					p.logger.Err(err)
					return
				}
				p.channelPool <- ch
			}
		}()

		body, err := json.Marshal(messageBody{
			Type:    message.EventName(),
			Payload: message,
		})
		if err != nil {
			return err
		}

		err = ch.channel.PublishWithContext(
			ctx,
			"chat.events",
			"chat.ws",
			true,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         body,
			},
		)
		if err != nil {
			return pkgErrors.WithStack(err)
		}

		for {
			select {
			case confirm := <-ch.confirms:
				if !confirm.Ack {
					return pkgErrors.New("Broker.Publish не подтвердил доставку сообщения")
				}
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}

	case <-ctx.Done():
		return ctx.Err()
	}
}

func exchangeDeclare(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		"chat.events",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return pkgErrors.WithStack(err)
	}

	return nil
}
