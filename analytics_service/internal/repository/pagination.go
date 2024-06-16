package repository

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name PaginationRepo
type PaginationRepo interface {
	GetRecordsCount(table string) (int, error)
}
