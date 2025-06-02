package api

import "go.uber.org/zap"

type Oi4Application interface {
	GetServiceType() ServiceType
	GetOI4Identifier() Oi4Identifier
	GetApplicationSource() ApplicationSource
	GetLogger() *zap.SugaredLogger
	ResourceChanged(resource ResourceType, source BaseSource, filter *Filter)
	SendPublicationMessage(publication PublicationMessage)
	GetIntervalPublicationScheduler() IntervalPublicationScheduler

	PublicationProvider
}
