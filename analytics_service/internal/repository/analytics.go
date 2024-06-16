package repository

import (
	"context"

	"analytics_service/internal/domain"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name AnalyticsRepo
type AnalyticsRepo interface {
	GetTopUrls(ctx context.Context, paginationParams domain.PaginationParams) ([]domain.TopURLData, error)
}
