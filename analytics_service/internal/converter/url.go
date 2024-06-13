package converter

import (
	"analytics_service/internal/domain"
	"analytics_service/internal/transport/rest/dto"
)

type TopURLConverter struct {
}

func NewTopURLConverter() *TopURLConverter {
	return &TopURLConverter{}
}

func (c *TopURLConverter) ConvertDomainToDto(d domain.TopURLData) dto.TopURLData {
	return dto.TopURLData{
		LongURL:     d.LongURL,
		ShortURL:    d.ShortURL,
		FollowCount: d.FollowCount,
		CreateCount: d.CreateCount,
	}
}

func (c *TopURLConverter) ConvertSliceDomainToDto(d []domain.TopURLData) []dto.TopURLData {
	dtos := make([]dto.TopURLData, len(d))

	for i := 0; i < len(d); i++ {
		dtos[i] = c.ConvertDomainToDto(d[i])
	}

	return dtos
}
