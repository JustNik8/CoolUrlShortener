package converter

import (
	"analytics_service/internal/domain"
	"analytics_service/internal/transport/rest/dto"
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
