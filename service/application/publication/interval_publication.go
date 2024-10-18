package publication

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"time"
)

type IntervalPublication interface {
	Publication
	GetNextPublicationTime() time.Time
	DueForPublication() bool
}

type IntervalPublicationImpl struct {
	PublicationImpl

	publicationInterval time.Duration
	lastPublication     time.Time
}

func (p *IntervalPublicationImpl) GetPublicationType() PublicationType {
	return Interval
}

func (p *IntervalPublicationImpl) GetNextPublicationTime() time.Time {
	return p.lastPublication.Add(p.publicationInterval)
}

func (p *IntervalPublicationImpl) DueForPublication() bool {
	return time.Now().After(p.GetNextPublicationTime())
}

func (p *IntervalPublicationImpl) ShouldPublicate(trigger Trigger) bool {
	if trigger == OnRequest {
		return true
	}

	mode := p.publicationMode
	if mode == api.PublicationMode_OFF_0 || //
		mode == api.PublicationMode_ON_REQUEST_1 {
		return false
	}

	interval := p.publicationInterval
	if interval == 0 && trigger != ByInterval || //
		interval != 0 && trigger == ByInterval {
		return true
	}

	return false
}

func (p *IntervalPublicationImpl) startPublicationTimer(newPublicationInterval time.Duration, newPublicationMode api.PublicationMode) {

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

	// Start new timer if following conditions match
	if p.application != nil && // we need a parent (the publication must be registered)
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
					p.triggerPublication(ByInterval, "")
				case <-stopChannel:
					ticker.Stop()
					return
				}
			}
		}()
		p.stopIntervalTicker = stopChannel
	}
}

func (p *IntervalPublicationImpl) Start() {
	p.startPublicationTimer(p.publicationInterval, p.publicationMode)
	p.PublicationImpl.Start()
}
