package subscription

import (
	"encoding/json"
	"github.com/OI4/oi4-oec-service-go/service/api"
	tp "github.com/OI4/oi4-oec-service-go/service/topic"
	"github.com/eclipse/paho.mqtt.golang"
)

type MessageHandlerImpl struct {
	handler        func(mqtt.Client, mqtt.Message)
	skipOwnMessage bool
}

func NewMessageHandler(app api.Oi4Application, handler func(resource api.ResourceType, source *api.Oi4Identifier, networkMessage api.NetworkMessage, topic *tp.Topic), opts ...func(*MessageHandlerImpl)) *MessageHandlerImpl {

	messageHandler := &MessageHandlerImpl{
		skipOwnMessage: true,
	}

	for _, opt := range opts {
		opt(messageHandler)
	}

	handle := func(_ mqtt.Client, message mqtt.Message) {
		networkMessage := api.NetworkMessage{}
		err := json.Unmarshal(message.Payload(), &networkMessage)
		if err != nil {
			app.GetLogger().Infof("%s %s topic:%s", "error unmarshalling network message", err, message.Topic())
			return
		}
		topic, err := tp.ParseTopic(message.Topic())

		if err != nil {
			app.GetLogger().Infof("topic:%s invalid with: %v", message.Topic(), err)
			return
		}

		if messageHandler.skipOwnMessage && topic.HasSameApplication(app.GetServiceType(), app.GetApplicationSource().GetOi4Identifier()) {
			app.GetLogger().Debugf("skipping own message: %s", message.Topic())
			return
		}

		handler(topic.Resource, topic.Source, networkMessage, topic)
	}

	messageHandler.handler = handle

	return messageHandler
}

func (m *MessageHandlerImpl) GetHandler() mqtt.MessageHandler {
	return m.handler
}

func WithSkipOwnMessage(skip bool) func(*MessageHandlerImpl) {
	return func(s *MessageHandlerImpl) {
		s.skipOwnMessage = skip
	}
}
