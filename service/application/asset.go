package application

import (
	"sync"
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
)

type Oi4Asset struct {
	parent *Oi4ApplicationImpl
	mam    *api.MasterAssetModel

	publicationsList map[api.ResourceType]Publication
	publicationMutex sync.RWMutex

	source api.Oi4Source
}

func CreateNewAsset(source api.Oi4Source) *Oi4Asset {
	mam := source.GetMasterAssetModel()
	asset := &Oi4Asset{
		parent:           nil,
		mam:              &mam,
		publicationsList: make(map[api.ResourceType]Publication),
		publicationMutex: sync.RWMutex{},
		source:           source,
	}

	assetSource := NewSourceImpl(mam)
	oi4Source := api.Oi4Source(assetSource)

	asset.RegisterPublication(CreatePublication(api.ResourceHealth, &oi4Source). //SetDataFunc(func() *api.Health {
		//	health := asset.source.GetHealth()
		//	return &health
		//}).
		SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5).SetPublicationInterval(60 * time.Second))
	asset.RegisterPublication(CreatePublication(api.ResourceMam, &oi4Source). //SetData(&mam).
											SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5))

	asset.RegisterPublication(CreatePublication(api.ResourceReferenceDesignation, &oi4Source). //SetDataFunc(func() *api.ReferenceDesignation {
		// Dummy implementation yet
		//	return &api.ReferenceDesignation{}
		//}).
		SetPublicationMode(api.PublicationMode_APPLICATION_SOURCE_5))

	asset.RegisterPublication(CreatePublication(api.ResourceProfile, &oi4Source). //SetDataFunc(func() *api.Profile {
		//	resources := make([]api.ResourceType, 0)
		//	for key := range asset.publicationsList {
		//		resources = append(resources, key)
		//	}
		//	profile := api.Profile{
		//		Resources: resources,
		//	}
		//	return &profile
		//}).
		SetPublicationMode(api.PublicationMode_APPLICATION_2))

	return asset
}

func (asset *Oi4Asset) RegisterPublication(publication Publication) error {
	asset.publicationMutex.Lock()
	defer asset.publicationMutex.Unlock()

	if publication.getParent() != nil {
		return ErrPublicationAlreadyRegisteredOnAnotherApplication
	}

	asset.publicationsList[publication.getResource()] = publication

	if asset.parent != nil {
		publication.setParent(asset.parent)
		publication.setSource(asset.mam.ToOi4Identifier())
		publication.start()
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
	for _, pub := range asset.publicationsList {
		if parent != nil {
			pub.setParent(asset.parent)
			pub.setSource(asset.mam.ToOi4Identifier())
			pub.start()
		} else {
			pub.setSource(nil)
			pub.stop()
		}
	}
	return nil
}
