package source

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"sync"
)

type ApplicationSourceImpl struct {
	BaseSourceImpl

	sources     map[api.Oi4Identifier]*api.AssetSource
	sourceMutex sync.RWMutex
}

func NewApplicationSourceImpl(mam api.MasterAssetModel, options ...Option) *ApplicationSourceImpl {
	options = append([]Option{WithProfile(api.ProfileApplication())}, options...)
	source := ApplicationSourceImpl{
		BaseSourceImpl: *newBaseImpl(mam, options...),
		sources:        make(map[api.Oi4Identifier]*api.AssetSource),
		sourceMutex:    sync.RWMutex{},
	}

	source.publicationProvider = &source

	return &source
}

func (source *ApplicationSourceImpl) GetSources() map[api.Oi4Identifier]*api.AssetSource {
	return source.sources
}

func (source *ApplicationSourceImpl) AddSource(sourceToAdd api.AssetSource) {
	source.sourceMutex.RLock()
	defer source.sourceMutex.RUnlock()

	source.sources[*sourceToAdd.GetOi4Identifier()] = &sourceToAdd
}

func (source *ApplicationSourceImpl) RemoveSource(sourceToRemove api.Oi4Identifier) {
	source.sourceMutex.RLock()
	defer source.sourceMutex.RUnlock()

	delete(source.sources, sourceToRemove)
}

func (source *ApplicationSourceImpl) GetPublications() []api.Publication {
	return source.application.GetPublications()
}
