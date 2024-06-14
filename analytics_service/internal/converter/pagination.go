package converter

import (
	"analytics_service/internal/domain"
	"analytics_service/internal/transport/rest/dto"
	analytics "analytics_service/pkg/proto"
)

type PaginationConverter struct {
}

func NewPaginationConverter() PaginationConverter {
	return PaginationConverter{}
}

func (c *PaginationConverter) MapDomainToDto(d domain.Pagination) dto.Pagination {
	return dto.Pagination{
		Next:          d.Next,
		Previous:      d.Previous,
		RecordPerPage: d.RecordPerPage,
		CurrentPage:   d.CurrentPage,
		TotalPage:     d.TotalPage,
	}
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
