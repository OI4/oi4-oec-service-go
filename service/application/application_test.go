package application

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/OI4/oi4-oec-service-go/service/application/source"
	"github.com/OI4/oi4-oec-service-go/service/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/url"
	"testing"
)

func TestWithMockMqttClient(t *testing.T) {
	observedZapCore, _ := observer.New(zap.DebugLevel)
	logger := zap.New(observedZapCore)

	applicationSource := source.NewApplicationSourceImpl(api.MasterAssetModel{})

	mqttClientMock := &MqttClientMock{}
	option := WithMqttClientFn(func(options *api.MqttClientOptions) (api.MqttClient, error) {
		return mqttClientMock, nil
	})
	app := CreateNewApplication(api.ServiceTypeUtility, applicationSource, logger.Sugar(), option)
	require.NotNil(t, app)

	require.Nil(t, app.mqttClient)

	storage := container.Storage{
		MessageBusStorage: &container.MessageBusStorage{
			BrokerConfiguration: &container.BrokerConfiguration{
				Address:    "mqtt.example.com",
				SecurePort: 8883,
			},
		},
		SecretStorage: &container.SecretStorage{
			MqttCredentials: url.UserPassword("testuser", "testpassword"),
		},
	}

	err := app.Start(storage)
	require.NoError(t, err)
	assert.Equal(t, mqttClientMock, app.mqttClient)
	assert.True(t, mqttClientMock.RegisterGetHandlerCalled)
}

type MqttClientMock struct {
	PublishResourceFunc      func(topic string, msg interface{}) error
	SubscribeFunc            func(sub api.Subscription) error
	RegisterGetHandlerCalled bool
}

func (m *MqttClientMock) RegisterGetHandler(_ api.ServiceType, _ api.Oi4Identifier, _ byte, _ api.MessageHandler) error {
	m.RegisterGetHandlerCalled = true
	return nil
}

func (m *MqttClientMock) SubscribeToTopic(_ string, _ byte, _ api.MessageHandler) error {
	panic("implement me")
}

func (m *MqttClientMock) Stop() {
	panic("implement me")
}

func (m *MqttClientMock) PublishResource(topic string, _ byte, msg interface{}) error {
	if m.PublishResourceFunc != nil {
		return m.PublishResourceFunc(topic, msg)
	}
	return nil
}

func (m *MqttClientMock) Subscribe(sub api.Subscription) error {
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(sub)
	}
	return nil
}
