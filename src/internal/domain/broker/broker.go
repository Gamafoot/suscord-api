package broker

import (
	"context"
	"suscord/internal/domain/broker/event"
)

type Broker interface {
	Publish(ctx context.Context, event event.Event) error
}
