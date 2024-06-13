package repository

type PaginationRepo interface {
	GetRecordsCount(table string) (int, error)
}
