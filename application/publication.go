package application

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

type Publisher interface {
	SendPublicationMessage(PublicationMessage)
}

// we definitely need a mutex there :D
type Publication struct {
	publishOnRegistration bool
	parent                Publisher
	resource              v1.Resource
	publicationMode       v1.PublicationMode
	publicationConfig     v1.PublicationConfig
	statusCode            v1.StatusCode
	dataSetWriterId       uint16
	data                  interface{}
	publicationInterval   time.Duration
	getDataFunc           func() interface{}
	stopIntervalTicker    chan struct{}
}

func CreatePublication(resource v1.Resource, publishOnRegistration bool) *Publication {
	return &Publication{
		resource:              resource,
		publicationMode:       v1.PublicationMode_ON_REQUEST_1,
		publishOnRegistration: publishOnRegistration,
		statusCode:            0,
		publicationConfig:     v1.PublicationConfig_NONE_0,
		publicationInterval:   0,
		dataSetWriterId:       getNextDataSetWriterId(),
	}
}

func (p *Publication) SetPublicationMode(newMode v1.PublicationMode) *Publication {
	p.startPublicationTimer(p.publicationInterval, newMode)

	return p
}

func (p *Publication) SetPublicationConfig(newConfig v1.PublicationConfig) *Publication {
	p.publicationConfig = newConfig

	return p
}

func (p *Publication) SetPublicationInterval(newPublicationInterval time.Duration) *Publication {
	p.startPublicationTimer(newPublicationInterval, p.publicationMode)

	return p
}

func (p *Publication) SetStatusCode(status v1.StatusCode) *Publication {
	p.statusCode = status

	p.triggerPublication(false, false, "")
	return p
}

func (p *Publication) SetData(data interface{}) *Publication {
	p.data = data

	p.triggerPublication(false, false, "")
	return p
}

func (p *Publication) SetDataFunc(getDataFunc func() interface{}) *Publication {
	p.getDataFunc = getDataFunc

	return p
}

func (p *Publication) Publish() *Publication {
	p.triggerPublication(false, false, "")
	return p
}

func (p *Publication) stop() {
	if p.stopIntervalTicker != nil {
		p.stopIntervalTicker <- struct{}{}
		p.stopIntervalTicker = nil
	}
}

func (p *Publication) start() {
	p.startPublicationTimer(p.publicationInterval, p.publicationMode)
}

func (p *Publication) triggerPublication(byInterval bool, onRequest bool, correlationId string) {
	if p.parent != nil && onRequest ||
		(p.publicationMode != v1.PublicationMode_OFF_0 && p.publicationMode != v1.PublicationMode_ON_REQUEST_1 &&
			((p.publicationInterval == 0 && !byInterval) ||
				(p.publicationInterval != 0 && byInterval))) {
		message := PublicationMessage{
			resource:        p.resource,
			statusCode:      p.statusCode,
			publicationMode: p.publicationMode,
			dataSetWriterId: p.dataSetWriterId,
			correlationId:   correlationId,
		}
		if p.getDataFunc != nil {
			message.data = p.getDataFunc()
		} else {
			message.data = p.data
		}
		p.parent.SendPublicationMessage(message)
	}
}

func (p *Publication) startPublicationTimer(newPublicationInterval time.Duration, newPublicationMode v1.PublicationMode) {

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
