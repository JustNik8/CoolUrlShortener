package repository

import (
	"context"

	"CoolUrlShortener/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name UrlRepo
type UrlRepo interface {
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	GetShortURLByLongURL(ctx context.Context, longURL string) (string, error)
	SaveURL(ctx context.Context, urlData domain.URLData) error
}
