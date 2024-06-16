package grpc

import (
	"context"
	"log/slog"

	"analytics_service/internal/converter"
	"analytics_service/internal/domain"
	"analytics_service/internal/service"
	analytics "analytics_service/pkg/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	urlEventsCounterTableName = "url_events_counter"
)

type AnalyticsServer struct {
	logger              *slog.Logger
	analyticsService    service.AnalyticsService
	paginationService   service.PaginationService
	topURLConverter     converter.TopURLConverter
	paginationConverter converter.PaginationConverter
	analytics.UnimplementedAnalyticsServer
}

func NewAnalyticsServer(
	logger *slog.Logger,
	analyticsService service.AnalyticsService,
	paginationService service.PaginationService,
	topURLConverter converter.TopURLConverter,
	paginationConverter converter.PaginationConverter,
) *AnalyticsServer {
	return &AnalyticsServer{
		logger:              logger,
		analyticsService:    analyticsService,
		paginationService:   paginationService,
		topURLConverter:     topURLConverter,
		paginationConverter: paginationConverter,
	}
}

func (s *AnalyticsServer) GetTopUrls(
	ctx context.Context,
	req *analytics.TopUrlsRequest,
) (*analytics.TopUrlsResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	paginationParams := domain.PaginationParams{
		Page:  int(req.Page),
		Limit: int(req.Limit),
	}

	topUrls, err := s.analyticsService.GetTopUrls(ctx, paginationParams)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	pagination, err := s.paginationService.GetPaginationInfo(urlEventsCounterTableName, paginationParams)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &analytics.TopUrlsResponse{
		TopUrlData: s.topURLConverter.MapSliceDomainToPb(topUrls),
		Pagination: s.paginationConverter.MapDomainToPb(pagination),
	}, nil
}
