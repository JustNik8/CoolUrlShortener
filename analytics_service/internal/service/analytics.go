package service

import (
	"context"

	"analytics_service/internal/domain"
	"analytics_service/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name AnalyticsService
type AnalyticsService interface {
	GetTopUrls(ctx context.Context, paginationParams domain.PaginationParams) ([]domain.TopURLData, error)
}

type analyticsService struct {
	analyticsRepo repository.AnalyticsRepo
}

func NewAnalyticsService(
	analyticsRepo repository.AnalyticsRepo,
) AnalyticsService {
	return &analyticsService{
		analyticsRepo: analyticsRepo,
	}
}

func (s *analyticsService) GetTopUrls(ctx context.Context, paginationParams domain.PaginationParams) ([]domain.TopURLData, error) {
	return s.analyticsRepo.GetTopUrls(ctx, paginationParams)
}
