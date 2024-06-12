package service

import (
	"log/slog"
	"sync"
	"time"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/repository"
)

const (
	defaultEventCap = 10
	eventTypeCreate = 1
	eventTypeFollow = 2
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name EventsServiceConsumer
type EventsServiceConsumer interface {
	ConsumeEvents()
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name EventsServiceProducer
type EventsServiceProducer interface {
	ProduceEvent(event domain.URLEvent)
}

type eventsServiceConsumer struct {
	logger       *slog.Logger
	eventsCh     <-chan domain.URLEvent
	periodCh     <-chan time.Time
	doneCh       <-chan struct{}
	eventsWriter repository.EventsWriter
}

func NewEventsServiceConsumer(
	logger *slog.Logger,
	eventsCh <-chan domain.URLEvent,
	periodCh <-chan time.Time,
	doneCh <-chan struct{},
	clickhouseWriter repository.EventsWriter,
) EventsServiceConsumer {
	return &eventsServiceConsumer{
		logger:       logger,
		eventsCh:     eventsCh,
		periodCh:     periodCh,
		doneCh:       doneCh,
		eventsWriter: clickhouseWriter,
	}
}

func (s *eventsServiceConsumer) ConsumeEvents() {
	go func() {
		events := make([]domain.URLEvent, 0, defaultEventCap)
		m := &sync.Mutex{}

		for {
			select {
			case event, isOpen := <-s.eventsCh:
				if !isOpen {
					return
				}
				m.Lock()
				events = append(events, event)
				m.Unlock()
			case _, isOpen := <-s.periodCh:
				if !isOpen {
					return
				}
				if len(events) == 0 {
					continue
				}
				err := s.eventsWriter.Insert(events)
				if err != nil {
					s.logger.Error(err.Error())
				}
				events = make([]domain.URLEvent, 0, defaultEventCap)
			case <-s.doneCh:
				return
			}
		}
	}()
}

type eventsServiceProducer struct {
	eventsCh chan<- domain.URLEvent
}

func NewEventsServiceProducer(eventsCh chan<- domain.URLEvent) EventsServiceProducer {
	return &eventsServiceProducer{eventsCh: eventsCh}
}

func (s *eventsServiceProducer) ProduceEvent(event domain.URLEvent) {
	s.eventsCh <- event
}
