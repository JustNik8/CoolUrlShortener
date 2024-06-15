package clickhouserepo

import (
	"context"
	"log/slog"

	"analytics_service/internal/domain"
	"analytics_service/internal/repository"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type analyticsRepoClickhouse struct {
	logger *slog.Logger
	conn   driver.Conn
}

func NewAnalyticsRepoClickhouse(
	logger *slog.Logger,
	conn driver.Conn,
) (repository.AnalyticsRepo, error) {
	//addr := fmt.Sprintf("%s:%s", host, port)
	//
	//conn, err := clickhouse.Open(&clickhouse.Options{
	//	Protocol: clickhouse.Native,
	//	Addr:     []string{addr},
	//	Auth: clickhouse.Auth{
	//		Database: database,
	//		Username: username,
	//		Password: password,
	//	},
	//	Debug:           true,
	//	DialTimeout:     30 * time.Second,
	//	MaxOpenConns:    10,
	//	MaxIdleConns:    5,
	//	ConnMaxLifetime: time.Hour,
	//})
	//
	//if err != nil {
	//	return nil, err
	//}

	return &analyticsRepoClickhouse{
		logger: logger,
		conn:   conn,
	}, nil
}

const getTopUrlsQuery = `select long_url, short_url, follow_count, create_count from url_events_counter FINAL 
ORDER BY (follow_count, create_count) DESC 
LIMIT $1
OFFSET $2;`

func (r *analyticsRepoClickhouse) GetTopUrls(
	ctx context.Context,
	paginationParams domain.PaginationParams,
) ([]domain.TopURLData, error) {
	offset := paginationParams.Limit * (paginationParams.Page - 1)

	rows, err := r.conn.Query(ctx, getTopUrlsQuery, paginationParams.Limit, offset)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			r.logger.Error(err.Error())
		}
	}()

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
