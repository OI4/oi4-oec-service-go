package service

import (
	"sync"
	"time"

	oi4 "github.com/mzeiher/oi4/api/pkg/types"
)

type Oi4Asset struct {
	parent *Oi4Application
	mam    *oi4.MasterAssetModel

	publicationsList map[oi4.Resource]Publication
	publicationMutex sync.RWMutex
}

func CreateNewAsset(mam *oi4.MasterAssetModel) *Oi4Asset {
	asset := &Oi4Asset{
		parent:           nil,
		mam:              mam,
		publicationsList: make(map[oi4.Resource]Publication),
		publicationMutex: sync.RWMutex{},
	}

	asset.RegisterPublication(CreatePublication[*oi4.Health](oi4.Resource_Health, true).SetData(&oi4.Health{Health: oi4.Health_Normal, HealthScore: 100}).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5).SetPublicationInterval(60 * time.Second))
	asset.RegisterPublication(CreatePublication[*oi4.MasterAssetModel](oi4.Resource_MAM, true).SetData(mam).SetPublicationMode(oi4.PublicationMode_APPLICATION_SOURCE_5))

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
	asset.publicationsList[oi4.Resource_Health].(*PublicationImpl[*oi4.Health]).SetData(&health)
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
