package publication

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"time"
)

type BaseBuilder interface {
	Source(source *api.Oi4Identifier) *BaseBuilder
	Resource(resource api.ResourceType) *BaseBuilder
	PublishOnRegistration(publishOnRegistration bool) *BaseBuilder
	PublicationMode(mode api.PublicationMode) *BaseBuilder
	PublicationConfig(config api.PublicationConfig) *BaseBuilder
	StatusCode(status api.StatusCode) *BaseBuilder
	Filter(filter api.Filter) *BaseBuilder
	DataFunc(getDataFunc func()) *BaseBuilder
}

type Builder interface {
	BaseBuilder
	Build() Publication
}

// BuilderImpl we definitely need a mutex there :D
type BuilderImpl struct {
	application api.Oi4Application

	oi4Source               api.Oi4Source
	resource                api.ResourceType
	doPublishOnRegistration bool

	filter            api.Filter //*string
	publicationMode   api.PublicationMode
	publicationConfig api.PublicationConfig
	statusCode        api.StatusCode
	getDataFunc       func() any
}

func NewBuilder(application api.Oi4Application) *BuilderImpl {
	builder := BuilderImpl{
		application:             application,
		publicationMode:         api.PublicationMode_ON_REQUEST_1,
		statusCode:              0,
		publicationConfig:       api.PublicationConfig_NONE_0,
		doPublishOnRegistration: false,
	}

	return &builder
}

func (p *BuilderImpl) Oi4Source(oi4Source api.Oi4Source) *BuilderImpl {
	p.oi4Source = oi4Source

	return p
}

func (p *BuilderImpl) Resource(resource api.ResourceType) *BuilderImpl {
	p.resource = resource

	return p
}

func (p *BuilderImpl) PublishOnRegistration(publishOnRegistration bool) *BuilderImpl {
	p.doPublishOnRegistration = publishOnRegistration

	return p
}

func (p *BuilderImpl) PublicationConfig(config api.PublicationConfig) *BuilderImpl {
	p.publicationConfig = config

	return p
}

func (p *BuilderImpl) PublicationMode(publicationMode api.PublicationMode) *BuilderImpl {
	p.publicationMode = publicationMode

	return p
}

func (p *BuilderImpl) StatusCode(status api.StatusCode) *BuilderImpl {
	p.statusCode = status

	return p
}

func (p *BuilderImpl) Filter(filter api.Filter) *BuilderImpl {
	p.filter = filter

	return p
}

func (p *BuilderImpl) DataFunc(getDataFunc func() any) *BuilderImpl {
	p.getDataFunc = getDataFunc

	return p
}

func (p *BuilderImpl) Build() *PublicationImpl {
	return &PublicationImpl{
		application: p.application,
		resource:    p.resource,
		source:      p.oi4Source.GetOi4Identifier(),
		filter:      p.filter,

		oi4Source:               p.oi4Source,
		doPublishOnRegistration: p.doPublishOnRegistration,

		publicationMode:   p.publicationMode,
		publicationConfig: p.publicationConfig,
		statusCode:        p.statusCode,
		getDataFunc:       p.getDataFunc,
	}
}

type IntervalBuilder interface {
	BaseBuilder
	PublicationInterval(publicationInterval time.Duration) *IntervalPublication
	Build() IntervalPublication
}

type IntervalBuilderImpl struct {
	BuilderImpl
	publicationInterval time.Duration
}

func NewIntervalBuilder(application api.Oi4Application, publicationInterval time.Duration) *IntervalBuilderImpl {
	builder := *NewBuilder(application)
	return &IntervalBuilderImpl{
		builder,
		publicationInterval,
	}
}

func (p *IntervalBuilderImpl) Oi4Source(oi4Source api.Oi4Source) *IntervalBuilderImpl {
	p.oi4Source = oi4Source

	return p
}

func (p *IntervalBuilderImpl) PublicationInterval(publicationInterval time.Duration) *IntervalBuilderImpl {
	p.publicationInterval = publicationInterval

	return p
}

func (p *IntervalBuilderImpl) Build() *IntervalPublicationImpl {
	publication := *p.BuilderImpl.Build()

	return &IntervalPublicationImpl{
		PublicationImpl:     publication,
		publicationInterval: p.publicationInterval,
	}
}
