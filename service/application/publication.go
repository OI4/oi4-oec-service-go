package application

import (
	"time"

	"github.com/OI4/oi4-oec-service-go/service/api"
)

// Trigger type of trigger for publication
type Trigger int

// Declare typed constants each with type of status
const (
	ByInterval Trigger = iota
	OnRequest
)

type PublicationMessage struct {
	resource      api.ResourceType
	source        *api.Oi4Identifier
	correlationId string
	//publicationMode api.PublicationMode
	data       interface{}
	statusCode api.StatusCode
	filter     api.Filter
}

type PublicationPublisher interface {
	sendPublicationMessage(PublicationMessage)
}

type Publication interface {
	getResource() api.ResourceType
	getSource() *api.Oi4Identifier
	getPublicationMode() api.PublicationMode
	shouldPublicate() bool
	//triggerPublication(byInterval bool, onRequest bool, correlationId string)
	stop()
	start()
	setParent(PublicationPublisher)
	setSource(*api.Oi4Identifier)
	publishOnRegistration() bool
	getParent() PublicationPublisher
}

// PublicationImpl we definitely need a mutex there :D
type PublicationImpl struct {
	Publication
	oi4Source               *api.Oi4Source
	doPublishOnRegistration bool
	parent                  PublicationPublisher
	resource                api.ResourceType
	filter                  api.Filter //*string
	publicationMode         api.PublicationMode
	publicationConfig       api.PublicationConfig
	statusCode              api.StatusCode
	source                  *api.Oi4Identifier
	//data                    T
	publicationInterval time.Duration
	//getDataFunc             func() T
	stopIntervalTicker chan struct{}
}

func CreatePublication(resource api.ResourceType, oi4Source *api.Oi4Source) *PublicationImpl {
	publication := PublicationImpl{
		oi4Source:               oi4Source,
		resource:                resource,
		publicationMode:         api.PublicationMode_ON_REQUEST_1,
		doPublishOnRegistration: false,
		statusCode:              0,
		publicationConfig:       api.PublicationConfig_NONE_0,
		publicationInterval:     0,
	}

	return &publication
}

func (p *PublicationImpl) PublishOnRegistration(publishOnRegistration bool) *PublicationImpl {
	p.doPublishOnRegistration = publishOnRegistration

	return p
}

func (p *PublicationImpl) SetPublicationMode(newMode api.PublicationMode) *PublicationImpl {
	p.startPublicationTimer(p.publicationInterval, newMode)

	return p
}

func (p *PublicationImpl) SetPublicationConfig(newConfig api.PublicationConfig) *PublicationImpl {
	p.publicationConfig = newConfig

	return p
}

func (p *PublicationImpl) SetPublicationInterval(newPublicationInterval time.Duration) *PublicationImpl {
	p.startPublicationTimer(newPublicationInterval, p.publicationMode)

	return p
}

func (p *PublicationImpl) SetStatusCode(status api.StatusCode) *PublicationImpl {
	p.statusCode = status

	p.triggerPublication(false, false, "")
	return p
}

func (p *PublicationImpl) SetFilter(filter api.Filter) *PublicationImpl {
	p.filter = filter

	p.triggerPublication(false, false, "")
	return p
}

//func (p *PublicationImpl) SetData(data T) *PublicationImpl {
//	p.data = data
//
//	p.triggerPublication(false, false, "")
//	return p
//}
//
//func (p *PublicationImpl) SetDataFunc(getDataFunc func() T) *PublicationImpl {
//	p.getDataFunc = getDataFunc
//
//	return p
//}

func (p *PublicationImpl) Publish() *PublicationImpl {
	p.triggerPublication(false, false, "")
	return p
}

func (p *PublicationImpl) stop() {
	if p.stopIntervalTicker != nil {
		p.stopIntervalTicker <- struct{}{}
		p.stopIntervalTicker = nil
	}
}

func (p *PublicationImpl) setParent(parent PublicationPublisher) {
	p.parent = parent
}

func (p *PublicationImpl) setSource(source *api.Oi4Identifier) {
	p.source = source
}

func (p *PublicationImpl) getParent() PublicationPublisher {
	return p.parent
}

func (p *PublicationImpl) getResource() api.ResourceType {
	return p.resource
}

func (p *PublicationImpl) getSource() *api.Oi4Identifier {
	return p.source
}

func (p *PublicationImpl) getPublicationMode() api.PublicationMode {
	return p.publicationMode
}

func (p *PublicationImpl) publishOnRegistration() bool {
	return p.doPublishOnRegistration
}

func (p *PublicationImpl) start() {
	p.startPublicationTimer(p.publicationInterval, p.publicationMode)
}

func (p *PublicationImpl) shouldPublicate() bool {
	return true
}

func (p *PublicationImpl) triggerPublication(byInterval bool, onRequest bool, correlationId string) {
	if p.parent != nil && (onRequest ||
		(p.publicationMode != api.PublicationMode_OFF_0 && p.publicationMode != api.PublicationMode_ON_REQUEST_1 &&
			((p.publicationInterval == 0 && !byInterval) ||
				(p.publicationInterval != 0 && byInterval)))) {
		message := PublicationMessage{
			resource:   p.resource,
			statusCode: p.statusCode,
			//publicationMode: p.publicationMode,
			correlationId: correlationId,
			source:        p.source,
			filter:        p.filter,
		}
		//p.oi4Source.Get(p.resource)

		src := *p.oi4Source
		message.data = src.Get(p.resource)
		//if p.getDataFunc != nil {
		//	message.data = p.getDataFunc()
		//} else {
		//	message.data = p.data
		//}
		p.parent.sendPublicationMessage(message)
	}
}

func (p *PublicationImpl) startPublicationTimer(newPublicationInterval time.Duration, newPublicationMode api.PublicationMode) {

	// really, really shaky this function :D
	resetTimerInterval := false
	resetTimerMode := false
	if p.publicationInterval != newPublicationInterval && newPublicationInterval != 0 {
		resetTimerInterval = true
	}
	p.publicationInterval = newPublicationInterval

	if p.publicationMode != newPublicationMode && (newPublicationMode == api.PublicationMode_OFF_0 || newPublicationMode == api.PublicationMode_ON_REQUEST_1) {
		resetTimerMode = true
	}
	p.publicationMode = newPublicationMode

	// first clean off old timer because either publication type changed or publication interval
	if resetTimerInterval || resetTimerMode {
		if p.stopIntervalTicker != nil {
			p.stopIntervalTicker <- struct{}{}
			p.stopIntervalTicker = nil
		}
	}

	// start new timer if following conditions match
	if p.parent != nil && // we need a parent (the publication must be registered)
		p.stopIntervalTicker == nil && // we need no other active timer
		p.publicationInterval != 0 && // the publication interval needs to be != 0
		p.publicationMode != api.PublicationMode_OFF_0 && // the publication interval should not be turned off by mode
		p.publicationMode != api.PublicationMode_ON_REQUEST_1 {
		ticker := time.NewTicker(p.publicationInterval)
		stopChannel := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					p.triggerPublication(true, false, "")
				case <-stopChannel:
					ticker.Stop()
					return
				}
			}
		}()
		p.stopIntervalTicker = stopChannel
	}
}
