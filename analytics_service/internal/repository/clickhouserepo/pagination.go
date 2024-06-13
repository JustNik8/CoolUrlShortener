package clickhouserepo

import (
	"context"
	"fmt"

	"analytics_service/internal/repository"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type paginationRepoClickhouse struct {
	conn driver.Conn
}

func NewPaginationRepoClickhouse(conn driver.Conn) repository.PaginationRepo {
	return &paginationRepoClickhouse{
		conn: conn,
	}
}

func (r *paginationRepoClickhouse) GetRecordsCount(table string) (int, error) {
	sqlTableQuery := fmt.Sprintf("SELECT count(*) FROM %s FINAL", table)
	row := r.conn.QueryRow(context.Background(), sqlTableQuery)

	var recordsCount uint64
	err := row.Scan(&recordsCount)
	if err != nil {
		return 0, err
	}

	return int(recordsCount), nil
}
