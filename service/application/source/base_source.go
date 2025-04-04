package source

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"maps"
	"slices"
)

type BaseSourceImpl struct {
	profile              api.Profile
	mam                  api.MasterAssetModel
	health               api.Health
	config               api.PublishConfig
	license              api.License
	licenseText          map[string]api.LicenseText
	rtLicense            api.RtLicense
	subscriptionList     []api.SubscriptionList
	referenceDesignation api.ReferenceDesignation
	data                 map[string]api.Data

	application api.Oi4Application

	publicationProvider api.PublicationProvider

	dataFn        func(source api.BaseSource, filter *api.Filter) []api.Data
	dataWrapperFn func([]api.Data) []any
	healthFn      func(source api.BaseSource) api.Health
}

func newBaseImpl(mam api.MasterAssetModel, options ...Option) *BaseSourceImpl {
	source := &BaseSourceImpl{
		mam:                  mam,
		health:               api.Health{Health: api.Health_Normal, HealthScore: 100},
		config:               api.PublishConfig{},
		license:              api.EmptyLicense(),
		licenseText:          make(map[string]api.LicenseText),
		rtLicense:            api.RtLicense{},
		subscriptionList:     make([]api.SubscriptionList, 0),
		referenceDesignation: api.ReferenceDesignation{},
		data:                 make(map[string]api.Data),
	}

	// Apply all the functional options to configure the client.
	for _, opt := range options {
		opt(source)
	}

	return source
}

func (source *BaseSourceImpl) GetOi4Identifier() *api.Oi4Identifier {
	return source.mam.ToOi4Identifier()
}

func (source *BaseSourceImpl) Equals(other api.BaseSource) bool {
	return source.GetOi4Identifier().Equals(other.GetOi4Identifier())
}

func (source *BaseSourceImpl) GetProfile() api.Profile {
	return source.profile
}

func (source *BaseSourceImpl) GetMasterAssetModel() api.MasterAssetModel {
	return source.mam
}

func (source *BaseSourceImpl) GetHealth() api.Health {
	if source.healthFn != nil {
		return source.healthFn(source)
	}
	return source.health
}

func (source *BaseSourceImpl) UpdateHealth(health api.Health) {
	source.health = health
	if source.application != nil {
		source.application.ResourceChanged(api.ResourceHealth, source, nil)
	}
}

func (source *BaseSourceImpl) GetData(filter *api.Filter) []api.Data {
	if source.dataFn != nil {
		return source.dataFn(source, filter)
	}
	if filter == nil {
		return slices.Collect(maps.Values(source.data))
	}
	if data, ok := source.data[filter.String()]; ok {
		return []api.Data{data}
	}
	return nil
}

func (source *BaseSourceImpl) UpdateData(data api.Data, dataTag string) {
	source.data[dataTag] = data
	if source.application != nil {
		source.application.ResourceChanged(api.ResourceData, source, api.NewFilter(dataTag))
	}
}

func (source *BaseSourceImpl) GetConfig() api.PublishConfig {
	return source.config
}

func (source *BaseSourceImpl) GetLicense() api.License {
	return source.license
}

func (source *BaseSourceImpl) GetLicenseText(filter *api.Filter) []api.LicenseText {
	if len(source.licenseText) == 0 || filter == nil {
		return nil
	}

	licenseText, ok := source.licenseText[filter.String()]
	if ok {
		return []api.LicenseText{licenseText}
	}

	return nil
}

func (source *BaseSourceImpl) GetLicenseTexts() map[string]api.LicenseText {
	if len(source.licenseText) == 0 {
		return make(map[string]api.LicenseText)
	}

	return source.licenseText
}

func (source *BaseSourceImpl) GetRtLicense() api.RtLicense {
	return source.rtLicense
}

func (source *BaseSourceImpl) GetPublicationList() []api.PublicationList {
	srcPublications := source.publicationProvider.GetPublications()
	publications := make([]api.PublicationList, len(srcPublications))
	for i, pub := range srcPublications {
		var filter *api.Filter
		if pub.GetFilter() != nil {
			filter = pub.GetFilter()
		}
		publications[i] = api.PublicationList{
			ResourceType:    pub.GetResource(),
			Source:          pub.GetSource().ToString(),
			Filter:          filter,
			DataSetWriterId: pub.GetDataSetWriterId(),
			Mode:            pub.GetPublicationMode(),
			//Interval:        pub.GetInterval(),
			//Precisions:      pub.GetPrecisions(),
			//Config:          pub.GetConfig(),
		}
	}
	return publications
}

func (source *BaseSourceImpl) GetSubscriptionList() []api.SubscriptionList {
	return source.subscriptionList
}

func (source *BaseSourceImpl) GetReferenceDesignation() api.ReferenceDesignation {
	return source.referenceDesignation
}

func (source *BaseSourceImpl) Get(resourceType api.ResourceType, filter *api.Filter) []any {
	getResource := func() any {
		switch resourceType {
		case api.ResourceProfile:
			return source.GetProfile()
		case api.ResourceMam:
			return source.GetMasterAssetModel()
		case api.ResourceHealth:
			return source.GetHealth()
		case api.ResourceConfig:
			return source.GetConfig()
		case api.ResourceLicense:
			return source.GetLicense()
		case api.ResourceLicenseText:
			return source.GetLicenseText(filter)
		case api.ResourceRtLicense:
			return source.GetRtLicense()
		case api.ResourcePublicationList:
			return source.GetPublicationList()
		case api.ResourceSubscriptionList:
			return source.GetSubscriptionList()
		case api.ResourceReferenceDesignation:
			return source.GetReferenceDesignation()
		case api.ResourceData:
			return source.wrapData(source.GetData(filter))
		default:
			return nil
		}
	}

	return toAnySlice(getResource())
}

func (source *BaseSourceImpl) wrapData(data []api.Data) []any {
	if data == nil || len(data) == 0 {
		return nil
	}

	if source.dataWrapperFn != nil {
		return source.dataWrapperFn(data)
	}

	result := make([]any, len(data))
	for i, v := range data {
		result[i] = v.GetData()
	}
	return result
}

func toAnySlice[T any](input ...T) []any {
	if input == nil || len(input) == 0 {
		return nil
	}
	result := make([]any, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}

func (source *BaseSourceImpl) SetOi4Application(application api.Oi4Application) {
	source.application = application
}

type Option func(*BaseSourceImpl)

func WithHealthFn(fn func(source api.BaseSource) api.Health) Option {
	return func(s *BaseSourceImpl) {
		s.healthFn = fn
	}
}

func WithDataFn(fn func(source api.BaseSource, filter *api.Filter) []api.Data) Option {
	return func(s *BaseSourceImpl) {
		s.dataFn = fn
	}
}

func WithDataWrapperFn(fn func([]api.Data) []any) Option {
	return func(s *BaseSourceImpl) {
		s.dataWrapperFn = fn
	}
}

func WithProfile(profile api.Profile) Option {
	return func(s *BaseSourceImpl) {
		s.profile = profile
	}
}

func WithConfig(config api.PublishConfig) Option {
	return func(s *BaseSourceImpl) {
		s.config = config
	}
}

func WithLicense(license api.License) Option {
	return func(s *BaseSourceImpl) {
		s.license = license
	}
}

func WithLicenseText(licenseText map[string]api.LicenseText) Option {
	return func(s *BaseSourceImpl) {
		s.licenseText = licenseText
	}
}

func WithRtLicense(rtLicense api.RtLicense) Option {
	return func(s *BaseSourceImpl) {
		s.rtLicense = rtLicense
	}
}

func WithReferenceDesignation(ref api.ReferenceDesignation) Option {
	return func(s *BaseSourceImpl) {
		s.referenceDesignation = ref
	}
}
