package eventbus

type Handler func(Event)

type Bus interface {
	Publish(events ...Event)
	Subscribe(eventName string, handler Handler)
}
