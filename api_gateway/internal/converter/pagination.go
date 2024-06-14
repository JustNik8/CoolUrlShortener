package converter

import (
	"api_gateway/internal/transport/rest/dto"
	analytics "api_gateway/pkg/proto"
)

type PaginationConverter struct {
}

func NewPaginationConverter() PaginationConverter {
	return PaginationConverter{}
}

func (c *PaginationConverter) MapPbToDto(pb *analytics.Pagination) dto.Pagination {
	return dto.Pagination{
		Next:          int(pb.Next),
		Previous:      int(pb.Previous),
		RecordPerPage: int(pb.RecordPerPage),
		CurrentPage:   int(pb.CurrentPage),
		TotalPage:     int(pb.TotalPage),
	}
}
