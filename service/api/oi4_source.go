package api

type Oi4ApplicationSource interface {
	Oi4Source

	GetSources() map[Oi4Identifier]*Oi4Source
	AddSource(Oi4Source)
	RemoveSource(Oi4Identifier)
}

type Oi4Source interface {
	GetOi4Identifier() *Oi4Identifier

	GetMasterAssetModel() MasterAssetModel

	GetHealth() Health
	UpdateHealth(Health)

	GetData() Data
	UpdateData(data Data, dataTag string)

	GetConfig() *PublishConfig

	GetProfile() Profile

	GetLicense() *[]License

	GetLicenseText() map[string]LicenseText

	GetRtLicense() *RtLicense

	GetPublicationList() []PublicationList

	GetSubscriptionList() []SubscriptionList

	GetReferenceDesignation() ReferenceDesignation

	Get(ResourceType) any

	SetOi4Application(Oi4Application)
}
