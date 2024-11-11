package publication

import (
	"github.com/OI4/oi4-oec-service-go/service/api"
	"go.uber.org/zap"
	"sync"
	"time"
)

type IntervalPublicationImpl struct {
	Impl

	publicationInterval time.Duration
	lastPublication     time.Time
}

func (p *IntervalPublicationImpl) GetPublicationType() api.PublicationType {
	return api.Interval
}

func (p *IntervalPublicationImpl) GetNextPublicationTime() time.Time {
	return p.lastPublication.Add(p.publicationInterval)
}

func (p *IntervalPublicationImpl) DueForPublication() bool {
	return time.Now().After(p.GetNextPublicationTime())
}

func (p *IntervalPublicationImpl) ShouldPublicate(trigger api.Trigger) bool {
	if trigger == api.OnRequest {
		return true
	}

	mode := getPublicationMode(p.publicationMode)
	if mode == api.PublicationMode_OFF_0 || //
		mode == api.PublicationMode_ON_REQUEST_1 {
		return false
	}

	interval := p.publicationInterval
	if interval == 0 && trigger != api.ByInterval || //
		interval != 0 && trigger == api.ByInterval {
		return true
	}

	return false
}

func (p *IntervalPublicationImpl) TriggerPublication(trigger api.Trigger, correlationId *string) bool {
	if !p.ShouldPublicate(trigger) {
		return false
	}

	p.triggerPublication(correlationId)

	if trigger == api.ByInterval {
		p.lastPublication = time.Now()
	}

	return true
}

func (p *IntervalPublicationImpl) Start() {
	p.application.GetIntervalPublicationScheduler().AddPublication(p)
}

func (p *IntervalPublicationImpl) Stop() {
	p.application.GetIntervalPublicationScheduler().RemovePublication(p)
}

type IntervalPublicationSchedulerImpl struct {
	publications         map[string]api.IntervalPublication
	workQueue            chan api.IntervalPublication
	doneQueue            chan string
	enqueuedPublications map[string]bool
	queueSize            int
	workerCount          int
	mu                   sync.RWMutex
	ticker               *time.Ticker
}

func NewIntervalPublicationSchedulerImpl(queueSize int, workerCount int) *IntervalPublicationSchedulerImpl {
	return &IntervalPublicationSchedulerImpl{
		publications:         make(map[string]api.IntervalPublication),
		enqueuedPublications: make(map[string]bool),
		queueSize:            queueSize,
		workerCount:          workerCount,
		mu:                   sync.RWMutex{},
	}
}

func (s *IntervalPublicationSchedulerImpl) Start() {
	s.workQueue = make(chan api.IntervalPublication, s.queueSize)
	s.doneQueue = make(chan string, s.workerCount)

	for i := 1; i <= s.workerCount; i++ {
		go Worker(i, s.workQueue, s.doneQueue)
	}

	s.ticker = time.NewTicker(50 * time.Millisecond)

	go func() {
		for range s.ticker.C {
			for _, pub := range s.publications {
				s.mu.Lock()
				if pub.DueForPublication() && !s.enqueuedPublications[pub.GetID()] {
					s.enqueuedPublications[pub.GetID()] = true
					s.workQueue <- pub
				}
				s.mu.Unlock()
			}

			// Process completed items
			select {
			case id := <-s.doneQueue:
				s.mu.Lock()
				delete(s.enqueuedPublications, id)
				s.mu.Unlock()
			default:
				// No completed items to process
			}
		}
	}()

}

func (s *IntervalPublicationSchedulerImpl) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}

	if s.workQueue != nil {
		close(s.workQueue)
	}

	if s.doneQueue != nil {
		close(s.doneQueue)
	}
}

func (s *IntervalPublicationSchedulerImpl) AddPublication(publication api.IntervalPublication) {
	s.publications[publication.GetID()] = publication
}

func (s *IntervalPublicationSchedulerImpl) RemovePublication(publication api.IntervalPublication) {
	delete(s.publications, publication.GetID())
}

func getLogger(publication api.IntervalPublication) *zap.SugaredLogger {
	return publication.GetApplication().GetLogger()
}

func Worker(id int, workQueue <-chan api.IntervalPublication, doneQueue chan<- string) {
	for pub := range workQueue {
		getLogger(pub).Debugf("Worker %d processing publication %s - %s", id, pub.GetSource().ToString(), pub.GetResource())
		pub.TriggerPublication(api.ByInterval, nil)

		doneQueue <- pub.GetID()
	}
}
