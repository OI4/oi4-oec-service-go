package api

import "go.uber.org/zap"

type Oi4Application interface {
	GetLogger() *zap.SugaredLogger
	ResourceChanged(resource ResourceType, source BaseSource, filter Filter)
	SendPublicationMessage(publication PublicationMessage)
	GetApplicationSource() ApplicationSource
	GetIntervalPublicationScheduler() IntervalPublicationScheduler

	PublicationProvider
}
