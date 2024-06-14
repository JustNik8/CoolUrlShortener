package converter

import (
	"analytics_service/internal/domain"
	"analytics_service/internal/transport/rest/dto"
	analytics "analytics_service/pkg/proto"
)

type TopURLConverter struct {
}

func NewTopURLConverter() TopURLConverter {
	return TopURLConverter{}
}

func (c *TopURLConverter) MapDomainToDto(d domain.TopURLData) dto.TopURLData {
	return dto.TopURLData{
		LongURL:     d.LongURL,
		ShortURL:    d.ShortURL,
		FollowCount: d.FollowCount,
		CreateCount: d.CreateCount,
	}
}

func (c *TopURLConverter) MapSliceDomainToDto(d []domain.TopURLData) []dto.TopURLData {
	dtos := make([]dto.TopURLData, len(d))

	for i := 0; i < len(d); i++ {
		dtos[i] = c.MapDomainToDto(d[i])
	}

	return dtos
}

func (c *TopURLConverter) MapDomainToPb(d domain.TopURLData) *analytics.TopUrlData {
	return &analytics.TopUrlData{
		LongUrl:     d.LongURL,
		ShortUrl:    d.ShortURL,
		FollowCount: d.FollowCount,
		CreateCount: d.CreateCount,
	}
}

func (c *TopURLConverter) MapSliceDomainToPb(d []domain.TopURLData) []*analytics.TopUrlData {
	pbs := make([]*analytics.TopUrlData, len(d))

	for i := 0; i < len(d); i++ {
		pbs[i] = c.MapDomainToPb(d[i])
	}

	return pbs
}
