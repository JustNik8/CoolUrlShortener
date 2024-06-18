package client

import (
	"context"
	"log/slog"

	"api_gateway/errs"
	"api_gateway/internal/converter"
	"api_gateway/internal/transport/rest/dto"
	"api_gateway/pkg/proto/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name AnalyticsClient
type AnalyticsClient interface {
	GetTopUrls(ctx context.Context, page int64, limit int64) (dto.TopURLDataResponse, error)
}

type grpcAnalyticsClient struct {
	logger              *slog.Logger
	grpcClient          analytics.AnalyticsClient
	topUrlConverter     converter.TopURLConverter
	paginationConverter converter.PaginationConverter
}

func NewGrpcAnalyticsClient(
	logger *slog.Logger,
	grpcClient analytics.AnalyticsClient,
	topUrlConverter converter.TopURLConverter,
	paginationConverter converter.PaginationConverter,
) AnalyticsClient {
	return &grpcAnalyticsClient{
		logger:              logger,
		grpcClient:          grpcClient,
		topUrlConverter:     topUrlConverter,
		paginationConverter: paginationConverter,
	}
}

func (g *grpcAnalyticsClient) GetTopUrls(ctx context.Context, page int64, limit int64) (dto.TopURLDataResponse, error) {
	topUrlsGrpcResp, err := g.grpcClient.GetTopUrls(context.Background(), &analytics.TopUrlsRequest{
		Page:  page,
		Limit: limit,
	})

	if err != nil {
		g.logger.Error(err.Error())

		st, ok := status.FromError(err)
		if !ok || st.Code() == codes.Internal {
			return dto.TopURLDataResponse{}, errs.ErrInternal
		}

		if st.Code() == codes.InvalidArgument {
			return dto.TopURLDataResponse{}, errs.ErrInvalidArgument
		}

		return dto.TopURLDataResponse{}, errs.ErrInternal
	}

	topUrlsResp := dto.TopURLDataResponse{
		TopURLData: g.topUrlConverter.MapSlicePbToDto(topUrlsGrpcResp.TopUrlData),
		Pagination: g.paginationConverter.MapPbToDto(topUrlsGrpcResp.Pagination),
	}

	return topUrlsResp, nil
}
