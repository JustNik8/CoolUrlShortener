package postgresql

import (
	"context"
	"errors"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type urlRepoPostgres struct {
	dbPool *pgxpool.Pool
}

func NewUrlRepoPostgres(
	dbPool *pgxpool.Pool,
) repository.UrlRepo {
	return &urlRepoPostgres{
		dbPool: dbPool,
	}
}

const getLongURLQuery = `SELECT long_url FROM url_data WHERE short_url = $1`

func (r *urlRepoPostgres) GetLongURL(ctx context.Context, shortUrl string) (string, error) {
	var longURL string
	row := r.dbPool.QueryRow(ctx, getLongURLQuery, shortUrl)

	err := row.Scan(&longURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errs.ErrNoURL
	}

	return longURL, err
}

const saveURLQuery = `INSERT INTO url_data (id, short_url, long_url, created_at) 
VALUES ($1, $2, $3, $4)`

const getShortURLByLongURL = `SELECT short_url FROM url_data WHERE long_url = $1`

func (r *urlRepoPostgres) GetShortURLByLongURL(ctx context.Context, longURL string) (string, error) {
	var shortURL string
	row := r.dbPool.QueryRow(ctx, getShortURLByLongURL, longURL)

	err := row.Scan(&shortURL)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", errs.ErrNoURL
	}

	return shortURL, err
}

func (r *urlRepoPostgres) SaveURL(ctx context.Context, urlData domain.URLData) error {
	_, err := r.dbPool.Exec(ctx, saveURLQuery, urlData.ID, urlData.ShortUrl, urlData.LongUrl, urlData.CreatedAt)
	return err
}
