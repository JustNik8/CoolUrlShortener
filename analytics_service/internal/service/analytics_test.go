package service

import (
	"context"
	"errors"
	"testing"

	"analytics_service/internal/domain"
	"analytics_service/internal/repository"
	"analytics_service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTopUrls(t *testing.T) {
	testTopUrlData := []domain.TopURLData{
		{LongURL: "http://test.long", ShortURL: "test", FollowCount: 10, CreateCount: 2},
	}
	errTest := errors.New("test error")

	testCases := []struct {
		name               string
		buildAnalyticsRepo func() repository.AnalyticsRepo
		paginationParams   domain.PaginationParams
		expectedUrlData    []domain.TopURLData
		expectedErr        error
	}{
		{
			name: "get top urls without error",
			buildAnalyticsRepo: func() repository.AnalyticsRepo {
				mockRepo := mocks.NewAnalyticsRepo(t)
				mockRepo.On("GetTopUrls", mock.Anything, mock.Anything).
					Return(testTopUrlData, nil)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 1, Limit: 10},
			expectedUrlData:  testTopUrlData,
			expectedErr:      nil,
		},
		{
			name: "get top urls error occurred",
			buildAnalyticsRepo: func() repository.AnalyticsRepo {
				mockRepo := mocks.NewAnalyticsRepo(t)
				mockRepo.On("GetTopUrls", mock.Anything, mock.Anything).
					Return(nil, errTest)

				return mockRepo
			},
			paginationParams: domain.PaginationParams{Page: 1, Limit: 10},
			expectedUrlData:  nil,
			expectedErr:      errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyticsService := NewAnalyticsService(tc.buildAnalyticsRepo())

			urlData, err := analyticsService.GetTopUrls(context.Background(), tc.paginationParams)
			assert.Equal(t, tc.expectedUrlData, urlData)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
