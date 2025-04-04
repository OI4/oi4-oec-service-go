package publication

import (
	"fmt"
	"github.com/OI4/oi4-oec-service-go/service/api"
	"github.com/OI4/oi4-oec-service-go/service/opc"
	"time"
)

func NewResourcePublication(application api.Oi4Application, oi4Source api.BaseSource, resourceType api.ResourceType) *Impl {
	return NewBuilder(application). //
					Oi4Source(oi4Source).                                      //
					Resource(resourceType).                                    //
					PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
					Build()
}

func NewResourcePublicationWithFilter(application api.Oi4Application, oi4Source api.BaseSource, resourceType api.ResourceType, filter *api.Filter) *Impl {
	return NewBuilder(application). //
					Oi4Source(oi4Source).                                      //
					Resource(resourceType).                                    //
					PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
					Filter(filter).                                            //
					Build()
}

func NewMAMPublication(application api.Oi4Application, oi4Source api.BaseSource) *Impl {
	mam := NewResourcePublication(application, oi4Source, api.ResourceMam)
	mam.doPublishOnRegistration = true
	return mam
}

func NewHealthPublication(application api.Oi4Application, oi4Source api.BaseSource) *IntervalPublicationImpl {
	return NewIntervalBuilder(application, 60*time.Second). //
								Oi4Source(oi4Source).                                      //
								Resource(api.ResourceHealth).                              //
								PublicationMode(api.PublicationMode_APPLICATION_SOURCE_5). //
								Build()
}

type BaseBuilder[T any] interface {
	Resource(resource api.ResourceType) T
	PublishOnRegistration(publishOnRegistration bool) T
	PublicationMode(mode api.PublicationMode) T
	PublicationConfig(config api.PublicationConfig) T
	StatusCode(status *api.StatusCode) T
	Filter(filter *api.Filter) T
	DataFunc(getDataFunc func() any) T
}

type Builder interface {
	BaseBuilder[Builder]
	Build() *Impl
}

// BuilderImpl we definitely need a mutex there :D
type BuilderImpl struct {
	application api.Oi4Application

	oi4Source               api.BaseSource
	resource                api.ResourceType
	doPublishOnRegistration bool

	filter            *api.Filter
	publicationMode   *api.PublicationMode
	publicationConfig api.PublicationConfig
	statusCode        *api.StatusCode
	getDataFunc       func() any
}

func NewBuilder(application api.Oi4Application) *BuilderImpl {
	mode := api.PublicationMode_ON_REQUEST_1
	builder := BuilderImpl{
		application:             application,
		publicationMode:         &mode,
		publicationConfig:       api.PublicationConfig_NONE_0,
		doPublishOnRegistration: false,
	}

	return &builder
}

func (p *BuilderImpl) Oi4Source(oi4Source api.BaseSource) Builder {
	p.oi4Source = oi4Source

	return p
}

func (p *BuilderImpl) Resource(resource api.ResourceType) Builder {
	p.resource = resource

	return p
}

func (p *BuilderImpl) PublishOnRegistration(publishOnRegistration bool) Builder {
	p.doPublishOnRegistration = publishOnRegistration

	return p
}

func (p *BuilderImpl) PublicationConfig(config api.PublicationConfig) Builder {
	p.publicationConfig = config

	return p
}

func (p *BuilderImpl) PublicationMode(publicationMode api.PublicationMode) Builder {
	p.publicationMode = &publicationMode

	return p
}

func (p *BuilderImpl) StatusCode(status *api.StatusCode) Builder {
	p.statusCode = status

	return p
}

func (p *BuilderImpl) Filter(filter *api.Filter) Builder {
	p.filter = filter

	return p
}

func (p *BuilderImpl) DataFunc(getDataFunc func() any) Builder {
	p.getDataFunc = getDataFunc

	return p
}

func (p *BuilderImpl) Build() *Impl {
	oi4Identifier := p.oi4Source.GetOi4Identifier()
	pub := Impl{
		application: p.application,
		resource:    p.resource,
		source:      oi4Identifier,
		filter:      p.filter,

		oi4Source:               p.oi4Source,
		doPublishOnRegistration: p.doPublishOnRegistration,
		dataSetWriterId:         opc.GetDataSetWriterId(p.resource, oi4Identifier),

		publicationMode:   p.publicationMode,
		publicationConfig: p.publicationConfig,
		statusCode:        p.statusCode,
		getDataFunc:       p.getDataFunc,
	}
	pub.id = fmt.Sprintf("%p", &pub)
	return &pub
}

type IntervalBuilder interface {
	BaseBuilder[IntervalBuilder]
	PublicationInterval(publicationInterval time.Duration) IntervalBuilder
	Build() *IntervalPublicationImpl
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

func (p *IntervalBuilderImpl) Oi4Source(oi4Source api.BaseSource) IntervalBuilder {
	p.oi4Source = oi4Source

	return p
}

func (p *IntervalBuilderImpl) Resource(resource api.ResourceType) IntervalBuilder {
	p.resource = resource

	return p
}

func (p *IntervalBuilderImpl) PublishOnRegistration(publishOnRegistration bool) IntervalBuilder {
	p.doPublishOnRegistration = publishOnRegistration

	return p
}

func (p *IntervalBuilderImpl) PublicationConfig(config api.PublicationConfig) IntervalBuilder {
	p.publicationConfig = config

	return p
}

func (p *IntervalBuilderImpl) PublicationMode(publicationMode api.PublicationMode) IntervalBuilder {
	p.publicationMode = &publicationMode

	return p
}

func (p *IntervalBuilderImpl) StatusCode(status *api.StatusCode) IntervalBuilder {
	p.statusCode = status

	return p
}

func (p *IntervalBuilderImpl) Filter(filter *api.Filter) IntervalBuilder {
	p.filter = filter

	return p
}

func (p *IntervalBuilderImpl) DataFunc(getDataFunc func() any) IntervalBuilder {
	p.getDataFunc = getDataFunc

	return p
}

func (p *IntervalBuilderImpl) PublicationInterval(publicationInterval time.Duration) IntervalBuilder {
	p.publicationInterval = publicationInterval

	return p
}

func (p *IntervalBuilderImpl) Build() *IntervalPublicationImpl {
	publication := *p.BuilderImpl.Build()

	return &IntervalPublicationImpl{
		Impl:                publication,
		publicationInterval: p.publicationInterval,
	}
}
