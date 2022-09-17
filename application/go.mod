module github.com/mzeiher/oi4/application

go 1.18

replace github.com/mzeiher/oi4/api v0.0.0 => ../api

require (
	github.com/eclipse/paho.mqtt.golang v1.4.1
	github.com/mzeiher/oi4/api v0.0.0
)

require (
	github.com/gorilla/websocket v1.4.2 // indirect
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)
