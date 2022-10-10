package service

import (
	"time"

	v1 "github.com/mzeiher/oi4/api/pkg/types"
)

var dataSetWriterId uint16 = 10

func getNextDataSetWriterId() uint16 {
	nextId := dataSetWriterId
	dataSetWriterId++
	// this could overflow....
	return nextId
}

type PublicationMessage struct {
	resource        v1.Resource
	source          *v1.Oi4Identifier
	correlationId   string
	publicationMode v1.PublicationMode
	data            interface{}
	statusCode      v1.StatusCode
	dataSetWriterId uint16
}

type PublicationPublisher interface {
	sendPublicationMessage(PublicationMessage)
}

type Publication interface {
	getResource() v1.Resource
	triggerPublication(byInterval bool, onRequest bool, correlationId string)
	stop()
	start()
	setParent(PublicationPublisher)
	setSource(*v1.Oi4Identifier)
	publishOnRegistration() bool
	getParent() PublicationPublisher
}

// we definitely need a mutex there :D
type PublicationImpl[T interface{}] struct {
	Publication
	doPublishOnRegistration bool
	parent                  PublicationPublisher
	resource                v1.Resource
	publicationMode         v1.PublicationMode
	publicationConfig       v1.PublicationConfig
	statusCode              v1.StatusCode
	source                  *v1.Oi4Identifier
	dataSetWriterId         uint16
	data                    T
	publicationInterval     time.Duration
	getDataFunc             func() T
	stopIntervalTicker      chan struct{}
}

func CreatePublication[T interface{}](resource v1.Resource, publishOnRegistration bool) *PublicationImpl[T] {
	return &PublicationImpl[T]{
		resource:                resource,
		publicationMode:         v1.PublicationMode_ON_REQUEST_1,
		doPublishOnRegistration: publishOnRegistration,
		statusCode:              0,
		publicationConfig:       v1.PublicationConfig_NONE_0,
		publicationInterval:     0,
		dataSetWriterId:         getNextDataSetWriterId(),
	}
}

func (p *PublicationImpl[T]) SetPublicationMode(newMode v1.PublicationMode) *PublicationImpl[T] {
	p.startPublicationTimer(p.publicationInterval, newMode)

	return p
}

func (p *PublicationImpl[T]) SetPublicationConfig(newConfig v1.PublicationConfig) *PublicationImpl[T] {
	p.publicationConfig = newConfig

	return p
}

func (p *PublicationImpl[T]) SetPublicationInterval(newPublicationInterval time.Duration) *PublicationImpl[T] {
	p.startPublicationTimer(newPublicationInterval, p.publicationMode)

	return p
}

func (p *PublicationImpl[T]) SetStatusCode(status v1.StatusCode) *PublicationImpl[T] {
	p.statusCode = status

	p.triggerPublication(false, false, "")
	return p
}

func (p *PublicationImpl[T]) SetData(data T) *PublicationImpl[T] {
	p.data = data

	p.triggerPublication(false, false, "")
	return p
}

func (p *PublicationImpl[T]) SetDataFunc(getDataFunc func() T) *PublicationImpl[T] {
	p.getDataFunc = getDataFunc

	return p
}

func (p *PublicationImpl[T]) Publish() *PublicationImpl[T] {
	p.triggerPublication(false, false, "")
	return p
}

func (p *PublicationImpl[T]) stop() {
	if p.stopIntervalTicker != nil {
		p.stopIntervalTicker <- struct{}{}
		p.stopIntervalTicker = nil
	}
}

func (p *PublicationImpl[T]) setParent(parent PublicationPublisher) {
	p.parent = parent
}

func (p *PublicationImpl[T]) setSource(source *v1.Oi4Identifier) {
	p.source = source
}

func (p *PublicationImpl[T]) getParent() PublicationPublisher {
	return p.parent
}

func (p *PublicationImpl[T]) getResource() v1.Resource {
	return p.resource
}

func (p *PublicationImpl[T]) publishOnRegistration() bool {
	return p.doPublishOnRegistration
}

func (p *PublicationImpl[T]) start() {
	p.startPublicationTimer(p.publicationInterval, p.publicationMode)
}

func (p *PublicationImpl[T]) triggerPublication(byInterval bool, onRequest bool, correlationId string) {
	if p.parent != nil && (onRequest ||
		(p.publicationMode != v1.PublicationMode_OFF_0 && p.publicationMode != v1.PublicationMode_ON_REQUEST_1 &&
			((p.publicationInterval == 0 && !byInterval) ||
				(p.publicationInterval != 0 && byInterval)))) {
		message := PublicationMessage{
			resource:        p.resource,
			statusCode:      p.statusCode,
			publicationMode: p.publicationMode,
			dataSetWriterId: p.dataSetWriterId,
			correlationId:   correlationId,
			source:          p.source,
		}
		if p.getDataFunc != nil {
			message.data = p.getDataFunc()
		} else {
			message.data = p.data
		}
		p.parent.sendPublicationMessage(message)
	}
}

func (p *PublicationImpl[T]) startPublicationTimer(newPublicationInterval time.Duration, newPublicationMode v1.PublicationMode) {

	// really really shaky this function :D
	resetTimerInterval := false
	resetTimerMode := false
	if p.publicationInterval != newPublicationInterval && newPublicationInterval != 0 {
		resetTimerInterval = true
	}
	p.publicationInterval = newPublicationInterval

	if p.publicationMode != newPublicationMode && (newPublicationMode == v1.PublicationMode_OFF_0 || newPublicationMode == v1.PublicationMode_ON_REQUEST_1) {
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
		p.publicationMode != v1.PublicationMode_OFF_0 && // the publication interval should not be turned off by mode
		p.publicationMode != v1.PublicationMode_ON_REQUEST_1 {
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
