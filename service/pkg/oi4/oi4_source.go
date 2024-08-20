package oi4

import v1 "github.com/OI4/oi4-oec-service-go/api/pkg/types"

type SourceImpl struct {
	profile              v1.Profile
	mam                  v1.MasterAssetModel
	health               v1.Health
	config               *v1.PublishConfig
	license              *[]v1.License
	licenseText          *map[string]v1.LicenseText
	rtLicense            *v1.RTLicense
	publicationList      []v1.PublicationList
	subscriptionList     []v1.SubscriptionList
	referenceDesignation v1.ReferenceDesignation
	data                 *any

	application v1.Oi4Application
}

func NewSourceImpl(mam v1.MasterAssetModel) *SourceImpl {
	return &SourceImpl{
		mam:    mam,
		health: v1.Health{Health: v1.Health_Normal, HealthScore: 100},
	}
}

func (source *SourceImpl) GetOi4Identifier() *v1.Oi4Identifier {
	return source.mam.ToOi4Identifier()
}

func (source *SourceImpl) GetProfile() v1.Profile {
	return source.profile
}

func (source *SourceImpl) GetMasterAssetModel() v1.MasterAssetModel {
	return source.mam
}

func (source *SourceImpl) GetHealth() v1.Health {
	return source.health
}

func (source *SourceImpl) UpdateHealth(health v1.Health) {
	source.health = health
	if source.application != nil {
		source.application.ResourceChanged(v1.ResourceHealth, source, nil)
	}
}

func (source *SourceImpl) GetData() *any {
	return source.data
}

func (source *SourceImpl) UpdateData(data *any, dataTag string) {
	source.data = data
	if source.application != nil {
		source.application.ResourceChanged(v1.ResourceData, source, &dataTag)
	}
}

func (source *SourceImpl) SetOi4Application(application v1.Oi4Application) {
	source.application = application
}
