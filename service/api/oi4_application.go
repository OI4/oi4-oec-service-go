package api

import "go.uber.org/zap"

type Oi4Application interface {
	GetLogger() *zap.SugaredLogger
	ResourceChanged(resource ResourceType, source Oi4Source, resourceTag *string)
	SendPublicationMessage(publication PublicationMessage)
}

type PublicationMessage struct {
	Resource      ResourceType
	Source        *Oi4Identifier
	CorrelationId string
	//publicationMode api.PublicationMode
	Data       interface{}
	StatusCode StatusCode
	Filter     Filter
}
