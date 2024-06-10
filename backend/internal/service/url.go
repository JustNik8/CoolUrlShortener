package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/repository"
)

const (
	alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

var (
	alphabetLen = len(alphabet)
)

type URLService interface {
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	SaveURL(ctx context.Context, longURL string) (string, error)
	shortenURL(longURL string) (int64, string)
}

type urlService struct {
	repo repository.UrlRepo
}

func NewURLService(repo repository.UrlRepo) URLService {
	return &urlService{
		repo: repo,
	}
}

func (s *urlService) GetLongURL(ctx context.Context, shortUrl string) (string, error) {
	return s.repo.GetLongURL(ctx, shortUrl)
}

func (s *urlService) SaveURL(ctx context.Context, longURL string) (string, error) {
	gotShortURL, err := s.repo.GetShortURLByLongURL(ctx, longURL)
	if err == nil {
		return gotShortURL, nil
	}
	if err != nil && !errors.Is(err, errs.ErrNoURL) {
		return "", err
	}

	id, shortUrl := s.shortenURL(longURL)
	urlData := domain.URLData{
		ID:        id,
		ShortUrl:  shortUrl,
		LongUrl:   longURL,
		CreatedAt: time.Now(),
	}

	err = s.repo.SaveURL(ctx, urlData)
	return shortUrl, err
}

func (s *urlService) shortenURL(longURL string) (int64, string) {
	sum := int(0)

	for _, r := range longURL {
		sum += int(r)
	}

	id := int64(sum)
	nums := make([]int, 0)
	for sum > 0 {
		rem := sum % alphabetLen
		nums = append(nums, rem)

		rem /= alphabetLen
	}

	var sb strings.Builder
	for i := range nums {
		sb.WriteByte(alphabet[i])
	}

	return id, sb.String()
}
