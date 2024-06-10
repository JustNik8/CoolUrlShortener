package repository

import (
	"context"

	"CoolUrlShortener/internal/domain"
)

type UrlShortenerRepo interface {
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	SaveURL(ctx context.Context, urlData domain.URLData) error
}
