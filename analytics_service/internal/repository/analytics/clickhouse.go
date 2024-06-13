package analytics

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"analytics_service/internal/domain"
	"analytics_service/internal/repository"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type analyticsRepoClickhouse struct {
	logger *slog.Logger
	conn   driver.Conn
}

func NewAnalyticsRepoClickhouse(
	logger *slog.Logger,
	database string,
	username string,
	password string,
	host string,
	port string,
) (repository.AnalyticsRepo, error) {
	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Protocol: clickhouse.Native,
		Addr:     []string{addr},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     30 * time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})

	if err != nil {
		return nil, err
	}

	return &analyticsRepoClickhouse{
		logger: logger,
		conn:   conn,
	}, nil
}

const getTopUrlsQuery = `select long_url, short_url, follow_count, create_count
from url_events_counter FINAL
ORDER BY (follow_count, create_count) DESC
LIMIT {limit:Int64};`

func (r *analyticsRepoClickhouse) GetTopUrls(ctx context.Context, limit int) ([]domain.TopURLData, error) {
	chCtx := clickhouse.Context(ctx, clickhouse.WithParameters(clickhouse.Parameters{
		"limit": strconv.Itoa(limit),
	}))

	rows, err := r.conn.Query(chCtx, getTopUrlsQuery)
	defer func() {
		err := rows.Close()
		if err != nil {
			r.logger.Error(err.Error())
		}
	}()
	if err != nil {
		return nil, err
	}

	topURLs := make([]domain.TopURLData, 0)
	for rows.Next() {
		var urlData domain.TopURLData
		err = rows.Scan(&urlData.LongURL, &urlData.ShortURL, &urlData.FollowCount, &urlData.CreateCount)
		if err != nil {
			r.logger.Error(err.Error())
			continue
		}

		topURLs = append(topURLs, urlData)
	}

	return topURLs, nil
}
