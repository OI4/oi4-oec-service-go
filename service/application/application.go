package application

import (
	"errors"
	"github.com/OI4/oi4-oec-service-go/service/api"
	pub "github.com/OI4/oi4-oec-service-go/service/application/publication"
	"github.com/OI4/oi4-oec-service-go/service/application/subscription"
	"github.com/OI4/oi4-oec-service-go/service/container"
	"github.com/OI4/oi4-oec-service-go/service/mqtt"
	"github.com/OI4/oi4-oec-service-go/service/opc"
	tp "github.com/OI4/oi4-oec-service-go/service/topic"
	"go.uber.org/zap"
	"maps"
	"slices"
	"sync"
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

	mqttClient api.MqttClient

	assets     map[api.Oi4Identifier]*AssetImpl
	assetMutex sync.RWMutex

	publications     map[api.ResourceType][]api.Publication
	publicationMutex sync.RWMutex

	subscriptions      map[string]api.Subscription
	subscriptionsMutex sync.RWMutex

	applicationSource api.ApplicationSource

	logger *zap.SugaredLogger

	scheduler api.IntervalPublicationScheduler

	createMqttClientFn func(options *api.MqttClientOptions) (api.MqttClient, error)
}

// CreateNewApplication Create a new Application host of a specific service type
func CreateNewApplication(serviceType api.ServiceType, applicationSource api.ApplicationSource, logger *zap.SugaredLogger, options ...Option) (*Oi4ApplicationImpl, error) {
	mam := applicationSource.GetMasterAssetModel()
	scheduler := pub.NewIntervalPublicationSchedulerImpl(50, 5)
	application := &Oi4ApplicationImpl{

		mam:           &mam,
		oi4Identifier: mam.ToOi4Identifier(),
		serviceType:   serviceType,

		assets:     make(map[api.Oi4Identifier]*AssetImpl),
		assetMutex: sync.RWMutex{},

		publications:     make(map[api.ResourceType][]api.Publication),
		publicationMutex: sync.RWMutex{},

		subscriptions:      make(map[string]api.Subscription),
		subscriptionsMutex: sync.RWMutex{},

		applicationSource: applicationSource,
		logger:            logger,
		scheduler:         scheduler,
	}
	applicationSource.SetOi4Application(application)

	for _, opt := range options {
		opt(application)
	}

	return application, nil
}

func (app *Oi4ApplicationImpl) GetServiceType() api.ServiceType {
	return app.serviceType
}

func (app *Oi4ApplicationImpl) GetOi4Identifier() api.Oi4Identifier {
	return *app.oi4Identifier
}

func (app *Oi4ApplicationImpl) GetApplicationSource() api.ApplicationSource {
	return app.applicationSource
}

func (app *Oi4ApplicationImpl) GetLogger() *zap.SugaredLogger {
	return app.logger
}

func (app *Oi4ApplicationImpl) GetIntervalPublicationScheduler() api.IntervalPublicationScheduler {
	return app.scheduler
}

// Start an application and connect to a broker
func (app *Oi4ApplicationImpl) Start(storage container.Storage) error {
	brokerConfig := storage.MessageBusStorage.BrokerConfiguration
	credentials := storage.SecretStorage.MqttCredentials
	pwd, _ := credentials.Password()
	mqttClientOptions := &api.MqttClientOptions{
		Host:     brokerConfig.Address,
		Port:     int(brokerConfig.SecurePort),
		Tls:      true,
		Username: credentials.Username(),
		Password: pwd,
	}

	var err error
	if app.mqttClient, err = app.newMqttClient(mqttClientOptions); err != nil {
		return err
	}

	if err = app.mqttClient.RegisterGetHandler(app.serviceType, *app.mam.ToOi4Identifier(), 1, app.GetHandler()); err != nil {
		return err
	}

	if err = app.registerPublications(); err != nil {
		return err
	}

	app.GetIntervalPublicationScheduler().Start()

	return nil
}

// Stop application and shutdown all publications and assets
func (app *Oi4ApplicationImpl) Stop() {
	for _, publication := range app.GetPublications() {
		publication.Stop()
	}
	app.sendGracefulShutdown()
	app.mqttClient.Stop()
}

// RegisterPublication Register a publisher for the specific application
// you can overwrite built-in publications like MAM, Health, etc...
func (app *Oi4ApplicationImpl) RegisterPublication(publication api.Publication) error {
	app.publicationMutex.Lock()
	defer app.publicationMutex.Unlock()

	resourcePublications := app.publications[publication.GetResource()]
	if resourcePublications == nil {
		resourcePublications = make([]api.Publication, 0)
	}

	found := false
	for i, current := range resourcePublications {
		if api.FilterEquals(current.GetFilter(), publication.GetFilter()) {
			resourcePublications[i] = publication
			found = true
			break
		}
	}

	if !found {
		resourcePublications = append(resourcePublications, publication)
	}

	app.publications[publication.GetResource()] = resourcePublications
	publication.Start()

	return nil
}

// GetPublications Return all registered publications
func (app *Oi4ApplicationImpl) GetPublications() []api.Publication {
	app.publicationMutex.RLock()
	defer app.publicationMutex.RUnlock()

	result := make([]api.Publication, 0)
	publications := slices.Collect(maps.Values(app.publications))
	for _, current := range publications {
		result = append(result, current...)
	}

	return result
}

// RegisterAsset Add new asset to the application
func (app *Oi4ApplicationImpl) RegisterAsset(asset *AssetImpl) {
	app.assetMutex.RLock()
	defer app.assetMutex.RUnlock()

	asset.setParent(app)
	oi4Id := asset.mam.ToOi4Identifier()
	app.assets[*oi4Id] = asset
}

// RemoveAsset remove an asset from the application
func (app *Oi4ApplicationImpl) RemoveAsset(asset *AssetImpl) {
	app.assetMutex.RLock()
	defer app.assetMutex.RUnlock()

	asset.setParent(nil)
	delete(app.assets, *asset.mam.ToOi4Identifier())
}

func (app *Oi4ApplicationImpl) UpdateHealth(health api.Health) {
	app.applicationSource.UpdateHealth(health)
}

func (app *Oi4ApplicationImpl) GetMam() *api.MasterAssetModel {
	return app.mam
}

func (app *Oi4ApplicationImpl) SendPublicationMessage(publication api.PublicationMessage) {
	if app.mqttClient == nil || publication.Content == nil || len(publication.Content) == 0 {
		return
	}

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

	topic := tp.NewTopic(
		app.serviceType,
		*app.mam.ToOi4Identifier(),
		api.MethodPub,
		publication.Resource,
		source,
		nil,
		publication.Filter,
	)

	err := app.mqttClient.PublishResource(topic.ToString(), opc.CreateNetworkMessage(app.mam.ToOi4Identifier(), app.serviceType, publication))
	if err != nil {
		return
	}
	app.logger.Debugf("Published message to topic: %s", topic.ToString())

}

func (app *Oi4ApplicationImpl) SendGetMessage(topic string, getMessage api.GetMessage) error {
	return app.mqttClient.PublishResource(topic, getMessage)
}

func (app *Oi4ApplicationImpl) GetHandler() api.MessageHandler {
	return subscription.NewMessageHandler(app, func(resource api.ResourceType, source *api.Oi4Identifier, networkMessage api.NetworkMessage, _ *tp.Topic) {
		sources := make([]api.BaseSource, 0)
		if source == nil {
			sources = append(sources, app.applicationSource)
			for _, asset := range app.assets {
				sources = append(sources, asset.source)
			}
		} else {
			if source.Equals(app.mam.ToOi4Identifier()) {
				sources = append(sources, app.applicationSource)
			} else if asset, ok := app.assets[*source]; ok {
				sources = append(sources, asset.source)
			}
		}

		if sources == nil || len(sources) == 0 {
			return
		}

		for _, current := range sources {
			var filter api.Filter
			// TODO retrieve the topic filter as input
			if len(networkMessage.Messages) > 0 {
				filter = networkMessage.Messages[0].Filter
			}

			app.triggerSourcePublication(current, resource, &filter, api.OnRequest, &networkMessage.MessageId)
		}
	})
}

func (app *Oi4ApplicationImpl) ResourceChanged(resource api.ResourceType, source api.BaseSource, filter *api.Filter) {
	app.triggerSourcePublication(source, resource, filter, api.OnRequest, nil)
}

func (app *Oi4ApplicationImpl) RegisterSubscription(subscription api.Subscription) error {
	app.subscriptionsMutex.Lock()
	defer app.subscriptionsMutex.Unlock()

	app.subscriptions[subscription.GetID()] = subscription

	return app.mqttClient.Subscribe(subscription)
}

func (app *Oi4ApplicationImpl) registerPublications() error {
	// register built-in publications
	err := app.RegisterPublication(pub.NewHealthPublication(app, app.applicationSource)) //

	//.SetDataFunc(func() *api.Health {health := application.applicationSource.GetHealth() return &health})
	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewMAMPublication(app, app.applicationSource))

	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewResourcePublication(app, app.applicationSource, api.ResourceLicense))
	//SetDataFunc(func() *api.License {// Dummy implementation yet		components := make([]api.LicenseComponent, 0)		return &api.License{			Components: components,		}	}).

	if err != nil {
		return err
	}

	for key, _ := range app.applicationSource.GetLicenseTexts() {
		err = app.RegisterPublication(pub.NewResourcePublicationWithFilter(app, app.applicationSource, api.ResourceLicenseText, api.NewFilter(key)))
		if err != nil {
			return err
		}
	}

	err = app.RegisterPublication(pub.NewResourcePublication(app, app.applicationSource, api.ResourcePublicationList))

	if err != nil {
		return err
	}

	err = app.RegisterPublication(pub.NewBuilder(app).
		Oi4Source(app.applicationSource).
		Resource(api.ResourceProfile).
		PublicationMode(api.PublicationMode_APPLICATION_2).
		Build())

	if err != nil {
		return err
	}

	return nil
}

func (app *Oi4ApplicationImpl) triggerSourcePublication(source api.BaseSource, resource api.ResourceType, filter *api.Filter, trigger api.Trigger, correlationId *string) {
	var publications []api.Publication
	if source.Equals(app.applicationSource) {
		publications = getPublications(app.publications, resource, filter)
	} else {
		var asset *AssetImpl

		for _, as := range app.assets {
			if as.source.Equals(source) {
				asset = as
			}
		}

		if asset == nil {
			return
		}

		publications = getPublications(asset.publications, resource, filter)
	}

	if publications == nil || len(publications) == 0 {
		return
	}

	for _, publication := range publications {
		publication.TriggerPublication(trigger, correlationId)
	}
}

func (app *Oi4ApplicationImpl) shouldPublicate(trigger api.Trigger, publication *api.PublicationList) bool {
	if trigger == api.OnRequest {
		return true
	}

	mode := *publication.Mode
	if mode == api.PublicationMode_OFF_0 || //
		mode == api.PublicationMode_ON_REQUEST_1 {
		return false
	}

	interval := *publication.Interval
	if interval == 0 && trigger != api.ByInterval || //
		interval != 0 && trigger == api.ByInterval {
		return true
	}

	return false
}

func (app *Oi4ApplicationImpl) sendGracefulShutdown() {
	code := api.Status_Good
	app.SendPublicationMessage(api.PublicationMessage{
		Resource: api.ResourceHealth,
		Source:   app.mam.ToOi4Identifier(),
		Content: []api.PublicationContent{
			{
				StatusCode: &code,
				Data:       &api.Health{Health: api.Health_Normal, HealthScore: 0},
			},
		},
		//publicationMode: api.PublicationMode_APPLICATION_SOURCE_5,
	})
}

func (app *Oi4ApplicationImpl) newMqttClient(options *api.MqttClientOptions) (api.MqttClient, error) {
	if app.createMqttClientFn != nil {
		return app.createMqttClientFn(options)
	}
	return mqtt.NewClient(options)
}

func getPublications(publications map[api.ResourceType][]api.Publication, resource api.ResourceType, filter *api.Filter) []api.Publication {
	resourcePublications := publications[resource]
	if resourcePublications == nil || filter == nil {
		return resourcePublications
	}

	for _, current := range resourcePublications {
		if api.FilterEquals(current.GetFilter(), filter) {
			return []api.Publication{current}
		}
	}

	return nil
}

/******************************************************
**Builder Option for creating a new Oi4ApplicationImpl  **
******************************************************/

type Option func(app *Oi4ApplicationImpl)

func WithMqttClientFn(fn func(options *api.MqttClientOptions) (api.MqttClient, error)) Option {
	return func(app *Oi4ApplicationImpl) {
		app.createMqttClientFn = fn
	}
}
