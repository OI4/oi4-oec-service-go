package application

import (
	"errors"
	"github.com/OI4/oi4-oec-service-go/service/container"
	"github.com/OI4/oi4-oec-service-go/service/topic"
	"sync"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/OI4/oi4-oec-service-go/service/mqtt"
	"github.com/OI4/oi4-oec-service-go/service/opc"
)

var (
	ErrPublisherAlreadyRegistered                       = errors.New("a publication with the same resource is already registered")
	ErrAssetAlreadyRegistered                           = errors.New("this asset is already assigned to a application")
	ErrPublicationAlreadyRegisteredOnAnotherApplication = errors.New("the publication is already registered on an asset or application")
)

// Oi4ApplicationImpl An OI4 Application host defined by the service type
type Oi4ApplicationImpl struct {
	PublicationPublisher

	mam           *api.MasterAssetModel
	oi4Identifier *api.Oi4Identifier
	serviceType   api.ServiceType

	mqttClient *mqtt.MQTTClient

	assets     map[api.Oi4Identifier]*Oi4Asset
	assetMutex sync.RWMutex

	publicationsList map[api.ResourceType]Publication
	publicationMutex sync.RWMutex

	applicationSource api.Oi4ApplicationSource
}

// CreateNewApplication Create a new Application host of a specific service type
func CreateNewApplication(serviceType api.ServiceType, applicationSource api.Oi4ApplicationSource) *Oi4ApplicationImpl {
	mam := applicationSource.GetMasterAssetModel()
	application := &Oi4ApplicationImpl{
		mam:           &mam,
		oi4Identifier: mam.ToOi4Identifier(),
		serviceType:   serviceType,

		assets:     make(map[api.Oi4Identifier]*Oi4Asset),
		assetMutex: sync.RWMutex{},

		publicationsList: make(map[api.ResourceType]Publication),
		publicationMutex: sync.RWMutex{},

		applicationSource: applicationSource,
	}
	applicationSource.SetOi4Application(application)

	source := api.Oi4Source(applicationSource)
	// register built-in publications
	application.RegisterPublication(CreatePublication(api.ResourceHealth, &source). //.SetDataFunc(func() *api.Health {health := application.applicationSource.GetHealth() return &health})
											SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5).SetPublicationInterval(60 * time.Second))
	application.RegisterPublication(CreatePublication(api.ResourceMam, &source). //SetData(&mam).
											SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5))
	application.RegisterPublication(CreatePublication(api.ResourceLicense, &source). //SetDataFunc(func() *api.License {// Dummy implementation yet		components := make([]api.LicenseComponent, 0)		return &api.License{			Components: components,		}	}).
												SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5))
	application.RegisterPublication(CreatePublication(api.ResourceLicenseText, &source). //SetDataFunc(func() *api.LicenseText {		// Dummy implementation yet		return &api.LicenseText{			LicenseText: "",		}	}).
												SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5))

	application.RegisterPublication(CreatePublication(api.ResourcePublicationList, &source). //SetDataFunc(func() *api.PublicationList {
		//publications := make([]api.PublicationList, len(application.publicationsList))
		//var publicationList api.PublicationList
		//for key := range application.publicationsList {
		//	publication := application.publicationsList[key]
		//	mode := publication.getPublicationMode()
		//	publicationList = api.PublicationList{
		//		ResourceType:    key,
		//		Source:          publication.getSource().ToString(),
		//		DataSetWriterId: opc.GetDataSetWriterId(publication.getResource(), *publication.getSource()),
		//		Mode:            &mode,
		//	}
		//}
		//return &publicationList
		//}).
		SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5))

	application.RegisterPublication(CreatePublication(api.ResourceProfile, &source). //SetDataFunc(func() *api.Profile {
		//resources := make([]api.ResourceType, 0)
		//for key := range application.publicationsList {
		//	resources = append(resources, key)
		//}
		//profile := api.Profile{
		//	Resources: resources,
		//}
		//return &profile
		//})
		SetPublicationMode(api.PublicationMode_APPLICATION_2))

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
func (app *Oi4ApplicationImpl) GetPublications() []api.ResourceType {
	app.publicationMutex.RLock()
	defer app.publicationMutex.RUnlock()

	resources := make([]api.ResourceType, len(app.publicationsList))
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

	//for _, publication := range asset.publicationsList {
	//if publication.publishOnRegistration() {
	//publication.triggerPublication(false, true, "")
	//app.triggerPublication(false, true, "")
	//}
	//}
}

// RemoveAsset remove an asset from the application
func (app *Oi4ApplicationImpl) RemoveAsset(asset *Oi4Asset) {
	app.assetMutex.RLock()
	defer app.assetMutex.RUnlock()

	asset.setParent(nil)
	delete(app.assets, *asset.mam.ToOi4Identifier())
}

func (app *Oi4ApplicationImpl) UpdateHealth(health api.Health) {
	// TODO implement me
	//app.publicationsList[api.ResourceHealth].(*PublicationImpl).SetData(&health)
}

func (app *Oi4ApplicationImpl) GetMam() *api.MasterAssetModel {
	return app.mam
}

func (app *Oi4ApplicationImpl) sendPublicationMessage(publication PublicationMessage) {
	if app.mqttClient != nil && publication.data != nil {
		// Deal with combined messages
		//var source *api.Oi4Identifier
		//if publication.source != nil &&
		//	(publication.publicationMode == api.PublicationMode_SOURCE_3 ||
		//		publication.publicationMode == api.PublicationMode_APPLICATION_SOURCE_FILTER_8 ||
		//		publication.publicationMode == api.PublicationMode_SOURCE_FILTER_7 ||
		//		publication.publicationMode == api.PublicationMode_APPLICATION_SOURCE_5) {
		//	source = publication.source
		//}
		source := publication.source

		tp := topic.NewTopic(
			app.serviceType,
			*app.mam.ToOi4Identifier(),
			api.MethodPub,
			publication.resource,
			source,
			nil,
			publication.filter,
		)

		dswId := opc.GetDataSetWriterId(publication.resource, *source)

		err := app.mqttClient.PublishResource(tp.ToString(), opc.CreateNetworkMessage(app.mam.ToOi4Identifier(), app.serviceType, publication.resource, publication.source, dswId, publication.correlationId, publication.data))
		if err != nil {
			return
		}

	}
}

// Start  application and connect to broker
func (app *Oi4ApplicationImpl) Start(storage container.Storage) error {
	brokerConfig := storage.MessageBusStorage.BrokerConfiguration
	credentials := storage.SecretStorage.MqttCredentials
	pwd, _ := credentials.Password()
	mqttClientOptions := &mqtt.MQTTClientOptions{
		Host:     brokerConfig.Address,
		Port:     int(brokerConfig.SecurePort),
		Tls:      true,
		Username: credentials.Username(),
		Password: pwd,
	}

	client, err := mqtt.NewMQTTClient(mqttClientOptions)
	if err != nil {
		return err
	}
	app.mqttClient = client

	//app.mqttClient.RegisterGetHandler(app.serviceType, api.Oi4IdentifierString(app.mam.ToOi4Identifier().ToString()), func(resource api.ResourceType, source api.Oi4IdentifierString, networkMessage api.NetworkMessage) {
	app.mqttClient.RegisterGetHandler(app.serviceType, *app.mam.ToOi4Identifier(), func(resource api.ResourceType, source *api.Oi4Identifier, networkMessage api.NetworkMessage) {
		if source == nil {
			// TODO return all resources
			return
		}
		var oi4Source api.Oi4Source
		if source.Equals(app.mam.ToOi4Identifier()) {
			oi4Source = app.applicationSource
		} else {
			oi4Source = app.assets[*source].source
		}

		if oi4Source == nil {
			return
		}

		var filter api.Filter
		if len(networkMessage.Messages) > 0 {
			filter = networkMessage.Messages[0].Filter
		}

		app.triggerPublication(oi4Source, resource, filter, OnRequest, networkMessage.MessageId)
	})

	// trigger publications for application
	//for _, publication := range app.publicationsList {
	//	if publication.publishOnRegistration() {
	//		app.triggerPublication(false, true, "")
	//	}
	//}
	// trigger publications for assets
	//for _, asset := range app.assets {
	//	for _, publication := range asset.publicationsList {
	//		if publication.publishOnRegistration() {
	//			app.triggerPublication(false, true, "")
	//		}
	//	}
	//}

	return nil
}

func (app *Oi4ApplicationImpl) ResourceChanged(resource api.ResourceType, source api.Oi4Source, _ *string) {
	app.triggerPublication(source, resource, nil, OnRequest, "")
}

// func (app *Oi4ApplicationImpl) triggerPublication(byInterval bool, onRequest bool, correlationId string) {
func (app *Oi4ApplicationImpl) triggerPublication(source api.Oi4Source, resource api.ResourceType, filter api.Filter, trigger Trigger, correlationId string) {
	var publication *api.PublicationList
	for _, pub := range source.GetPublicationList() {
		if pub.ResourceType == resource {
			publication = &pub
		}
	}
	if publication == nil {
		return
	}

	if !app.shouldPublicate(trigger, publication) {
		return
	}

	message := PublicationMessage{
		resource:   resource,
		statusCode: api.Status_Good,
		//publicationMode: p.publicationMode,
		correlationId: correlationId,
		source:        source.GetOi4Identifier(),
		filter:        filter,
	}
	message.data = source.Get(resource)

	app.sendPublicationMessage(message)
}

func (app *Oi4ApplicationImpl) shouldPublicate(trigger Trigger, publication *api.PublicationList) bool {
	if trigger == OnRequest {
		return true
	}

	mode := *publication.Mode
	if mode == api.PublicationMode_OFF_0 || //
		mode == api.PublicationMode_ON_REQUEST_1 {
		return false
	}

	interval := *publication.Interval
	if interval == 0 && trigger != ByInterval || //
		interval != 0 && trigger == ByInterval {
		return true
	}

	return false
}

func (app *Oi4ApplicationImpl) sendGracefulShutdown() {
	app.sendPublicationMessage(PublicationMessage{
		resource:   api.ResourceHealth,
		statusCode: api.Status_Good,
		source:     app.mam.ToOi4Identifier(),
		//publicationMode: api.PublicationMode_APPLICATION_SOURCE_5,
		data: &api.Health{Health: api.Health_Normal, HealthScore: 0},
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
