package publication

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	//"github.com/OI4/oi4-oec-service-go/service/application"
)

// Trigger type of trigger for publication
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
	GetPublicationType() PublicationType
	ShouldPublicate(trigger Trigger) bool
	GetOi4Source() api.Oi4Source
	GetResource() api.ResourceType
	GetSource() *api.Oi4Identifier
	GetFilter() api.Filter

	Stop()
	Start()

	TriggerPublication(trigger Trigger, correlationId string)

	getApplication() api.Oi4Application
	getPublicationMode() api.PublicationMode
	//triggerSourcePublication(byInterval bool, onRequest bool, correlationId string)
	publishOnRegistration() bool
}

// PublicationImpl we definitely need a mutex there :D
type PublicationImpl struct {
	Publication
	application             api.Oi4Application
	oi4Source               api.Oi4Source
	doPublishOnRegistration bool
	resource                api.ResourceType
	filter                  api.Filter //*string
	publicationMode         api.PublicationMode
	publicationConfig       api.PublicationConfig
	statusCode              api.StatusCode
	source                  *api.Oi4Identifier
	//Data                    T
	getDataFunc        func() any
	stopIntervalTicker chan struct{}
}

func (p *PublicationImpl) GetPublicationType() PublicationType {
	return General
}

func (p *PublicationImpl) Stop() {
	if p.stopIntervalTicker != nil {
		p.stopIntervalTicker <- struct{}{}
		p.stopIntervalTicker = nil
	}
}

func (p *PublicationImpl) getApplication() api.Oi4Application {
	return p.application
}

func (p *PublicationImpl) GetOi4Source() api.Oi4Source {
	return p.oi4Source
}

func (p *PublicationImpl) GetResource() api.ResourceType {
	return p.resource
}

func (p *PublicationImpl) GetSource() *api.Oi4Identifier {
	return p.source
}

func (p *PublicationImpl) GetFilter() api.Filter {
	return p.filter
}

func (p *PublicationImpl) getPublicationMode() api.PublicationMode {
	return p.publicationMode
}

func (p *PublicationImpl) publishOnRegistration() bool {
	return p.doPublishOnRegistration
}

func (p *PublicationImpl) Start() {
	if p.doPublishOnRegistration {
		p.triggerPublication(OnRequest, "")
	}
}

func (p *PublicationImpl) ShouldPublicate(trigger Trigger) bool {
	if trigger == OnRequest {
		return true
	}

	mode := p.publicationMode
	if mode == api.PublicationMode_OFF_0 || //
		mode == api.PublicationMode_ON_REQUEST_1 {
		return false
	}

	//interval := p.publicationInterval
	//if interval == 0 && trigger != ByInterval || //
	//	interval != 0 && trigger == ByInterval {
	//	return true
	//}

	return false
}

func (p *PublicationImpl) triggerPublication(trigger Trigger, correlationId string) { //(byInterval bool, onRequest bool, correlationId string) {
	if p.application == nil {
		return
	}
	p.TriggerPublication(trigger, correlationId)
}

func (p *PublicationImpl) TriggerPublication(trigger Trigger, correlationId string) {

	if !p.ShouldPublicate(trigger) {
		return
	}

	resource := p.GetResource()
	source := p.GetOi4Source()

	message := api.PublicationMessage{
		Resource:   resource,
		StatusCode: api.Status_Good,
		//publicationMode: p.publicationMode,
		CorrelationId: correlationId,
		Source:        source.GetOi4Identifier(),
		Filter:        p.GetFilter(),
	}
	message.Data = source.Get(resource)

	p.application.SendPublicationMessage(message)
}
