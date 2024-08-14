package oi4

import (
	v1 "github.com/OI4/oi4-oec-service-go/api/pkg/types"
	"sync"
)

type ApplicationSourceImpl struct {
	SourceImpl

	sources     map[v1.Oi4Identifier]*v1.Oi4Source
	sourceMutex sync.RWMutex
}

func NewApplicationSourceImpl(mam v1.MasterAssetModel) *ApplicationSourceImpl {
	return &ApplicationSourceImpl{
		SourceImpl: SourceImpl{
			mam:    mam,
			health: v1.Health{Health: v1.Health_Normal, HealthScore: 100},
		},
		sources:     make(map[v1.Oi4Identifier]*v1.Oi4Source),
		sourceMutex: sync.RWMutex{},
	}
}

func (source *ApplicationSourceImpl) GetSources() map[v1.Oi4Identifier]*v1.Oi4Source {
	return source.sources
}

func (source *ApplicationSourceImpl) AddSource(sourceToAdd v1.Oi4Source) {
	source.sourceMutex.RLock()
	defer source.sourceMutex.RUnlock()

	source.sources[*sourceToAdd.GetOi4Identifier()] = &sourceToAdd
}

func (source *ApplicationSourceImpl) RemoveSource(sourceToRemove v1.Oi4Identifier) {
	source.sourceMutex.RLock()
	defer source.sourceMutex.RUnlock()

	delete(source.sources, sourceToRemove)
}
