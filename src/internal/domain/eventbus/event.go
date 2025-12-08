package eventbus

type Event interface {
	EventName() string
}
