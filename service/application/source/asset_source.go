package source

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
)

type AssetSourceImpl struct {
	BaseSourceImpl
	asset api.Asset
}

func NewAssetSourceImpl(mam api.MasterAssetModel, options ...Option) *AssetSourceImpl {
	options = append([]Option{WithProfile(api.ProfileDevice())}, options...)
	source := AssetSourceImpl{
		BaseSourceImpl: *newBaseImpl(mam, options...),
	}

	source.publicationProvider = &source

	return &source
}

func (source *AssetSourceImpl) GetPublications() []api.Publication {
	return source.asset.GetPublications()
}
