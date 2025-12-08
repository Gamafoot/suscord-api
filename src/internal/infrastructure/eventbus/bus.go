package eventbus

import "suscord/internal/domain/eventbus"

type bus struct {
	handlers map[string][]eventbus.Handler
}

func NewBus() eventbus.Bus {
	return &bus{
		handlers: make(map[string][]eventbus.Handler),
	}
}

func (b *bus) Publish(events ...eventbus.Event) {
	for _, evt := range events {
		if hs, ok := b.handlers[evt.EventName()]; ok {
			for _, h := range hs {
				go h(evt)
			}
		}
	}
}

func (b *bus) Subscribe(eventName string, handler eventbus.Handler) {
	b.handlers[eventName] = append(b.handlers[eventName], handler)
}
