package subscription

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/google/uuid"
)

type Impl struct {
	api.Subscription
	id       string
	topic    string
	interval uint32
	config   api.SubscriptionConfig
	qos      byte
	handler  api.MessageHandler
}

func NewTopicSubscription(topic string, handler api.MessageHandler, opts ...func(*Impl)) *Impl {
	subscription := &Impl{
		id:       uuid.New().String(),
		topic:    topic,
		interval: 0,
		config:   api.SubsciptionConfig_CONF_1,
		qos:      1,
		handler:  handler,
	}

	for _, opt := range opts {
		opt(subscription)
	}

	return subscription
}

func (s *Impl) GetTopic() string {
	return s.topic
}

func (s *Impl) GetQoS() byte {
	return s.qos
}

func (s *Impl) GetID() string {
	return s.topic
}

func (s *Impl) GetHandler() api.MessageHandler {
	return s.handler
}

func WithInterval(interval int) func(*Impl) {
	return func(s *Impl) {
		s.interval = uint32(interval)
	}
}

func WithConfig(config api.SubscriptionConfig) func(*Impl) {
	return func(s *Impl) {
		s.config = config
	}
}

func WithQoS(qos byte) func(*Impl) {
	return func(s *Impl) {
		s.qos = qos
	}
}
