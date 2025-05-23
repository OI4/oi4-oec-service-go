package api

type Subscription interface {
	GetID() string
	GetTopic() string
	GetQoS() byte
	GetHandler() MessageHandler
}
