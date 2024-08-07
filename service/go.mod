module github.com/OI4/oi4-oec-service-go/service

go 1.22.6

replace github.com/OI4/oi4-oec-service-go/api v0.0.0 => ../api

require (
	github.com/OI4/oi4-oec-service-go/api v0.0.0
	github.com/eclipse/paho.mqtt.golang v1.5.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
)
