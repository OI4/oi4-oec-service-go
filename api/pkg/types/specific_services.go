package types

type SpecificService string

const (
	Read          SpecificService = "Read"
	Write         SpecificService = "Write"
	Subscribe     SpecificService = "Subscribe"
	Unsubscribe   SpecificService = "Unsubscribe"
	GenericMethod SpecificService = "GenericMethod"
)
