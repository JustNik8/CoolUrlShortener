package service

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/repository"
	"CoolUrlShortener/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConsumeEvents(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	currentTime := time.Now()

	testCases := []struct {
		name              string
		buildEventsWriter func(actualBatches map[int][]domain.URLEvent) repository.EventsWriter
		sendEvents        func(eventsCh chan<- domain.URLEvent, events []domain.URLEvent)
		batchPeriod       time.Duration
		waitTime          time.Duration
		buildEvents       func() []domain.URLEvent
		expectedBatches   map[int][]domain.URLEvent
	}{
		{
			name: "send 3 events in 1 batch",
			buildEventsWriter: func(actualBatches map[int][]domain.URLEvent) repository.EventsWriter {
				mockWriter := mocks.NewEventsWriter(t)
				i := 0
				mockWriter.On("Insert", mock.AnythingOfType("[]domain.URLEvent")).
					Run(func(args mock.Arguments) {
						actualBatches[i] = args.Get(0).([]domain.URLEvent)
						i++
					}).
					Return(nil).
					Once()

				return mockWriter
			},
			sendEvents: func(eventsCh chan<- domain.URLEvent, events []domain.URLEvent) {
				for i := range events {
					eventsCh <- events[i]
				}
			},
			batchPeriod: 200 * time.Millisecond,
			waitTime:    500 * time.Millisecond,
			buildEvents: func() []domain.URLEvent {
				return []domain.URLEvent{
					{LongURL: "https://test.long", ShortURL: "short", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long2", ShortURL: "short2", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long3", ShortURL: "short3", EventTime: currentTime, EventType: eventTypeFollow},
				}
			},
			expectedBatches: map[int][]domain.URLEvent{
				0: {
					{LongURL: "https://test.long", ShortURL: "short", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long2", ShortURL: "short2", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long3", ShortURL: "short3", EventTime: currentTime, EventType: eventTypeFollow},
				},
			},
		},
		{
			name: "send 6 events in 2 batches",
			buildEventsWriter: func(actualBatches map[int][]domain.URLEvent) repository.EventsWriter {
				mockWriter := mocks.NewEventsWriter(t)
				i := 0
				mockWriter.On("Insert", mock.AnythingOfType("[]domain.URLEvent")).
					Run(func(args mock.Arguments) {
						actualBatches[i] = args.Get(0).([]domain.URLEvent)
						i++
					}).
					Return(nil).
					Times(2)

				return mockWriter
			},
			sendEvents: func(eventsCh chan<- domain.URLEvent, events []domain.URLEvent) {
				size := len(events)
				for i := 0; i < size/2; i++ {
					eventsCh <- events[i]
				}
				time.Sleep(250 * time.Millisecond)

				for i := size / 2; i < size; i++ {
					eventsCh <- events[i]
				}
			},
			batchPeriod: 200 * time.Millisecond,
			waitTime:    1 * time.Second,
			buildEvents: func() []domain.URLEvent {
				return []domain.URLEvent{
					{LongURL: "https://test.long", ShortURL: "short", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long2", ShortURL: "short2", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long3", ShortURL: "short3", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long4", ShortURL: "short4", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long5", ShortURL: "short5", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long6", ShortURL: "short6", EventTime: currentTime, EventType: eventTypeFollow},
				}
			},
			expectedBatches: map[int][]domain.URLEvent{
				0: {
					{LongURL: "https://test.long", ShortURL: "short", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long2", ShortURL: "short2", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long3", ShortURL: "short3", EventTime: currentTime, EventType: eventTypeFollow},
				},
				1: {
					{LongURL: "https://test.long4", ShortURL: "short4", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long5", ShortURL: "short5", EventTime: currentTime, EventType: eventTypeFollow},
					{LongURL: "https://test.long6", ShortURL: "short6", EventTime: currentTime, EventType: eventTypeFollow},
				},
			},
		},
		{
			name: "do not send batch. 0 events",
			buildEventsWriter: func(actualBatches map[int][]domain.URLEvent) repository.EventsWriter {
				mockWriter := mocks.NewEventsWriter(t)
				return mockWriter
			},
			sendEvents: func(eventsCh chan<- domain.URLEvent, events []domain.URLEvent) {
				for i := range events {
					eventsCh <- events[i]
				}
			},
			batchPeriod: 200 * time.Millisecond,
			waitTime:    500 * time.Millisecond,
			buildEvents: func() []domain.URLEvent {
				return []domain.URLEvent{}
			},
			expectedBatches: map[int][]domain.URLEvent{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eventsCh := make(chan domain.URLEvent)
			events := tc.buildEvents()

			actualBatches := make(map[int][]domain.URLEvent)
			eventsWriter := tc.buildEventsWriter(actualBatches)

			doneCh := make(chan struct{})

			periodCh := time.NewTicker(tc.batchPeriod).C
			eventsConsumer := NewEventsServiceConsumer(
				logger,
				eventsCh,
				periodCh,
				doneCh,
				eventsWriter,
			)

			eventsConsumer.ConsumeEvents()
			tc.sendEvents(eventsCh, events)

			time.Sleep(tc.waitTime)
			doneCh <- struct{}{}
			close(eventsCh)

			assert.Equal(t, len(tc.expectedBatches), len(actualBatches))
			for k, expectedBatch := range tc.expectedBatches {
				actualBatch, exists := actualBatches[k]
				assert.True(t, exists)

				assert.Equal(t, expectedBatch, actualBatch)
			}
		})
	}

}
