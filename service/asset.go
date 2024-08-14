package service

import (
	"sync"
	"time"

	oi4 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
)

type Oi4Asset struct {
	parent *Oi4ApplicationImpl
	mam    *oi4.MasterAssetModel

	publicationsList map[oi4.ResourceType]Publication
	publicationMutex sync.RWMutex

	source oi4.Oi4Source
}

func CreateNewAsset(source oi4.Oi4Source) *Oi4Asset {
	mam := source.GetMasterAssetModel()
	asset := &Oi4Asset{
		parent:           nil,
		mam:              &mam,
		publicationsList: make(map[oi4.ResourceType]Publication),
		publicationMutex: sync.RWMutex{},
		source:           source,
	}

	asset.RegisterPublication(CreatePublication[*oi4.Health](oi4.ResourceHealth, true).SetDataFunc(func() *oi4.Health {
		health := asset.source.GetHealth()
		return &health
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5).SetPublicationInterval(60 * time.Second))
	asset.RegisterPublication(CreatePublication[*oi4.MasterAssetModel](oi4.ResourceMam, true).SetData(&mam).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))

	asset.RegisterPublication(CreatePublication[*oi4.ReferenceDesignation](oi4.ResourceReferenceDesignation, false).SetDataFunc(func() *oi4.ReferenceDesignation {
		// Dummy implementation yet
		return &oi4.ReferenceDesignation{}
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))

	asset.RegisterPublication(CreatePublication[*oi4.Profile](oi4.ResourceProfile, false).SetDataFunc(func() *oi4.Profile {
		resources := make([]oi4.ResourceType, 0)
		for key := range asset.publicationsList {
			resources = append(resources, key)
		}
		profile := oi4.Profile{
			Resources: resources,
		}
		return &profile
	}).SetPublicationMode(oi4.PublicationMode_APPLICATION_2))

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

func (asset *Oi4Asset) UpdateHealth(health oi4.Health) {
	asset.publicationsList[oi4.ResourceHealth].(*PublicationImpl[*oi4.Health]).SetData(&health)
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
