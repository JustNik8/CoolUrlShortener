package repository

import "context"

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name URLCache
type URLCache interface {
	SetLongURL(ctx context.Context, shortURL string, longURL string) error
	GetLongURL(ctx context.Context, shortURL string) (string, error)
}
