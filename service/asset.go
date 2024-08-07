package service

import (
	"sync"
	"time"

	oi4 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
)

type Oi4Asset struct {
	parent *Oi4Application
	mam    *oi4.MasterAssetModel

	publicationsList map[oi4.ResourceType]Publication
	publicationMutex sync.RWMutex
}

func CreateNewAsset(mam *oi4.MasterAssetModel) *Oi4Asset {
	asset := &Oi4Asset{
		parent:           nil,
		mam:              mam,
		publicationsList: make(map[oi4.ResourceType]Publication),
		publicationMutex: sync.RWMutex{},
	}

	asset.RegisterPublication(CreatePublication[*oi4.Health](oi4.ResourceHealth, true).SetData(&oi4.Health{Health: oi4.Health_Normal, HealthScore: 100}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5).SetPublicationInterval(60 * time.Second))
	asset.RegisterPublication(CreatePublication[*oi4.MasterAssetModel](oi4.ResourceMam, true).SetData(mam).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))

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

func (asset *Oi4Asset) setParent(parent *Oi4Application) error {
	if asset.parent != nil && parent != nil {
		return ErrAssetAlreadyRegistered
	}

	if parent != nil {
		asset.parent = parent
		for _, pub := range asset.publicationsList {
			pub.setParent(asset.parent)
			pub.setSource(asset.mam.ToOi4Identifier())
			pub.start()
		}
	} else {
		asset.parent = nil
		for _, pub := range asset.publicationsList {
			pub.setSource(nil)
			pub.stop()
		}

	}
	return nil
}
