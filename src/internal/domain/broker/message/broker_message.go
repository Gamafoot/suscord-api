package message

type BrokerMessage interface {
	EventName() string
}
