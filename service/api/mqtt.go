package api

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MessageHandler interface {
	GetHandler() mqtt.MessageHandler
}
