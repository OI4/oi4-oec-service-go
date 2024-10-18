package application

import (
	pub "github.com/OI4/oi4-oec-service-go/service/application/publication"
	"sync"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
)

type Oi4Asset struct {
	parent *Oi4ApplicationImpl
	mam    *api.MasterAssetModel

	publicationsList map[api.ResourceType]pub.Publication
	publicationMutex sync.RWMutex

	source api.Oi4Source
}

func CreateNewAsset(source api.Oi4Source, app *Oi4ApplicationImpl) *Oi4Asset {
	mam := source.GetMasterAssetModel()
	asset := &Oi4Asset{
		parent:           app,
		mam:              &mam,
		publicationsList: make(map[api.ResourceType]pub.Publication),
		publicationMutex: sync.RWMutex{},
		source:           source,
	}

	assetSource := NewSourceImpl(mam)

	// register built-in publications
	err := asset.RegisterPublication(pub.NewIntervalBuilder(app, 60*time.Second). //
											Oi4Source(assetSource).                                    //
											Resource(api.ResourceHealth).                              //
											PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
											Build())
	//SetDataFunc(func() *api.Health {
	//	health := asset.source.GetHealth()
	//	return &health
	//}).

	if err != nil {
		return nil
	}

	err = asset.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(assetSource).                                    //
								Resource(api.ResourceMam).                                 //
								PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
								PublishOnRegistration(true).                               //
								Build())

	if err != nil {
		return nil
	}

	err = asset.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(assetSource).                                    //
								Resource(api.ResourceReferenceDesignation).                //
								PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
								Build())
	//SetDataFunc(func() *api.ReferenceDesignation {
	// Dummy implementation yet
	//	return &api.ReferenceDesignation{}
	//}).

	if err != nil {
		return nil
	}

	err = asset.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(assetSource).                             //
								Resource(api.ResourceProfile).                      //
								PublicationMode(api.PublicationMode_APPLICATION_2). //
								Build())
	//SetDataFunc(func() *api.Profile {
	//	resources := make([]api.ResourceType, 0)
	//	for key := range asset.publicationsList {
	//		resources = append(resources, key)
	//	}
	//	profile := api.Profile{
	//		Resources: resources,
	//	}
	//	return &profile
	//}).

	if err != nil {
		return nil
	}

	return asset
}

func (asset *Oi4Asset) RegisterPublication(publication pub.Publication) error {
	asset.publicationMutex.Lock()
	defer asset.publicationMutex.Unlock()

	asset.publicationsList[publication.GetResource()] = publication

	if asset.parent != nil {
		publication.Start()
	}

	return nil
}

func (asset *Oi4Asset) UpdateHealth(health api.Health) {
	// TODO
	// asset.publicationsList[api.ResourceHealth].(*PublicationImpl).SetData(&health)
}

func (asset *Oi4Asset) setParent(parent *Oi4ApplicationImpl) error {
	if asset.parent != nil && parent != nil {
		return ErrAssetAlreadyRegistered
	}

	asset.parent = parent
	asset.source.SetOi4Application(parent)
	for _, publication := range asset.publicationsList {
		if parent != nil {
			publication.Start()
		} else {
			publication.Stop()
		}
	}
	return nil
}
