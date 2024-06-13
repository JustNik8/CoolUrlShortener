package service

import (
	"context"

	"analytics_service/internal/domain"
	"analytics_service/internal/repository"
)

type AnalyticsService interface {
	GetTopUrls(ctx context.Context, limit int) ([]domain.TopURLData, error)
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

func (s *analyticsService) GetTopUrls(ctx context.Context, limit int) ([]domain.TopURLData, error) {
	return s.analyticsRepo.GetTopUrls(ctx, limit)
}
