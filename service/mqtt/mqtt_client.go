package mqtt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/OI4/oi4-oec-service-go/service/tls"
	tp "github.com/OI4/oi4-oec-service-go/service/topic"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	ErrNoAuthInformation = errors.New("no auth information provided, please provide either a mTLS certificate or username/password")
)

type ClientOptions struct {
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

type Client struct {
	client mqtt.Client
}

func NewClient(options *ClientOptions) (*Client, error) {
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

func (client *Client) RegisterGetHandler(serviceType api.ServiceType, appId api.Oi4Identifier, handler func(resource api.ResourceType, source *api.Oi4Identifier, networkMessage api.NetworkMessage)) {
	topic := fmt.Sprintf("Oi4/%s/%s/Get/#", serviceType, appId.ToString())
	client.client.Subscribe(topic, 1, func(_ mqtt.Client, message mqtt.Message) {
		networkMessage := api.NetworkMessage{}
		err := json.Unmarshal(message.Payload(), &networkMessage)
		if err != nil {
			log.Printf("%s %s topic:%s", "error unmarshalling network message", err, message.Topic())
			return
		}
		topic, err := tp.ParseTopic(message.Topic())

		if err != nil {
			log.Printf("topic:%s invalid with: %v", message.Topic(), err)
			return
		}

		handler(topic.Resource, topic.Source, networkMessage)
	})
}

func (client *Client) Stop() {
	client.client.Disconnect(1000)
}
