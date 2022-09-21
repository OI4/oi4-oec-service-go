package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	oi4 "github.com/mzeiher/oi4/api/pkg/types"
	"github.com/mzeiher/oi4/service/pkg/mqtt"
	opcmessages "github.com/mzeiher/oi4/service/pkg/opc_messages"
)

var (
	ErrPublisherAlreadyRegistered                       = errors.New("a publication with the same resource is already registered")
	ErrPublicationAlreadyRegisteredOnAnotherApplication = errors.New("the publication is already registered on an asset or application")
)

type Oi4Application struct {
	Publisher

	mam         *oi4.MasterAssetModel
	serviceType oi4.ServiceType

	mqttClient *mqtt.MQTTClient

	assets     map[oi4.Oi4IdentifierPath]*Oi4Asset
	assetMutex sync.Mutex

	publicationsList map[oi4.Resource]*Publication
	publicationMutex sync.Mutex

	stopHealthInterval chan struct{}
}

// Create a new Application
//
// The application needs to be started and connected with the Connect Method
//
func CreateNewApplication(serviceType oi4.ServiceType, mam *oi4.MasterAssetModel) *Oi4Application {
	application := &Oi4Application{
		mam:         mam,
		serviceType: serviceType,

		assets:     make(map[oi4.Oi4IdentifierPath]*Oi4Asset),
		assetMutex: sync.Mutex{},

		publicationsList: make(map[oi4.Resource]*Publication),
		publicationMutex: sync.Mutex{},
	}

	application.RegisterPublication(CreatePublication(oi4.Resource_Health, true).SetData(oi4.Health{Health: oi4.Health_Normal, HealthScore: 100}).SetPublicationMode(oi4.PublicationMode_APPLICATION_2).SetPublicationInterval(60 * time.Second))
	application.RegisterPublication(CreatePublication(oi4.Resource_MAM, true).SetData(mam).SetPublicationMode(oi4.PublicationMode_APPLICATION_2))

	return application
}

func (app *Oi4Application) RegisterPublication(publication *Publication) error {
	app.publicationMutex.Lock()
	defer app.publicationMutex.Unlock()

	if app.publicationsList[publication.resource] != nil {
		return ErrPublisherAlreadyRegistered
	}

	if publication.parent != nil {
		return ErrPublicationAlreadyRegisteredOnAnotherApplication
	}

	publication.parent = app
	app.publicationsList[publication.resource] = publication
	publication.start()

	return nil
}

// // Do we need this?!
// func (app *Oi4Application) RemovePublication(publication *Publication) *Publication {
// 	app.publicationMutex.Lock()
// 	defer app.publicationMutex.Unlock()

// 	delete(app.publicationsList, publication.resource)
// 	publication.parent = nil
// 	publication.stop()

// 	return publication
// }

func (app *Oi4Application) GetPublications() []oi4.Resource {
	app.publicationMutex.Lock()
	defer app.publicationMutex.Unlock()

	resources := make([]oi4.Resource, len(app.publicationsList))
	i := 0
	for key, _ := range app.publicationsList {
		resources[i] = key
		i++
	}

	return resources
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
	app.publicationsList[oi4.Resource_Health].SetData(health)
}

func (app *Oi4Application) UpdateData(data oi4.Oi4Data) {
	app.publicationsList[oi4.Resource_Data].SetData(data)
}

func (app *Oi4Application) sendPublicationMessage(publication PublicationMessage) {
	if app.mqttClient != nil && publication.data != nil {
		topic := fmt.Sprintf("Oi4/%s/%s/Pub/%s", app.serviceType, app.mam.ToOi4Identifier().ToString(), publication.resource)
		if publication.source != nil &&
			(publication.publicationMode == oi4.PublicationMode_SOURCE_3 ||
				publication.publicationMode == oi4.PublicationMode_APPLICATION_SOURCE_FILTER_8 ||
				publication.publicationMode == oi4.PublicationMode_SOURCE_FILTER_7 ||
				publication.publicationMode == oi4.PublicationMode_APPLICATION_SOURCE_5) {
			topic = fmt.Sprintf("%s/%s", topic, publication.source.ToString())
		}

		app.mqttClient.PublishResource(topic, opcmessages.CreateNetworkMessage(app.mam.ToOi4Identifier(), app.serviceType, publication.resource, nil, publication.dataSetWriterId, publication.correlationId, publication.data))

	}
}

func (app *Oi4Application) Start(mqttClientOptions *mqtt.MQTTClientOptions) error {

	client, err := mqtt.NewMQTTClient(mqttClientOptions)
	if err != nil {
		return err
	}
	app.mqttClient = client

	// trigger publications
	for _, publication := range app.publicationsList {
		if publication.publishOnRegistration {
			publication.triggerPublication(false, true, "")
		}
	}

	app.mqttClient.RegisterGetHandler(app.serviceType, oi4.Oi4IdentifierPath(app.mam.ToOi4Identifier().ToString()), func(resource oi4.Resource, source oi4.Oi4IdentifierPath, networkMessage oi4.NetworkMessage) {
		if publication := app.publicationsList[resource]; publication != nil {
			publication.triggerPublication(false, true, networkMessage.MessageId)
		}
	})

	return nil
}

func (app *Oi4Application) Stop() {
	for _, publication := range app.publicationsList {
		publication.stop()
	}
	app.mqttClient.Stop()
}
