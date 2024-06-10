package postgresql

import (
	"context"
	"log/slog"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type urlShortenerRepoPostgres struct {
	logger *slog.Logger
	dbPool *pgxpool.Pool
}

func NewUrlShortenerRepoPostgres(
	dbPool *pgxpool.Pool,
) repository.UrlShortenerRepo {
	return &urlShortenerRepoPostgres{
		dbPool: dbPool,
	}
}

const getLongURLQuery = `SELECT long_url FROM url_data WHERE short_url = $1`

func (r *urlShortenerRepoPostgres) GetLongURL(ctx context.Context, shortUrl string) (string, error) {
	var longUrl string
	row := r.dbPool.QueryRow(ctx, getLongURLQuery, shortUrl)

	err := row.Scan(&longUrl)
	return longUrl, err
}

const saveURLQuery = `INSERT INTO url_data (id, short_url, long_url, created_at) 
VALUES ($1, $2, $3, $4)`

func (r *urlShortenerRepoPostgres) SaveURL(ctx context.Context, urlData domain.URLData) error {
	_, err := r.dbPool.Exec(ctx, saveURLQuery, urlData.ID, urlData.ShortUrl, urlData.LongUrl, urlData.CreatedAt)
	return err
}
