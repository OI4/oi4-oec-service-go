package api

type ApplicationSource interface {
	BaseSource

	GetSources() map[Oi4Identifier]*AssetSource
	AddSource(AssetSource)
	RemoveSource(Oi4Identifier)
}

type AssetSource interface {
	BaseSource

	SetAsset(asset Asset)
}

type BaseSource interface {
	GetOi4Identifier() *Oi4Identifier

	GetMasterAssetModel() MasterAssetModel

	GetHealth() Health
	UpdateHealth(Health)

	GetData() Data
	UpdateData(data Data, dataTag string)

	GetConfig() PublishConfig

	GetProfile() Profile

	GetLicense() License

	GetLicenseText(filter Filter) []LicenseText

	GetRtLicense() RtLicense

	GetPublicationList() []PublicationList

	GetSubscriptionList() []SubscriptionList

	GetReferenceDesignation() ReferenceDesignation

	Get(resource ResourceType, filter Filter) []any

	SetOi4Application(Oi4Application)

	Equals(BaseSource) bool
}
