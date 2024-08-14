package types

type Oi4Application interface {
	ResourceChanged(resource ResourceType, source Oi4Source)
}
