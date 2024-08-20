package types

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

	GetData() *any
	UpdateData(data *any, dataTag string)

	SetOi4Application(Oi4Application)
}
