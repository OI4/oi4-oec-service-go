package subscription

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	tp "github.com/OI4/oi4-oec-service-go/service/topic"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"testing"
)

var validTopic = "Oi4/OTConnector/acme.com/FBC/fbc%183z/FBC#123/Get/MAM/acme.com/matches/m/42-A/F234#862"

func TestExecuteValidMessage(t *testing.T) {
	logger := zaptest.NewLogger(t)

	app := applicationMock(logger.Sugar())
	handlerCalled := false

	message := messageMock(validPayload(), validTopic)

	handler := func(resource api.ResourceType, source *api.Oi4Identifier, networkMessage api.NetworkMessage, topic *tp.Topic) {
		handlerCalled = true
		assert.Equal(t, api.ResourceMam, resource)
		assert.NotNil(t, source)
		assert.Equal(t, "123", networkMessage.MessageId)
		assert.NotNil(t, topic)
	}

	messageHandler := NewMessageHandler(app, handler)
	messageHandler.GetHandler()(nil, message)

	assert.True(t, handlerCalled, "Handler should have been called")
}

func validPayload() string {
	return `{"MessageID":"123"}`
}

/**
Mocks
*/

func messageMock(payload string, topic string) mqtt.Message {
	return &mqttMessageMock{
		payload: []byte(payload),
		topic:   topic,
	}
}

type mqttMessageMock struct {
	payload []byte
	topic   string
}

func (m *mqttMessageMock) Duplicate() bool   { return false }
func (m *mqttMessageMock) Qos() byte         { return 0 }
func (m *mqttMessageMock) Retained() bool    { return false }
func (m *mqttMessageMock) Topic() string     { return m.topic }
func (m *mqttMessageMock) MessageID() uint16 { return 0 }
func (m *mqttMessageMock) Payload() []byte   { return m.payload }
func (m *mqttMessageMock) Ack()              {}

func applicationMock(logger *zap.SugaredLogger) api.Oi4Application {
	return &applicationMockImpl{
		logger: logger,
	}
}

type applicationMockImpl struct {
	logger *zap.SugaredLogger
}

func (a *applicationMockImpl) GetOi4Identifier() api.Oi4Identifier {
	panic("implement me")
}

func (a *applicationMockImpl) SendGetMessage(_ string, _ api.GetMessage) error {
	panic("implement me")
}

func (a *applicationMockImpl) GetServiceType() api.ServiceType {
	return api.ServiceTypeUtility
}

func (a *applicationMockImpl) ResourceChanged(resource api.ResourceType, source api.BaseSource, filter *api.Filter) {
	panic("implement me")
}

func (a *applicationMockImpl) SendPublicationMessage(publication api.PublicationMessage) {
	panic("implement me")
}

func (a *applicationMockImpl) GetIntervalPublicationScheduler() api.IntervalPublicationScheduler {
	panic("implement me")
}

func (a *applicationMockImpl) GetPublications() []api.Publication {
	panic("implement me")
}

func (a *applicationMockImpl) GetLogger() *zap.SugaredLogger {
	return a.logger
}

func (a *applicationMockImpl) GetApplicationSource() api.ApplicationSource {
	return &applicationSourceMock{}
}

type applicationSourceMock struct{}

func (a *applicationSourceMock) GetMasterAssetModel() api.MasterAssetModel {
	panic("implement me")
}

func (a *applicationSourceMock) GetHealth() api.Health {
	panic("implement me")
}

func (a *applicationSourceMock) UpdateHealth(health api.Health) {
	panic("implement me")
}

func (a *applicationSourceMock) GetData(filter *api.Filter) []api.Data {
	panic("implement me")
}

func (a *applicationSourceMock) UpdateData(data api.Data, dataTag string) {
	panic("implement me")
}

func (a *applicationSourceMock) GetConfig() api.PublishConfig {
	panic("implement me")
}

func (a *applicationSourceMock) GetProfile() api.Profile {
	panic("implement me")
}

func (a *applicationSourceMock) GetLicense() api.License {
	panic("implement me")
}

func (a *applicationSourceMock) GetLicenseText(filter *api.Filter) []api.LicenseText {
	panic("implement me")
}

func (a *applicationSourceMock) GetLicenseTexts() map[string]api.LicenseText {
	panic("implement me")
}

func (a *applicationSourceMock) GetRtLicense() api.RtLicense {
	panic("implement me")
}

func (a *applicationSourceMock) GetPublicationList() []api.PublicationList {
	panic("implement me")
}

func (a *applicationSourceMock) GetSubscriptionList() []api.SubscriptionList {
	panic("implement me")
}

func (a *applicationSourceMock) GetReferenceDesignation() api.ReferenceDesignation {
	panic("implement me")
}

func (a *applicationSourceMock) Get(resource api.ResourceType, filter *api.Filter) []any {
	panic("implement me")
}

func (a *applicationSourceMock) SetOi4Application(application api.Oi4Application) {
	panic("implement me")
}

func (a *applicationSourceMock) Equals(source api.BaseSource) bool {
	panic("implement me")
}

func (a *applicationSourceMock) GetSources() map[api.Oi4Identifier]*api.AssetSource {
	panic("implement me")
}

func (a *applicationSourceMock) AddSource(source api.AssetSource) {
	panic("implement me")
}

func (a *applicationSourceMock) RemoveSource(identifier api.Oi4Identifier) {
	panic("implement me")
}

func (a *applicationSourceMock) GetOi4Identifier() *api.Oi4Identifier {
	return &api.Oi4Identifier{}
}
