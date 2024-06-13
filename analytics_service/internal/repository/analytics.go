package repository

import (
	"context"

	"analytics_service/internal/domain"
)

type AnalyticsRepo interface {
	GetTopUrls(ctx context.Context, limit int) ([]domain.TopURLData, error)
}
