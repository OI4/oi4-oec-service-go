package application

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"sync"
)

type ApplicationSourceImpl struct {
	SourceImpl

	sources     map[api.Oi4Identifier]*api.Oi4Source
	sourceMutex sync.RWMutex
}

func NewApplicationSourceImpl(mam api.MasterAssetModel) *ApplicationSourceImpl {
	return &ApplicationSourceImpl{
		SourceImpl: SourceImpl{
			mam:    mam,
			health: api.Health{Health: api.Health_Normal, HealthScore: 100},
		},
		sources:     make(map[api.Oi4Identifier]*api.Oi4Source),
		sourceMutex: sync.RWMutex{},
	}
}

func (source *ApplicationSourceImpl) GetSources() map[api.Oi4Identifier]*api.Oi4Source {
	return source.sources
}

func (source *ApplicationSourceImpl) AddSource(sourceToAdd api.Oi4Source) {
	source.sourceMutex.RLock()
	defer source.sourceMutex.RUnlock()

	source.sources[*sourceToAdd.GetOi4Identifier()] = &sourceToAdd
}

func (source *ApplicationSourceImpl) RemoveSource(sourceToRemove api.Oi4Identifier) {
	source.sourceMutex.RLock()
	defer source.sourceMutex.RUnlock()

	delete(source.sources, sourceToRemove)
}
