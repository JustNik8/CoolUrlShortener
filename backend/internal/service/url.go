package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/repository"
	"CoolUrlShortener/internal/repository/models"
	"CoolUrlShortener/pkg/shortener"
	"github.com/google/uuid"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name URLService
type URLService interface {
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	SaveURL(ctx context.Context, longURL string) (string, error)
}

type urlService struct {
	logger         *slog.Logger
	urlRepo        repository.UrlRepo
	urlCache       repository.URLCache
	eventsProducer repository.EventsProducer
	urlShortener   shortener.URLShortener
}

func NewURLService(
	logger *slog.Logger,
	repo repository.UrlRepo,
	urlCache repository.URLCache,
	eventsProducer repository.EventsProducer,
	urlShortener shortener.URLShortener,
) URLService {
	return &urlService{
		logger:         logger,
		urlRepo:        repo,
		urlCache:       urlCache,
		eventsProducer: eventsProducer,
		urlShortener:   urlShortener,
	}
}

func (s *urlService) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	longURLCache, err := s.urlCache.GetLongURL(ctx, shortURL)
	if err == nil {
		s.eventsProducer.ProduceEvent(
			models.URLEvent{
				LongURL:   longURLCache,
				ShortURL:  shortURL,
				EventTime: time.Now().Unix(),
				EventType: models.EventTypeFollow,
			},
		)
		return longURLCache, nil
	}

	longURL, err := s.urlRepo.GetLongURL(ctx, shortURL)
	if err != nil {
		return "", err
	}
	err = s.urlCache.SetLongURL(ctx, shortURL, longURL)
	if err != nil {
		s.logger.Error(err.Error())
	}

	s.eventsProducer.ProduceEvent(
		models.URLEvent{
			LongURL:   longURL,
			ShortURL:  shortURL,
			EventTime: time.Now().Unix(),
			EventType: models.EventTypeFollow,
		},
	)
	return longURL, nil
}

func (s *urlService) SaveURL(ctx context.Context, longURL string) (string, error) {
	gotShortURL, err := s.urlRepo.GetShortURLByLongURL(ctx, longURL)
	if err == nil {
		s.eventsProducer.ProduceEvent(
			models.URLEvent{
				LongURL:   longURL,
				ShortURL:  gotShortURL,
				EventTime: time.Now().Unix(),
				EventType: models.EventTypeCreate,
			},
		)
		return gotShortURL, nil
	}
	if err != nil && !errors.Is(err, errs.ErrNoURL) {
		return "", err
	}

	id := uuid.New().ID()
	shortUrl := s.urlShortener.ShortenURL(id)
	urlData := domain.URLData{
		ID:        int64(id),
		ShortUrl:  shortUrl,
		LongUrl:   longURL,
		CreatedAt: time.Now(),
	}

	err = s.urlRepo.SaveURL(ctx, urlData)
	if err != nil {
		return "", err
	}
	err = s.urlCache.SetLongURL(ctx, shortUrl, longURL)
	if err != nil {
		s.logger.Error(err.Error())
	}

	s.eventsProducer.ProduceEvent(
		models.URLEvent{
			LongURL:   longURL,
			ShortURL:  shortUrl,
			EventTime: time.Now().Unix(),
			EventType: models.EventTypeCreate,
		},
	)
	return shortUrl, nil
}
