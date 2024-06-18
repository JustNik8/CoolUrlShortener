package rest

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"api_gateway/internal/client"
	"api_gateway/internal/client/mocks"
	"api_gateway/internal/transport/rest/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTopURLs(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	basePath := "/api/top_urls"

	testTopUrlDataResp := dto.TopURLDataResponse{
		TopURLData: []dto.TopURLData{
			{LongURL: "http://test.long", ShortURL: "short", FollowCount: 10, CreateCount: 1},
			{LongURL: "http://test.long2", ShortURL: "short2", FollowCount: 20, CreateCount: 2},
			{LongURL: "http://test.long3", ShortURL: "short3", FollowCount: 30, CreateCount: 3},
		},
		Pagination: dto.Pagination{
			Next:          2,
			Previous:      0,
			RecordPerPage: 3,
			CurrentPage:   1,
			TotalPage:     10,
		},
	}
	testErr := errors.New("test error")

	testCases := []struct {
		name                 string
		buildAnalyticsClient func() client.AnalyticsClient
		page                 string
		limit                string
		expectedCode         int
	}{
		{
			name: "Get top urls without error. 200 OK",
			buildAnalyticsClient: func() client.AnalyticsClient {
				mockClient := mocks.NewAnalyticsClient(t)
				mockClient.On("GetTopUrls", mock.Anything, mock.Anything, mock.Anything).
					Return(testTopUrlDataResp, nil)

				return mockClient
			},
			page:         "",
			limit:        "",
			expectedCode: http.StatusOK,
		},
		{
			name: "Get top urls when internal error happened. 500 Internal Server Error",
			buildAnalyticsClient: func() client.AnalyticsClient {
				mockClient := mocks.NewAnalyticsClient(t)
				mockClient.On("GetTopUrls", mock.Anything, mock.Anything, mock.Anything).
					Return(dto.TopURLDataResponse{}, testErr)

				return mockClient
			},
			page:         "",
			limit:        "",
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Invalid page. 400 Bad Request",
			buildAnalyticsClient: func() client.AnalyticsClient {
				mockClient := mocks.NewAnalyticsClient(t)

				return mockClient
			},
			page:         "test",
			limit:        "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Invalid limit. 400 Bad Request",
			buildAnalyticsClient: func() client.AnalyticsClient {
				mockClient := mocks.NewAnalyticsClient(t)

				return mockClient
			},
			page:         "",
			limit:        "test",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewAnalyticsHandler(
				logger,
				tc.buildAnalyticsClient(),
			)

			req := httptest.NewRequest(http.MethodGet, basePath, nil)
			q := req.URL.Query()
			if tc.page != "" {
				q.Add("page", tc.page)
			}
			if tc.limit != "" {
				q.Add("limit", tc.limit)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()

			handler.GetTopURLs(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
