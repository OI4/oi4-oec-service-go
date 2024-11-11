package api

import "time"

// PublicationType type of trigger for publication
type PublicationType int

// Declare typed constants each with type of status
const (
	General PublicationType = iota
	Interval
)

// Trigger type of trigger for publication
type Trigger int

// Declare typed constants each with type of status
const (
	ByInterval Trigger = iota
	OnRequest
)

type Publication interface {
	GetApplication() Oi4Application
	GetPublicationType() PublicationType
	ShouldPublicate(trigger Trigger) bool
	GetOi4Source() BaseSource
	GetResource() ResourceType
	GetSource() *Oi4Identifier
	GetFilter() Filter
	GetID() string
	GetDataSetWriterId() uint16
	GetPublicationMode() *PublicationMode

	Stop()
	Start()

	TriggerPublication(trigger Trigger, correlationId *string) bool

	//triggerSourcePublication(byInterval bool, onRequest bool, correlationId string)
	publishOnRegistration() bool
}

type IntervalPublication interface {
	Publication
	GetNextPublicationTime() time.Time
	DueForPublication() bool
}

type PublicationMessage struct {
	Resource      ResourceType
	Source        *Oi4Identifier
	CorrelationId *string
	//publicationMode api.PublicationMode
	//Data       any
	//StatusCode StatusCode
	Filter
	Content []PublicationContent
}

type PublicationContent struct {
	*StatusCode
	Data any
}

type IntervalPublicationScheduler interface {
	Start()
	Stop()
	AddPublication(publication IntervalPublication)
	RemovePublication(publication IntervalPublication)
}

type IntervalPublicationWorker interface {
}

type PublicationProvider interface {
	GetPublications() []Publication
}
