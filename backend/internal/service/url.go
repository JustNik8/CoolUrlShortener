package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/repository"
	"github.com/google/uuid"
)

const (
	alphabet          = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	protocolSeparator = "://"
	defaultProtocol   = "https"
)

var (
	alphabetLen = len(alphabet)
)

var (
	availableProtocols = map[string]struct{}{
		"http":  {},
		"https": {},
	}
)

type URLService interface {
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	SaveURL(ctx context.Context, longURL string) (string, error)
	shortenURL(id uint32) string
}

type urlService struct {
	logger                *slog.Logger
	urlRepo               repository.UrlRepo
	urlCache              repository.URLCache
	eventsServiceProducer EventsServiceProducer
}

func NewURLService(
	logger *slog.Logger,
	repo repository.UrlRepo,
	urlCache repository.URLCache,
	eventsServiceProducer EventsServiceProducer,
) URLService {
	return &urlService{
		logger:                logger,
		urlRepo:               repo,
		urlCache:              urlCache,
		eventsServiceProducer: eventsServiceProducer,
	}
}

func (s *urlService) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	longURLCache, err := s.urlCache.GetLongURL(ctx, shortURL)
	if err == nil {
		s.eventsServiceProducer.ProduceEvent(
			domain.URLEvent{
				LongURL:   longURLCache,
				ShortURL:  shortURL,
				EventTime: time.Now(),
				EventType: eventTypeFollow,
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

	s.eventsServiceProducer.ProduceEvent(
		domain.URLEvent{
			LongURL:   longURL,
			ShortURL:  shortURL,
			EventTime: time.Now(),
			EventType: eventTypeFollow,
		},
	)
	return longURL, nil
}

func (s *urlService) SaveURL(ctx context.Context, longURL string) (string, error) {
	gotShortURL, err := s.urlRepo.GetShortURLByLongURL(ctx, longURL)
	if err == nil {
		s.eventsServiceProducer.ProduceEvent(
			domain.URLEvent{
				LongURL:   longURL,
				ShortURL:  gotShortURL,
				EventTime: time.Now(),
				EventType: eventTypeCreate,
			},
		)
		return gotShortURL, nil
	}
	if err != nil && !errors.Is(err, errs.ErrNoURL) {
		return "", err
	}

	id := uuid.New().ID()
	shortUrl := s.shortenURL(id)
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

	s.eventsServiceProducer.ProduceEvent(
		domain.URLEvent{
			LongURL:   longURL,
			ShortURL:  shortUrl,
			EventTime: time.Now(),
			EventType: eventTypeCreate,
		},
	)
	return shortUrl, nil
}

func (s *urlService) shortenURL(id uint32) string {
	nums := make([]int, 0)
	for id > 0 {
		rem := int(id) % alphabetLen
		nums = append(nums, rem)

		id /= uint32(alphabetLen)
	}

	var sb strings.Builder
	for i := range nums {
		idx := nums[i]
		sb.WriteByte(alphabet[idx])
	}

	return sb.String()
}
