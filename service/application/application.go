package application

import (
	"errors"
	"github.com/OI4/oi4-oec-service-go/service/api"
	pub "github.com/OI4/oi4-oec-service-go/service/application/publication"
	"github.com/OI4/oi4-oec-service-go/service/container"
	"github.com/OI4/oi4-oec-service-go/service/mqtt"
	"github.com/OI4/oi4-oec-service-go/service/opc"
	"github.com/OI4/oi4-oec-service-go/service/topic"
	"go.uber.org/zap"
	"sync"
	"time"
)

var (
	ErrPublisherAlreadyRegistered                       = errors.New("a publication with the same resource is already registered")
	ErrAssetAlreadyRegistered                           = errors.New("this asset is already assigned to a application")
	ErrPublicationAlreadyRegisteredOnAnotherApplication = errors.New("the publication is already registered on an asset or application")
)

// Oi4ApplicationImpl An OI4 Application host defined by the service type
type Oi4ApplicationImpl struct {
	mam           *api.MasterAssetModel
	oi4Identifier *api.Oi4Identifier
	serviceType   api.ServiceType

	mqttClient *mqtt.Client

	assets     map[api.Oi4Identifier]*Oi4Asset
	assetMutex sync.RWMutex

	publicationsList map[api.ResourceType]pub.Publication
	publicationMutex sync.RWMutex

	applicationSource api.Oi4ApplicationSource

	logger *zap.SugaredLogger
}

// CreateNewApplication Create a new Application host of a specific service type
func CreateNewApplication(serviceType api.ServiceType, applicationSource api.Oi4ApplicationSource, logger *zap.SugaredLogger) (*Oi4ApplicationImpl, error) {
	mam := applicationSource.GetMasterAssetModel()
	application := &Oi4ApplicationImpl{

		mam:           &mam,
		oi4Identifier: mam.ToOi4Identifier(),
		serviceType:   serviceType,

		assets:     make(map[api.Oi4Identifier]*Oi4Asset),
		assetMutex: sync.RWMutex{},

		publicationsList: make(map[api.ResourceType]pub.Publication),
		publicationMutex: sync.RWMutex{},

		applicationSource: applicationSource,
		logger:            logger,
	}
	applicationSource.SetOi4Application(application)

	return application, nil
}

func (app *Oi4ApplicationImpl) GetLogger() *zap.SugaredLogger {
	return app.logger
}

// RegisterPublication Register a publisher for the specific application
// you can overwrite built-in publications like MAM, Health etc...
func (app *Oi4ApplicationImpl) RegisterPublication(publication pub.Publication) error {
	app.publicationMutex.Lock()
	defer app.publicationMutex.Unlock()

	app.publicationsList[publication.GetResource()] = publication
	publication.Start()

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
	//publication.triggerSourcePublication(false, true, "")
	//app.triggerSourcePublication(false, true, "")
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

func (app *Oi4ApplicationImpl) SendPublicationMessage(publication api.PublicationMessage) {
	if app.mqttClient != nil && publication.Data != nil {
		// Deal with combined messages
		//var source *api.Oi4Identifier
		//if publication.source != nil &&
		//	(publication.publicationMode == api.PublicationMode_SOURCE_3 ||
		//		publication.publicationMode == api.PublicationMode_APPLICATION_SOURCE_FILTER_8 ||
		//		publication.publicationMode == api.PublicationMode_SOURCE_FILTER_7 ||
		//		publication.publicationMode == api.PublicationMode_APPLICATION_SOURCE_5) {
		//	source = publication.source
		//}
		source := publication.Source

		tp := topic.NewTopic(
			app.serviceType,
			*app.mam.ToOi4Identifier(),
			api.MethodPub,
			publication.Resource,
			source,
			nil,
			publication.Filter,
		)

		dswId := opc.GetDataSetWriterId(publication.Resource, *source)

		err := app.mqttClient.PublishResource(tp.ToString(), opc.CreateNetworkMessage(app.mam.ToOi4Identifier(), app.serviceType, publication.Resource, publication.Source, dswId, publication.CorrelationId, publication.Data))
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
	mqttClientOptions := &mqtt.ClientOptions{
		Host:     brokerConfig.Address,
		Port:     int(brokerConfig.SecurePort),
		Tls:      true,
		Username: credentials.Username(),
		Password: pwd,
	}

	client, err := mqtt.NewClient(mqttClientOptions)
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

		app.triggerSourcePublication(oi4Source, resource, filter, pub.OnRequest, networkMessage.MessageId)
	})

	err = app.registerPublications()
	if err != nil {
		return err
	}

	//applicationTicker := time.NewTicker(100 * time.Millisecond)
	//go func() {
	//	for {
	//		<-applicationTicker.C
	//
	//		data := api.NewOi4Data(rand.Float64())
	//
	//		addValue := func(key string, value any) {
	//			dErr := data.AddSecondaryData(key, &value)
	//
	//			if dErr != nil {
	//				logger.Error("Failed to add secondary data:", dErr)
	//			}
	//		}
	//
	//		addValue("Sv1", rand.Float64())
	//		addValue("Sv2", rand.Float64())
	//
	//		applicationSource.UpdateData(data, "Oi4Data")
	//	}
	//
	//}()
	return nil
}

func (app *Oi4ApplicationImpl) ResourceChanged(resource api.ResourceType, source api.Oi4Source, _ *string) {
	app.triggerSourcePublication(source, resource, nil, pub.OnRequest, "")
}

func (app *Oi4ApplicationImpl) registerPublications() error {
	// register built-in publications
	err := app.RegisterPublication(pub.NewIntervalBuilder(app, 60*time.Second). //
											Oi4Source(app.applicationSource).                          //
											Resource(api.ResourceData).                                //
											PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
											Build())
	//.SetDataFunc(func() *api.Health {health := application.applicationSource.GetHealth() return &health})
	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(app.applicationSource).                          //
								Resource(api.ResourceMam).                                 //
								PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
								PublishOnRegistration(true).                               //
								Build())

	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(app.applicationSource).                          //
								Resource(api.ResourceLicense).                             //
								PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
								Build())
	//SetDataFunc(func() *api.License {// Dummy implementation yet		components := make([]api.LicenseComponent, 0)		return &api.License{			Components: components,		}	}).

	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(app.applicationSource).                          //
								Resource(api.ResourceLicenseText).                         //
								PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
								Build())
	//SetDataFunc(func() *api.LicenseText {		// Dummy implementation yet		return &api.LicenseText{			LicenseText: "",		}	}).

	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(app.applicationSource).                          //
								Resource(api.ResourcePublicationList).                     //
								PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
								Build())
	//SetDataFunc(func() *api.PublicationList {

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
	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(app.applicationSource).                   //
								Resource(api.ResourceProfile).                      //
								PublicationMode(api.PublicationMode_APPLICATION_2). //
								Build())
	//SetDataFunc(func() *api.Profile {
	//resources := make([]api.ResourceType, 0)
	//for key := range application.publicationsList {
	//	resources = append(resources, key)
	//}
	//profile := api.Profile{
	//	Resources: resources,
	//}
	//return &profile
	//})
	if err != nil {
		return err
	}

	return nil
}

// func (app *Oi4ApplicationImpl) triggerSourcePublication(byInterval bool, onRequest bool, correlationId string) {
func (app *Oi4ApplicationImpl) triggerSourcePublication(source api.Oi4Source, resource api.ResourceType, filter api.Filter, trigger pub.Trigger, correlationId string) {
	publication := app.publicationsList[resource]
	if publication == nil {
		return
	}

	publication.TriggerPublication(trigger, correlationId)
}

func (app *Oi4ApplicationImpl) shouldPublicate(trigger pub.Trigger, publication *api.PublicationList) bool {
	if trigger == pub.OnRequest {
		return true
	}

	mode := *publication.Mode
	if mode == api.PublicationMode_OFF_0 || //
		mode == api.PublicationMode_ON_REQUEST_1 {
		return false
	}

	interval := *publication.Interval
	if interval == 0 && trigger != pub.ByInterval || //
		interval != 0 && trigger == pub.ByInterval {
		return true
	}

	return false
}

func (app *Oi4ApplicationImpl) sendGracefulShutdown() {
	app.SendPublicationMessage(api.PublicationMessage{
		Resource:   api.ResourceHealth,
		StatusCode: api.Status_Good,
		Source:     app.mam.ToOi4Identifier(),
		//publicationMode: api.PublicationMode_APPLICATION_SOURCE_5,
		Data: &api.Health{Health: api.Health_Normal, HealthScore: 0},
	})
}

// Stop application and shutdown all publications and assets
func (app *Oi4ApplicationImpl) Stop() {
	for _, publication := range app.publicationsList {
		publication.Stop()
	}
	app.sendGracefulShutdown()
	app.mqttClient.Stop()
}
