package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mzeiher/oi4/application/pkg/tls"
)

var (
	ErrNoAuthInformation = errors.New("no auth information provided, please provide either a mTLS certificate or username/password")
)

type MQTTClientOptions struct {
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

type MQTTClient struct {
	client mqtt.Client
}

func NewMQTTClient(options *MQTTClientOptions) (*MQTTClient, error) {
	client_options := mqtt.NewClientOptions()

	client_options.SetClientID("client")
	if options.Tls {
		client_options.AddBroker(fmt.Sprintf("ssl://%s:%d", options.Host, options.Port))
		tls_config, err := tls.NewTLSConfig(options.Ca_certificate_pem, options.Client_certificate_pem, options.Client_private_key_pem, options.TlsVerify)
		if err != nil {
			return nil, err
		}
		client_options.SetTLSConfig(tls_config)
	} else {
		client_options.AddBroker(fmt.Sprintf("tcp://%s:%d", options.Host, options.Port))
	}

	if options.Username != "" && options.Password != "" {
		client_options.Username = options.Username
		client_options.Password = options.Password
	}

	client := mqtt.NewClient(client_options)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	} else {
		return &MQTTClient{client: client}, nil
	}
}

func (client *MQTTClient) PublishResource(topic string, data interface{}) error {
	marshalledString, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if token := client.client.Publish(topic, 0, false, string(marshalledString)); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (client *MQTTClient) Stop() {
	client.client.Disconnect(1000)
}
