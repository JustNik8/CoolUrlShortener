package repository

import (
	"context"

	"CoolUrlShortener/internal/domain"
)

type UrlRepo interface {
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	GetShortURLByLongURL(ctx context.Context, longURL string) (string, error)
	SaveURL(ctx context.Context, urlData domain.URLData) error
}
