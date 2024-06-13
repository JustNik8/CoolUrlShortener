package domain

type PaginationParams struct {
	Page  int
	Limit int
}

type Pagination struct {
	Next          int
	Previous      int
	RecordPerPage int
	CurrentPage   int
	TotalPage     int
}
