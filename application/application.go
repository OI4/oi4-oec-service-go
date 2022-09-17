package application

import (
	"sync"
	"time"

	oi4 "github.com/mzeiher/oi4/api/v1"
	"github.com/mzeiher/oi4/application/pkg/mqtt"
)

type Oi4Application struct {
	mam         *oi4.MasterAssetModel
	serviceType oi4.ServiceType

	mqttClient *mqtt.MQTTClient

	assets     map[oi4.Oi4IdentifierPath]*Oi4Asset
	assetMutex sync.Mutex

	health oi4.Health
	data   oi4.Oi4Data

	stopHealthInterval chan struct{}
}

func CreateNewApplication(serviceType oi4.ServiceType, mam *oi4.MasterAssetModel) *Oi4Application {
	return &Oi4Application{
		mam:                mam,
		serviceType:        serviceType,
		assetMutex:         sync.Mutex{},
		assets:             make(map[oi4.Oi4IdentifierPath]*Oi4Asset),
		health:             oi4.Health{Health: oi4.Health_Normal, HealthScore: 100},
		stopHealthInterval: make(chan struct{}),
	}
}

func (app *Oi4Application) RegisterAsset(asset *Oi4Asset) {
	app.assetMutex.Lock()
	defer app.assetMutex.Unlock()

	asset.parent = app
	app.assets[oi4.Oi4IdentifierPath(asset.mam.ToOi4Identifier().ToString())] = asset
}

func (app *Oi4Application) RemoveAsset(asset *Oi4Asset) {
	app.assetMutex.Lock()
	defer app.assetMutex.Unlock()

	asset.parent = nil
	delete(app.assets, oi4.Oi4IdentifierPath(asset.mam.ToOi4Identifier().ToString()))
}

func (app *Oi4Application) UpdateHealth(health oi4.Health) {
	app.health = health
}

func (app *Oi4Application) UpdateData(data oi4.Oi4Data) {
	app.data = data
}

func (app *Oi4Application) AddCallReplyHandler(handler interface{}) {
	panic("not implemented")
}

func (app *Oi4Application) Start(mqttClientOptions *mqtt.MQTTClientOptions) error {
	client, err := mqtt.NewMQTTClient(mqttClientOptions)
	if err != nil {
		return err
	}
	app.mqttClient = client

	client.PublishResource(app.mam.ToOi4Identifier(), app.serviceType, oi4.Resource_MAM, nil, app.mam)

	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for {
			select {
			case <-ticker.C:
				client.PublishResource(app.mam.ToOi4Identifier(), app.serviceType, oi4.Resource_Health, nil, app.health)
			case <-app.stopHealthInterval:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

func (app *Oi4Application) Stop() {
	app.stopHealthInterval <- struct{}{}
	app.mqttClient.Stop()
}
