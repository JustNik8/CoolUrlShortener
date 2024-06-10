package events

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/repository"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type eventsWriterClickhouse struct {
	logger *slog.Logger
	conn   driver.Conn
}

func NewEventsWriterClickhouse(
	logger *slog.Logger,
	database string,
	username string,
	password string,
	host string,
	port string,
) (repository.EventsWriter, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Println(addr)

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

	v, err := conn.ServerVersion()
	log.Println(v)
	log.Println(err)

	return &eventsWriterClickhouse{
		logger: logger,
		conn:   conn,
	}, nil
}

const insertQuery = `INSERT INTO url_events`

func (c *eventsWriterClickhouse) Insert(events []domain.URLEvent) error {
	ctx := context.Background()

	batch, err := c.conn.PrepareBatch(ctx, insertQuery)
	if err != nil {
		return err
	}

	for _, event := range events {
		err = batch.Append(
			event.LongURL,
			event.ShortURL,
			event.EventTime,
			event.EventType,
		)

		if err != nil {
			c.logger.Error(err.Error())
		}
	}

	return batch.Send()
}
