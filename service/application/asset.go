package application

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	pub "github.com/OI4/oi4-oec-service-go/service/application/publication"
	"maps"
	"slices"
	"sync"
)

type AssetImpl struct {
	parent *Oi4ApplicationImpl
	mam    *api.MasterAssetModel

	publications     map[api.ResourceType][]api.Publication
	publicationMutex sync.RWMutex

	source api.AssetSource
}

func CreateNewAsset(source api.AssetSource, app *Oi4ApplicationImpl) *AssetImpl {
	mam := source.GetMasterAssetModel()
	asset := &AssetImpl{
		parent:           app,
		mam:              &mam,
		publications:     make(map[api.ResourceType][]api.Publication),
		publicationMutex: sync.RWMutex{},
		source:           source,
	}

	source.SetAsset(asset)

	err := asset.RegisterPublication(pub.NewHealthPublication(app, source))

	if err != nil {
		return nil
	}

	err = asset.RegisterPublication(pub.NewMAMPublication(app, source))

	if err != nil {
		return nil
	}

	err = asset.RegisterPublication(pub.NewResourcePublication(app, source, api.ResourceReferenceDesignation))

	if err != nil {
		return nil
	}

	err = asset.RegisterPublication(pub.NewResourcePublication(app, source, api.ResourcePublicationList))

	if err != nil {
		return nil
	}

	err = asset.RegisterPublication(pub.NewBuilder(app). //
								Oi4Source(source).                                  //
								Resource(api.ResourceProfile).                      //
								PublicationMode(api.PublicationMode_APPLICATION_2). //
								Build())

	if err != nil {
		return nil
	}

	return asset
}

func (asset *AssetImpl) RegisterPublication(publication api.Publication) error {
	asset.publicationMutex.Lock()
	defer asset.publicationMutex.Unlock()

	resourcePublications := asset.publications[publication.GetResource()]
	if resourcePublications == nil {
		resourcePublications = make([]api.Publication, 0)
	}

	found := false
	for i, current := range resourcePublications {
		if current.GetFilter().Equals(publication.GetFilter()) {
			resourcePublications[i] = publication
			found = true
			break
		}
	}

	if !found {
		resourcePublications = append(resourcePublications, publication)
	}

	asset.publications[publication.GetResource()] = resourcePublications

	if asset.parent != nil {
		publication.Start()
	}

	return nil
}

// GetPublications Return all registered publications
func (asset *AssetImpl) GetPublications() []api.Publication {
	asset.publicationMutex.RLock()
	defer asset.publicationMutex.RUnlock()

	result := make([]api.Publication, 0)
	publications := slices.Collect(maps.Values(asset.publications))
	for _, current := range publications {
		result = append(result, current...)
	}

	return result
}

func (asset *AssetImpl) UpdateHealth(health api.Health) {
	asset.source.UpdateHealth(health)
}

func (asset *AssetImpl) setParent(parent *Oi4ApplicationImpl) {
	asset.parent = parent
	asset.source.SetOi4Application(parent)
	for _, publication := range asset.publications {
		for _, current := range publication {
			if parent != nil {
				current.Start()
			} else {
				current.Stop()
			}
		}
	}
}
