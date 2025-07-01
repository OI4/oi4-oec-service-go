package api

import (
	"go.uber.org/zap"
)

type Oi4Application interface {
	GetServiceType() ServiceType
	GetOi4Identifier() Oi4Identifier
	GetApplicationSource() ApplicationSource
	GetLogger() *zap.SugaredLogger
	ResourceChanged(resource ResourceType, source BaseSource, filter *Filter)
	SendPublicationMessage(publication PublicationMessage)
	SendGetMessage(topic string, getMessage GetMessage) error
	GetIntervalPublicationScheduler() IntervalPublicationScheduler

	PublicationProvider
}
