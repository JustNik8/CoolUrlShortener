package repository

import "context"

type URLCache interface {
	SetLongURL(ctx context.Context, shortURL string, longURL string) error
	GetLongURL(ctx context.Context, shortURL string) (string, error)
}
