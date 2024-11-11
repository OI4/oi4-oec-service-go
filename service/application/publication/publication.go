package publication

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
)

// Impl PublicationImpl we definitely need a mutex there :D
type Impl struct {
	api.Publication
	id                      string
	application             api.Oi4Application
	oi4Source               api.BaseSource
	doPublishOnRegistration bool
	resource                api.ResourceType
	filter                  api.Filter
	publicationMode         *api.PublicationMode
	publicationConfig       api.PublicationConfig
	statusCode              *api.StatusCode
	source                  *api.Oi4Identifier
	dataSetWriterId         uint16
	//Data                    T
	getDataFunc        func() any
	stopIntervalTicker chan struct{}
}

func (p *Impl) GetPublicationType() api.PublicationType {
	return api.General
}

func (p *Impl) Stop() {
	if p.stopIntervalTicker != nil {
		p.stopIntervalTicker <- struct{}{}
		p.stopIntervalTicker = nil
	}
}

func (p *Impl) GetApplication() api.Oi4Application {
	return p.application
}

func (p *Impl) GetOi4Source() api.BaseSource {
	return p.oi4Source
}

func (p *Impl) GetResource() api.ResourceType {
	return p.resource
}

func (p *Impl) GetSource() *api.Oi4Identifier {
	return p.source
}

func (p *Impl) GetDataSetWriterId() uint16 {
	return p.dataSetWriterId
}

func (p *Impl) GetFilter() api.Filter {
	return p.filter
}

func (p *Impl) GetID() string {
	return p.id
}

func (p *Impl) GetPublicationMode() *api.PublicationMode {
	return p.publicationMode
}

func (p *Impl) publishOnRegistration() bool {
	return p.doPublishOnRegistration
}

func (p *Impl) Start() {
	if p.doPublishOnRegistration {
		p.triggerPublication(nil)
	}
}

func (p *Impl) ShouldPublicate(trigger api.Trigger) bool {
	if trigger == api.OnRequest {
		return true
	}

	mode := getPublicationMode(p.publicationMode)
	if mode == api.PublicationMode_OFF_0 || //
		mode == api.PublicationMode_ON_REQUEST_1 {
		return false
	}

	return false
}

func (p *Impl) TriggerPublication(trigger api.Trigger, correlationId *string) bool {
	if !p.ShouldPublicate(trigger) {
		return false
	}

	p.triggerPublication(correlationId)
	return true
}

func (p *Impl) triggerPublication(correlationId *string) {
	if p.application == nil {
		return
	}

	resource := p.GetResource()
	source := p.GetOi4Source()

	data := source.Get(resource)
	content := make([]api.PublicationContent, len(data))
	for i, single := range data {
		if single == nil {
			continue
		}
		content[i] = api.PublicationContent{
			StatusCode: p.statusCode,
			Data:       single,
		}
	}

	message := api.PublicationMessage{
		Filter:   p.filter,
		Resource: resource,
		//StatusCode: api.Status_Good,
		//publicationMode: p.publicationMode,
		CorrelationId: correlationId,
		Source:        source.GetOi4Identifier(),
		//Filter:        p.GetFilter(),
		Content: content,
	}
	//message.Data = source.Get(resource)

	p.application.SendPublicationMessage(message)
}

func getPublicationMode(mode *api.PublicationMode) api.PublicationMode {
	if mode == nil {
		return api.PublicationMode_OFF_0
	}
	return *mode
}
