package broker

import (
	"context"
	"suscord/internal/domain/broker/message"
)

type Broker interface {
	Publish(ctx context.Context, message message.BrokerMessage) error
}
