package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/OI4/oi4-oec-service-go/service/tls"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	ErrNoAuthInformation = errors.New("no auth information provided, please provide either a mTLS certificate or username/password")
)

type Client struct {
	client mqtt.Client
}

func NewClient(options *api.MqttClientOptions) (*Client, error) {
	clientOptions := mqtt.NewClientOptions()

	clientOptions.SetClientID("client")
	if options.Tls {
		clientOptions.AddBroker(fmt.Sprintf("ssl://%s:%d", options.Host, options.Port))
		tlsConfig, err := tls.NewTLSConfig(options.Ca_certificate_pem, options.Client_certificate_pem, options.Client_private_key_pem, options.TlsVerify)
		if err != nil {
			return nil, err
		}
		clientOptions.SetTLSConfig(tlsConfig)
	} else {
		clientOptions.AddBroker(fmt.Sprintf("tcp://%s:%d", options.Host, options.Port))
	}

	if options.Username != "" && options.Password != "" {
		clientOptions.Username = options.Username
		clientOptions.Password = options.Password
	}

	client := mqtt.NewClient(clientOptions)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	} else {
		return &Client{client: client}, nil
	}
}

func (client *Client) PublishResource(topic string, data interface{}) error {
	marshalledString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if token := client.client.Publish(topic, 0, false, string(marshalledString)); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (client *Client) RegisterGetHandler(serviceType api.ServiceType, appId api.Oi4Identifier, qos byte, handler api.MessageHandler) error {
	// TODO Parse topic and provide to handler
	topic := fmt.Sprintf("Oi4/%s/%s/Get/#", serviceType, appId.ToString())
	return client.SubscribeToTopic(topic, qos, handler)
}

func (client *Client) Subscribe(subscription api.Subscription) error {
	return client.SubscribeToTopic(subscription.GetTopic(), subscription.GetQoS(), subscription.GetHandler())
}

func (client *Client) SubscribeToTopic(topic string, qos byte, handler api.MessageHandler) error {
	token := client.client.Subscribe(topic, qos, handler.GetHandler())
	return token.Error()
}

func (client *Client) Stop() {
	client.client.Disconnect(1000)
}
