package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	oi4 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
	"github.com/OI4/oi4-oec-service-go/service/pkg/mqtt"
	opcmessages "github.com/OI4/oi4-oec-service-go/service/pkg/opc_messages"
)

var (
	ErrPublisherAlreadyRegistered                       = errors.New("a publication with the same resource is already registered")
	ErrAssetAlreadyRegistered                           = errors.New("this asset is already assigned to a application")
	ErrPublicationAlreadyRegisteredOnAnotherApplication = errors.New("the publication is already registered on an asset or application")
)

// An OI4 Application host defined by the service type
type Oi4Application struct {
	PublicationPublisher

	mam         *oi4.MasterAssetModel
	serviceType oi4.ServiceType

	mqttClient *mqtt.MQTTClient

	assets     map[oi4.Oi4IdentifierPath]*Oi4Asset
	assetMutex sync.RWMutex

	publicationsList map[oi4.Resource]Publication
	publicationMutex sync.RWMutex
}

// Create a new Application host of a specific service type
func CreateNewApplication(serviceType oi4.ServiceType, mam *oi4.MasterAssetModel) *Oi4Application {
	application := &Oi4Application{
		mam:         mam,
		serviceType: serviceType,

		assets:     make(map[oi4.Oi4IdentifierPath]*Oi4Asset),
		assetMutex: sync.RWMutex{},

		publicationsList: make(map[oi4.Resource]Publication),
		publicationMutex: sync.RWMutex{},
	}

	// register built-in publications
	application.RegisterPublication(CreatePublication[*oi4.Health](oi4.Resource_Health, true).SetData(&oi4.Health{Health: oi4.Health_Normal, HealthScore: 100}).SetPublicationMode(oi4.PublicationMode_APPLICATION_2).SetPublicationInterval(60 * time.Second))
	application.RegisterPublication(CreatePublication[*oi4.MasterAssetModel](oi4.Resource_MAM, true).SetData(mam).SetPublicationMode(oi4.PublicationMode_APPLICATION_2))

	application.RegisterPublication(CreatePublication[[]oi4.Resource](oi4.Resource_Profile, false).SetDataFunc(func() []oi4.Resource {
		resources := make([]oi4.Resource, 0)
		for key := range application.publicationsList {
			resources = append(resources, key)
		}
		return resources
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_2))

	return application
}

// Register a publisher for the specific application
// you can overwrite built-in publications like MAM, Health etc...
func (app *Oi4Application) RegisterPublication(publication Publication) error {
	app.publicationMutex.Lock()
	defer app.publicationMutex.Unlock()

	if publication.getParent() != nil {
		return ErrPublicationAlreadyRegisteredOnAnotherApplication
	}

	publication.setParent(app)
	app.publicationsList[publication.getResource()] = publication
	publication.start()

	return nil
}

// Return all resources where a publication is registered
func (app *Oi4Application) GetPublications() []oi4.Resource {
	app.publicationMutex.RLock()
	defer app.publicationMutex.RUnlock()

	resources := make([]oi4.Resource, len(app.publicationsList))
	i := 0
	for key := range app.publicationsList {
		resources[i] = key
		i++
	}

	return resources
}

// Add new asset to the application
func (app *Oi4Application) RegisterAsset(asset *Oi4Asset) {
	app.assetMutex.RLock()
	defer app.assetMutex.RUnlock()

	asset.setParent(app)
	app.assets[oi4.Oi4IdentifierPath(asset.mam.ToOi4Identifier().ToString())] = asset

	for _, publication := range asset.publicationsList {
		if publication.publishOnRegistration() {
			publication.triggerPublication(false, true, "")
		}
	}
}

// remove an asset from the application
func (app *Oi4Application) RemoveAsset(asset *Oi4Asset) {
	app.assetMutex.RLock()
	defer app.assetMutex.RUnlock()

	asset.setParent(nil)
	delete(app.assets, oi4.Oi4IdentifierPath(asset.mam.ToOi4Identifier().ToString()))
}

func (app *Oi4Application) UpdateHealth(health oi4.Health) {
	app.publicationsList[oi4.Resource_Health].(*PublicationImpl[*oi4.Health]).SetData(&health)
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

		app.mqttClient.PublishResource(topic, opcmessages.CreateNetworkMessage(app.mam.ToOi4Identifier(), app.serviceType, publication.resource, publication.source, publication.dataSetWriterId, publication.correlationId, publication.data))

	}
}

// start application and connect to broker
func (app *Oi4Application) Start(mqttClientOptions *mqtt.MQTTClientOptions) error {

	client, err := mqtt.NewMQTTClient(mqttClientOptions)
	if err != nil {
		return err
	}
	app.mqttClient = client

	// trigger publications for application
	for _, publication := range app.publicationsList {
		if publication.publishOnRegistration() {
			publication.triggerPublication(false, true, "")
		}
	}
	// trigger publications for assets
	for _, asset := range app.assets {
		for _, publication := range asset.publicationsList {
			if publication.publishOnRegistration() {
				publication.triggerPublication(false, true, "")
			}
		}
	}

	app.mqttClient.RegisterGetHandler(app.serviceType, oi4.Oi4IdentifierPath(app.mam.ToOi4Identifier().ToString()), func(resource oi4.Resource, source oi4.Oi4IdentifierPath, networkMessage oi4.NetworkMessage) {
		if source != "" {
			if asset := app.assets[source]; asset != nil {
				if publication := app.assets[source].publicationsList[resource]; publication != nil {
					publication.triggerPublication(false, true, networkMessage.MessageId)
				}
			}
		} else {
			if publication := app.publicationsList[resource]; publication != nil {
				publication.triggerPublication(false, true, networkMessage.MessageId)
			}
		}
	})

	return nil
}

// stop application and shutdown all publications and assets
func (app *Oi4Application) Stop() {
	for _, publication := range app.publicationsList {
		publication.stop()
	}
	app.mqttClient.Stop()
}
