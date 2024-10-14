package application

import "github.com/OI4/oi4-oec-service-go/service/api"

type SourceImpl struct {
	profile              api.Profile
	mam                  api.MasterAssetModel
	health               api.Health
	config               *api.PublishConfig
	license              *[]api.License
	licenseText          map[string]api.LicenseText
	rtLicense            *api.RtLicense
	publicationList      []api.PublicationList
	subscriptionList     []api.SubscriptionList
	referenceDesignation api.ReferenceDesignation
	data                 api.Data

	application api.Oi4Application
}

func NewSourceImpl(mam api.MasterAssetModel) *SourceImpl {
	return &SourceImpl{
		mam:         mam,
		health:      api.Health{Health: api.Health_Normal, HealthScore: 100},
		licenseText: make(map[string]api.LicenseText),
	}
}

func (source *SourceImpl) GetOi4Identifier() *api.Oi4Identifier {
	return source.mam.ToOi4Identifier()
}

func (source *SourceImpl) GetProfile() api.Profile {
	return source.profile
}

func (source *SourceImpl) GetMasterAssetModel() api.MasterAssetModel {
	return source.mam
}

func (source *SourceImpl) GetHealth() api.Health {
	return source.health
}

func (source *SourceImpl) UpdateHealth(health api.Health) {
	source.health = health
	if source.application != nil {
		source.application.ResourceChanged(api.ResourceHealth, source, nil)
	}
}

func (source *SourceImpl) GetData() api.Data {
	return source.data
}

func (source *SourceImpl) UpdateData(data api.Data, dataTag string) {
	source.data = data
	if source.application != nil {
		source.application.ResourceChanged(api.ResourceData, source, &dataTag)
	}
}

func (source *SourceImpl) GetConfig() *api.PublishConfig {
	return source.config
}

func (source *SourceImpl) GetLicense() *[]api.License {
	return source.license
}

func (source *SourceImpl) GetLicenseText() map[string]api.LicenseText {
	return source.licenseText
}

func (source *SourceImpl) GetRtLicense() *api.RtLicense {
	return source.rtLicense
}

func (source *SourceImpl) GetPublicationList() []api.PublicationList {
	return source.publicationList
}

func (source *SourceImpl) GetSubscriptionList() []api.SubscriptionList {
	return source.subscriptionList
}

func (source *SourceImpl) GetReferenceDesignation() api.ReferenceDesignation {
	return source.referenceDesignation
}

func (source *SourceImpl) Get(resourceType api.ResourceType) any {
	switch resourceType {
	case api.ResourceProfile:
		return source.profile
	case api.ResourceMam:
		return source.GetMasterAssetModel()
	case api.ResourceHealth:
		return source.GetHealth()
	case api.ResourceConfig:
		return source.GetConfig()
	case api.ResourceLicense:
		return source.GetLicense()
	case api.ResourceLicenseText:
		return source.licenseText
	case api.ResourceRtLicense:
		return source.rtLicense
	case api.ResourcePublicationList:
		return source.publicationList
	case api.ResourceSubscriptionList:
		return source.subscriptionList
	case api.ResourceReferenceDesignation:
		return source.referenceDesignation
	case api.ResourceData:
		return source.GetData()
	default:
		return nil
	}
}

func (source *SourceImpl) SetOi4Application(application api.Oi4Application) {
	source.application = application
}
