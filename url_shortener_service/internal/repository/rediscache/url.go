package rediscache

import (
	"context"
	"time"

	"CoolUrlShortener/internal/repository"
	"github.com/redis/go-redis/v9"
)

type urlCacheRedis struct {
	client *redis.Client
}

func NewURLCacheRedis(client *redis.Client) repository.URLCache {
	return &urlCacheRedis{
		client: client,
	}
}

func (u *urlCacheRedis) SetLongURL(ctx context.Context, shortURL string, longURL string) error {
	return u.client.Set(ctx, shortURL, longURL, 10*time.Minute).Err()
}

func (u *urlCacheRedis) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	longURL, err := u.client.Get(ctx, shortURL).Result()
	return longURL, err
}
