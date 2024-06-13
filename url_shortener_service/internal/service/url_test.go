package service

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/repository"
	"CoolUrlShortener/internal/repository/mocks"
	"CoolUrlShortener/pkg/shortener"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	shortenermocks "CoolUrlShortener/pkg/shortener/mocks"
)

func TestGetLongURL(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	urlShortener := shortenermocks.NewURLShortener(t)

	testLongURL := "https://test.longurl"
	testShortURL := "short"

	testCases := []struct {
		name                string
		buildURLRepo        func() repository.UrlRepo
		buildURLCache       func() repository.URLCache
		buildEventsProducer func() repository.EventsProducer
		expectedLongURL     string
		expectedErr         error
	}{
		{
			name: "Get long url from cache",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)
				mockCache.On("GetLongURL", mock.Anything, mock.Anything).
					Return(testLongURL, nil).
					Once()

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)
				mockEventsServiceProducer.On("ProduceEvent", mock.Anything).
					Once()

				return mockEventsServiceProducer
			},
			expectedLongURL: testLongURL,
			expectedErr:     nil,
		},
		{
			name: "Get long url from database",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetLongURL", mock.Anything, testShortURL).
					Return(testLongURL, nil).
					Once()

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)
				mockCache.On("GetLongURL", mock.Anything, mock.Anything).
					Return("", errors.New("no long url in cache")).
					Once()

				mockCache.On("SetLongURL", mock.Anything, testShortURL, testLongURL).
					Return(nil).
					Once()

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)
				mockEventsServiceProducer.On("ProduceEvent", mock.Anything).
					Once()

				return mockEventsServiceProducer
			},
			expectedLongURL: testLongURL,
			expectedErr:     nil,
		},
		{
			name: "long url not found in db. Should be error",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetLongURL", mock.Anything, testShortURL).
					Return("", errs.ErrNoURL).
					Once()

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)
				mockCache.On("GetLongURL", mock.Anything, mock.Anything).
					Return("", errors.New("no long url in cache")).
					Once()

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)

				return mockEventsServiceProducer
			},
			expectedLongURL: "",
			expectedErr:     errs.ErrNoURL,
		},
		{
			name: "could not write to cache. Should not be error",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetLongURL", mock.Anything, testShortURL).
					Return(testLongURL, nil).
					Once()

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)
				mockCache.On("GetLongURL", mock.Anything, mock.Anything).
					Return("", errors.New("no long url in cache")).
					Once()

				mockCache.On("SetLongURL", mock.Anything, testShortURL, testLongURL).
					Return(errors.New("unexpected error"))

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)
				mockEventsServiceProducer.On("ProduceEvent", mock.Anything).
					Once()

				return mockEventsServiceProducer
			},
			expectedLongURL: testLongURL,
			expectedErr:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			urlService := NewURLService(
				logger,
				tc.buildURLRepo(),
				tc.buildURLCache(),
				tc.buildEventsProducer(),
				urlShortener,
			)

			longURL, err := urlService.GetLongURL(context.Background(), testShortURL)
			assert.Equal(t, tc.expectedLongURL, longURL)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestSaveURL(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	testLongURL := "https://test.longurl"
	testShortURL := "short"

	unexpectedErr := errors.New("unexpected error")

	testCases := []struct {
		name                string
		buildURLRepo        func() repository.UrlRepo
		buildURLCache       func() repository.URLCache
		buildEventsProducer func() repository.EventsProducer
		buildURLShortener   func() shortener.URLShortener
		expectedShortURL    string
		expectedErr         error
	}{
		{
			name: "Short url exists. Should return existing short url",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetShortURLByLongURL", mock.Anything, testLongURL).
					Return(testShortURL, nil)

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)
				mockEventsServiceProducer.On("ProduceEvent", mock.Anything).
					Once()

				return mockEventsServiceProducer
			},
			buildURLShortener: func() shortener.URLShortener {
				mockURLShortener := shortenermocks.NewURLShortener(t)
				return mockURLShortener
			},
			expectedShortURL: testShortURL,
			expectedErr:      nil,
		},
		{
			name: "unexpected error when reading db",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetShortURLByLongURL", mock.Anything, testLongURL).
					Return("", unexpectedErr)

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)

				return mockEventsServiceProducer
			},
			buildURLShortener: func() shortener.URLShortener {
				mockURLShortener := shortenermocks.NewURLShortener(t)
				return mockURLShortener
			},
			expectedShortURL: "",
			expectedErr:      unexpectedErr,
		},
		{
			name: "create new short url without error",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetShortURLByLongURL", mock.Anything, testLongURL).
					Return("", errs.ErrNoURL)

				mockRepo.On("SaveURL", mock.Anything, mock.Anything).
					Return(nil)

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)
				mockCache.On("SetLongURL", mock.Anything, mock.Anything, mock.Anything).
					Return(nil)

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)
				mockEventsServiceProducer.On("ProduceEvent", mock.Anything).
					Once()

				return mockEventsServiceProducer
			},
			buildURLShortener: func() shortener.URLShortener {
				mockURLShortener := shortenermocks.NewURLShortener(t)
				mockURLShortener.On("ShortenURL", mock.AnythingOfType("uint32")).
					Return(testShortURL)

				return mockURLShortener
			},
			expectedShortURL: testShortURL,
			expectedErr:      nil,
		},
		{
			name: "error while saving url to db. Should return error",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetShortURLByLongURL", mock.Anything, testLongURL).
					Return("", errs.ErrNoURL)

				mockRepo.On("SaveURL", mock.Anything, mock.Anything).
					Return(unexpectedErr)

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)

				return mockEventsServiceProducer
			},
			buildURLShortener: func() shortener.URLShortener {
				mockURLShortener := shortenermocks.NewURLShortener(t)
				mockURLShortener.On("ShortenURL", mock.AnythingOfType("uint32")).
					Return(testShortURL)

				return mockURLShortener
			},
			expectedShortURL: "",
			expectedErr:      unexpectedErr,
		},
		{
			name: "error while saving url to cache. Should not return error",
			buildURLRepo: func() repository.UrlRepo {
				mockRepo := mocks.NewUrlRepo(t)
				mockRepo.On("GetShortURLByLongURL", mock.Anything, testLongURL).
					Return("", errs.ErrNoURL)

				mockRepo.On("SaveURL", mock.Anything, mock.Anything).
					Return(nil)

				return mockRepo
			},
			buildURLCache: func() repository.URLCache {
				mockCache := mocks.NewURLCache(t)
				mockCache.On("SetLongURL", mock.Anything, testShortURL, testLongURL).
					Return(unexpectedErr)

				return mockCache
			},
			buildEventsProducer: func() repository.EventsProducer {
				mockEventsServiceProducer := mocks.NewEventsProducer(t)
				mockEventsServiceProducer.On("ProduceEvent", mock.Anything).
					Once()

				return mockEventsServiceProducer
			},
			buildURLShortener: func() shortener.URLShortener {
				mockURLShortener := shortenermocks.NewURLShortener(t)
				mockURLShortener.On("ShortenURL", mock.AnythingOfType("uint32")).
					Return(testShortURL)

				return mockURLShortener
			},
			expectedShortURL: testShortURL,
			expectedErr:      nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			urlService := NewURLService(
				logger,
				tc.buildURLRepo(),
				tc.buildURLCache(),
				tc.buildEventsProducer(),
				tc.buildURLShortener(),
			)

			shortURL, err := urlService.SaveURL(context.Background(), testLongURL)
			assert.Equal(t, tc.expectedShortURL, shortURL)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
