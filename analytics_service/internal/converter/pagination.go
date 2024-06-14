package converter

import (
	"analytics_service/internal/domain"
	analytics "analytics_service/pkg/proto"
)

type PaginationConverter struct {
}

func NewPaginationConverter() PaginationConverter {
	return PaginationConverter{}
}

func (c *PaginationConverter) MapDomainToPb(d domain.Pagination) *analytics.Pagination {
	return &analytics.Pagination{
		Next:          int64(d.Next),
		Previous:      int64(d.Previous),
		RecordPerPage: int64(d.RecordPerPage),
		CurrentPage:   int64(d.CurrentPage),
		TotalPage:     int64(d.TotalPage),
	}
}
