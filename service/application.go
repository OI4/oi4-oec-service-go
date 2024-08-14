package service

import (
	"errors"
	"github.com/OI4/oi4-oec-service-go/service/pkg/topic"
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

// Oi4ApplicationImpl An OI4 Application host defined by the service type
type Oi4ApplicationImpl struct {
	PublicationPublisher

	mam           *oi4.MasterAssetModel
	oi4Identifier *oi4.Oi4Identifier
	serviceType   oi4.ServiceType

	mqttClient *mqtt.MQTTClient

	assets     map[oi4.Oi4Identifier]*Oi4Asset
	assetMutex sync.RWMutex

	publicationsList map[oi4.ResourceType]Publication
	publicationMutex sync.RWMutex

	applicationSource oi4.Oi4ApplicationSource
}

// CreateNewApplication Create a new Application host of a specific service type
func CreateNewApplication(serviceType oi4.ServiceType, applicationSource oi4.Oi4ApplicationSource) *Oi4ApplicationImpl {
	mam := applicationSource.GetMasterAssetModel()
	application := &Oi4ApplicationImpl{
		mam:           &mam,
		oi4Identifier: mam.ToOi4Identifier(),
		serviceType:   serviceType,

		assets:     make(map[oi4.Oi4Identifier]*Oi4Asset),
		assetMutex: sync.RWMutex{},

		publicationsList: make(map[oi4.ResourceType]Publication),
		publicationMutex: sync.RWMutex{},

		applicationSource: applicationSource,
	}
	applicationSource.SetOi4Application(application)

	// register built-in publications
	application.RegisterPublication(CreatePublication[*oi4.Health](oi4.ResourceHealth, true).SetDataFunc(func() *oi4.Health {
		health := application.applicationSource.GetHealth()
		return &health
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5).SetPublicationInterval(60 * time.Second))
	application.RegisterPublication(CreatePublication[*oi4.MasterAssetModel](oi4.ResourceMam, true).SetData(&mam).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))
	application.RegisterPublication(CreatePublication[*oi4.License](oi4.ResourceLicense, false).SetDataFunc(func() *oi4.License {
		// Dummy implementation yet
		components := make([]oi4.LicenseComponent, 0)
		return &oi4.License{
			Components: components,
		}
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))
	application.RegisterPublication(CreatePublication[*oi4.LicenseText](oi4.ResourceLicenseText, false).SetDataFunc(func() *oi4.LicenseText {
		// Dummy implementation yet
		return &oi4.LicenseText{
			LicenseText: "",
		}
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))

	application.RegisterPublication(CreatePublication[*oi4.PublicationList](oi4.ResourcePublicationList, false).SetDataFunc(func() *oi4.PublicationList {
		//publications := make([]oi4.PublicationList, len(application.publicationsList))
		var publicationList oi4.PublicationList
		for key := range application.publicationsList {
			publication := application.publicationsList[key]
			mode := publication.getPublicationMode()
			publicationList = oi4.PublicationList{
				ResourceType:    key,
				Source:          publication.getSource().ToString(),
				DataSetWriterId: opcmessages.GetDataSetWriterId(publication.getResource(), *publication.getSource()),
				Mode:            &mode,
			}
		}
		return &publicationList
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))

	application.RegisterPublication(CreatePublication[*oi4.Profile](oi4.ResourceProfile, false).SetDataFunc(func() *oi4.Profile {
		resources := make([]oi4.ResourceType, 0)
		for key := range application.publicationsList {
			resources = append(resources, key)
		}
		profile := oi4.Profile{
			Resources: resources,
		}
		return &profile
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_2))

	return application
}

// RegisterPublication Register a publisher for the specific application
// you can overwrite built-in publications like MAM, Health etc...
func (app *Oi4ApplicationImpl) RegisterPublication(publication Publication) error {
	app.publicationMutex.Lock()
	defer app.publicationMutex.Unlock()

	if publication.getParent() != nil {
		return ErrPublicationAlreadyRegisteredOnAnotherApplication
	}

	publication.setParent(app)
	publication.setSource(app.mam.ToOi4Identifier())
	app.publicationsList[publication.getResource()] = publication
	publication.start()

	return nil
}

// GetPublications Return all resources where a publication is registered
func (app *Oi4ApplicationImpl) GetPublications() []oi4.ResourceType {
	app.publicationMutex.RLock()
	defer app.publicationMutex.RUnlock()

	resources := make([]oi4.ResourceType, len(app.publicationsList))
	i := 0
	for key := range app.publicationsList {
		resources[i] = key
		i++
	}

	return resources
}

// RegisterAsset Add new asset to the application
func (app *Oi4ApplicationImpl) RegisterAsset(asset *Oi4Asset) {
	app.assetMutex.RLock()
	defer app.assetMutex.RUnlock()

	asset.setParent(app)
	oi4Id := asset.mam.ToOi4Identifier()
	app.assets[*oi4Id] = asset

	for _, publication := range asset.publicationsList {
		if publication.publishOnRegistration() {
			publication.triggerPublication(false, true, "")
		}
	}
}

// RemoveAsset remove an asset from the application
func (app *Oi4ApplicationImpl) RemoveAsset(asset *Oi4Asset) {
	app.assetMutex.RLock()
	defer app.assetMutex.RUnlock()

	asset.setParent(nil)
	delete(app.assets, *asset.mam.ToOi4Identifier())
}

func (app *Oi4ApplicationImpl) UpdateHealth(health oi4.Health) {
	app.publicationsList[oi4.ResourceHealth].(*PublicationImpl[*oi4.Health]).SetData(&health)
}

func (app *Oi4ApplicationImpl) GetMam() *oi4.MasterAssetModel {
	return app.mam
}

func (app *Oi4ApplicationImpl) sendPublicationMessage(publication PublicationMessage) {
	if app.mqttClient != nil && publication.data != nil {
		// Deal with combined messages
		//var source *oi4.Oi4Identifier
		//if publication.source != nil &&
		//	(publication.publicationMode == oi4.PublicationMode_SOURCE_3 ||
		//		publication.publicationMode == oi4.PublicationMode_APPLICATION_SOURCE_FILTER_8 ||
		//		publication.publicationMode == oi4.PublicationMode_SOURCE_FILTER_7 ||
		//		publication.publicationMode == oi4.PublicationMode_APPLICATION_SOURCE_5) {
		//	source = publication.source
		//}
		source := publication.source

		tp := topic.NewTopic(
			app.serviceType,
			*app.mam.ToOi4Identifier(),
			oi4.MethodPub,
			publication.resource,
			source,
			nil,
			nil,
		)

		dswId := opcmessages.GetDataSetWriterId(publication.resource, *source)

		err := app.mqttClient.PublishResource(tp.ToString(), opcmessages.CreateNetworkMessage(app.mam.ToOi4Identifier(), app.serviceType, publication.resource, publication.source, dswId, publication.correlationId, publication.data))
		if err != nil {
			return
		}

	}
}

// Start  application and connect to broker
func (app *Oi4ApplicationImpl) Start(mqttClientOptions *mqtt.MQTTClientOptions) error {

	client, err := mqtt.NewMQTTClient(mqttClientOptions)
	if err != nil {
		return err
	}
	app.mqttClient = client

	//app.mqttClient.RegisterGetHandler(app.serviceType, oi4.Oi4IdentifierString(app.mam.ToOi4Identifier().ToString()), func(resource oi4.ResourceType, source oi4.Oi4IdentifierString, networkMessage oi4.NetworkMessage) {
	app.mqttClient.RegisterGetHandler(app.serviceType, *app.mam.ToOi4Identifier(), func(resource oi4.ResourceType, source *oi4.Oi4Identifier, networkMessage oi4.NetworkMessage) {
		if source == nil {
			// TODO return all resources
		} else if source.Equals(app.mam.ToOi4Identifier()) {
			if publication := app.publicationsList[resource]; publication != nil {
				publication.triggerPublication(false, true, networkMessage.MessageId)
			}

		} else {
			if asset := app.assets[*source]; asset != nil {
				if publication := app.assets[*source].publicationsList[resource]; publication != nil {
					publication.triggerPublication(false, true, networkMessage.MessageId)
				}
			}
		}
	})

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

	return nil
}

func (app *Oi4ApplicationImpl) ResourceChanged(resource oi4.ResourceType, source oi4.Oi4Source) {
	var list map[oi4.ResourceType]Publication
	sourceId := source.GetOi4Identifier()
	if app.oi4Identifier.Equals(sourceId) {
		list = app.publicationsList
	} else {
		asset := app.assets[*sourceId]
		if asset != nil {
			list = asset.publicationsList
		}
	}
	if list == nil {
		return
	}

	for key := range list {
		if key == resource {
			list[key].triggerPublication(false, true, "")
		}
	}
}

func (app *Oi4ApplicationImpl) sendGracefulShutdown() {
	app.sendPublicationMessage(PublicationMessage{
		resource:        oi4.ResourceHealth,
		statusCode:      oi4.Status_Good,
		source:          app.mam.ToOi4Identifier(),
		publicationMode: oi4.PublicationMode_APPLICATION_SOURCE_5,
		data:            &oi4.Health{Health: oi4.Health_Normal, HealthScore: 0},
	})
}

// Stop application and shutdown all publications and assets
func (app *Oi4ApplicationImpl) Stop() {
	for _, publication := range app.publicationsList {
		publication.stop()
	}
	app.sendGracefulShutdown()
	app.mqttClient.Stop()
}
