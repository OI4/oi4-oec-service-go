package api

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MqttClientOptions struct {
	Host                          string
	Tls                           bool
	Port                          int
	Username                      string
	Password                      string
	Client_private_key_pem        string
	Client_private_key_passphrase string
	Client_certificate_pem        string
	Ca_certificate_pem            string
	TlsVerify                     bool
}

type MessageHandler interface {
	GetHandler() mqtt.MessageHandler
}
type MqttClient interface {
	PublishResource(topic string, data interface{}) error
	RegisterGetHandler(serviceType ServiceType, appId Oi4Identifier, qos byte, handler MessageHandler) error
	Subscribe(subscription Subscription) error
	SubscribeToTopic(topic string, qos byte, handler MessageHandler) error
	Stop()
}
