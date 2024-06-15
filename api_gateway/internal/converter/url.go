package converter

import (
	"api_gateway/internal/transport/rest/dto"
	"api_gateway/pkg/proto/analytics"
)

type TopURLConverter struct {
}

func NewTopURLConverter() TopURLConverter {
	return TopURLConverter{}
}

func (c *TopURLConverter) MapPbToDto(pb *analytics.TopUrlData) dto.TopURLData {
	return dto.TopURLData{
		LongURL:     pb.LongUrl,
		ShortURL:    pb.ShortUrl,
		FollowCount: pb.FollowCount,
		CreateCount: pb.CreateCount,
	}
}

func (c *TopURLConverter) MapSlicePbToDto(pbs []*analytics.TopUrlData) []dto.TopURLData {
	dtos := make([]dto.TopURLData, len(pbs))

	for i := 0; i < len(pbs); i++ {
		dtos[i] = c.MapPbToDto(pbs[i])
	}

	return dtos
}
