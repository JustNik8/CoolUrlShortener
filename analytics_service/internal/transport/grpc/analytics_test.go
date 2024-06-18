package grpc

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"

	"analytics_service/internal/converter"
	"analytics_service/internal/domain"
	"analytics_service/internal/service"
	"analytics_service/internal/service/mocks"
	analytics "analytics_service/pkg/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func initAnalyticsClient(
	logger *slog.Logger,
	analyticsService service.AnalyticsService,
	paginationService service.PaginationService,
) (analytics.AnalyticsClient, func()) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	topUrlConverter := converter.NewTopURLConverter()
	paginationConverter := converter.NewPaginationConverter()

	analyticsServer := NewAnalyticsServer(
		logger,
		analyticsService,
		paginationService,
		topUrlConverter,
		paginationConverter,
	)

	baseServer := grpc.NewServer()

	analytics.RegisterAnalyticsServer(baseServer, analyticsServer)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("Server exited with error: %v", err)
		}
	}()

	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	transportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient("passthrough://bufnet", grpc.WithContextDialer(bufDialer), transportOpt)
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		err = conn.Close()
		if err != nil {
			log.Printf("error closing conn: %v", err)
		}
		baseServer.Stop()
	}

	client := analytics.NewAnalyticsClient(conn)

	return client, closer
}

func TestGetTopUrls(t *testing.T) {
	testPaginationParams := domain.PaginationParams{Page: 1, Limit: 3}
	testPaginationParamsReq := &analytics.TopUrlsRequest{Page: 1, Limit: 3}

	testTopUrls := []domain.TopURLData{
		{LongURL: "http://test.long1", ShortURL: "test", FollowCount: 10, CreateCount: 1},
		{LongURL: "http://test.long2", ShortURL: "test2", FollowCount: 20, CreateCount: 2},
		{LongURL: "http://test.long3", ShortURL: "tes3", FollowCount: 30, CreateCount: 3},
	}

	testTopUrlsResp := []*analytics.TopUrlData{
		{LongUrl: "http://test.long1", ShortUrl: "test", FollowCount: 10, CreateCount: 1},
		{LongUrl: "http://test.long2", ShortUrl: "test2", FollowCount: 20, CreateCount: 2},
		{LongUrl: "http://test.long3", ShortUrl: "tes3", FollowCount: 30, CreateCount: 3},
	}

	testPagination := domain.Pagination{
		Next: 2, RecordPerPage: 3, CurrentPage: 1, TotalPage: 10,
	}
	testPaginationResp := &analytics.Pagination{
		Next: 2, RecordPerPage: 3, CurrentPage: 1, TotalPage: 10,
	}

	testErr := errors.New("test error")

	testCases := []struct {
		name                   string
		buildAnalyticsService  func() service.AnalyticsService
		buildPaginationService func() service.PaginationService
		request                *analytics.TopUrlsRequest
		expectedResp           *analytics.TopUrlsResponse
		isErrExpected          bool
		expectedCode           codes.Code
	}{
		{
			name: "test get top urls without error",
			buildAnalyticsService: func() service.AnalyticsService {
				mockService := mocks.NewAnalyticsService(t)
				mockService.On("GetTopUrls", mock.Anything, mock.Anything).
					Return(testTopUrls, nil)

				return mockService
			},
			buildPaginationService: func() service.PaginationService {
				mockService := mocks.NewPaginationService(t)
				mockService.On("GetPaginationInfo", mock.Anything, testPaginationParams).
					Return(testPagination, nil)

				return mockService
			},
			request: testPaginationParamsReq,
			expectedResp: &analytics.TopUrlsResponse{
				TopUrlData: testTopUrlsResp,
				Pagination: testPaginationResp,
			},
			isErrExpected: false,
			expectedCode:  codes.OK,
		},
		{
			name: "Given empty page should return error. 3 Invalid Argument",
			buildAnalyticsService: func() service.AnalyticsService {
				mockService := mocks.NewAnalyticsService(t)
				return mockService
			},
			buildPaginationService: func() service.PaginationService {
				mockService := mocks.NewPaginationService(t)
				return mockService
			},
			request:       &analytics.TopUrlsRequest{Page: 0, Limit: 10},
			expectedResp:  &analytics.TopUrlsResponse{},
			isErrExpected: true,
			expectedCode:  codes.InvalidArgument,
		},
		{
			name: "Given empty limit should return error. 3 Invalid Argument",
			buildAnalyticsService: func() service.AnalyticsService {
				mockService := mocks.NewAnalyticsService(t)
				return mockService
			},
			buildPaginationService: func() service.PaginationService {
				mockService := mocks.NewPaginationService(t)
				return mockService
			},
			request:       &analytics.TopUrlsRequest{Page: 1, Limit: 0},
			expectedResp:  &analytics.TopUrlsResponse{},
			isErrExpected: true,
			expectedCode:  codes.InvalidArgument,
		},
		{
			name: "internal error when get top urls. 13 Internal",
			buildAnalyticsService: func() service.AnalyticsService {
				mockService := mocks.NewAnalyticsService(t)
				mockService.On("GetTopUrls", mock.Anything, mock.Anything).
					Return(nil, testErr)

				return mockService
			},
			buildPaginationService: func() service.PaginationService {
				mockService := mocks.NewPaginationService(t)
				return mockService
			},
			request:       &analytics.TopUrlsRequest{Page: 1, Limit: 10},
			expectedResp:  &analytics.TopUrlsResponse{},
			isErrExpected: true,
			expectedCode:  codes.Internal,
		},
		{
			name: "internal error when get pagination. 13 Internal",
			buildAnalyticsService: func() service.AnalyticsService {
				mockService := mocks.NewAnalyticsService(t)
				mockService.On("GetTopUrls", mock.Anything, mock.Anything).
					Return(testTopUrls, nil)

				return mockService
			},
			buildPaginationService: func() service.PaginationService {
				mockService := mocks.NewPaginationService(t)
				mockService.On("GetPaginationInfo", mock.Anything, mock.Anything).
					Return(domain.Pagination{}, testErr)

				return mockService
			},
			request:       &analytics.TopUrlsRequest{Page: 1, Limit: 10},
			expectedResp:  &analytics.TopUrlsResponse{},
			isErrExpected: true,
			expectedCode:  codes.Internal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)

			analyticsClient, cancel := initAnalyticsClient(
				logger,
				tc.buildAnalyticsService(),
				tc.buildPaginationService(),
			)
			defer cancel()

			resp, err := analyticsClient.GetTopUrls(context.Background(), tc.request)
			isErrorHappened := err != nil

			assert.Equal(t, tc.isErrExpected, isErrorHappened)
			if tc.isErrExpected {
				st, ok := status.FromError(err)

				assert.Equal(t, ok, true)
				assert.Equal(t, tc.expectedCode, st.Code())
				return
			}

			assert.Equal(t, len(tc.expectedResp.TopUrlData), len(resp.TopUrlData))

			for i := range resp.TopUrlData {
				expectedData := tc.expectedResp.TopUrlData[i]
				actualData := resp.TopUrlData[i]

				assert.Equal(t, expectedData.LongUrl, actualData.LongUrl)
				assert.Equal(t, expectedData.ShortUrl, actualData.ShortUrl)
				assert.Equal(t, expectedData.FollowCount, actualData.FollowCount)
				assert.Equal(t, expectedData.CreateCount, actualData.CreateCount)
			}

			assert.Equal(t, tc.expectedResp.Pagination.TotalPage, resp.Pagination.TotalPage)
			assert.Equal(t, tc.expectedResp.Pagination.Next, resp.Pagination.Next)
			assert.Equal(t, tc.expectedResp.Pagination.CurrentPage, resp.Pagination.CurrentPage)
			assert.Equal(t, tc.expectedResp.Pagination.RecordPerPage, resp.Pagination.RecordPerPage)
			assert.Equal(t, tc.expectedResp.Pagination.Previous, resp.Pagination.Previous)
		})
	}
}
