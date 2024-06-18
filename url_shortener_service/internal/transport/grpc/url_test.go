package grpc

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"os"
	"testing"

	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/service"
	"CoolUrlShortener/internal/service/mocks"
	url "CoolUrlShortener/pkg/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func initUrlClient(
	logger *slog.Logger,
	urlService service.URLService,
) (url.UrlClient, func()) {
	const bufSize = 1024 * 1024
	lis := bufconn.Listen(bufSize)

	urlServer := NewUrlServer(
		logger, urlService,
	)

	baseServer := grpc.NewServer()

	url.RegisterUrlServer(baseServer, urlServer)
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

	client := url.NewUrlClient(conn)

	return client, closer
}

func TestShortenUrl(t *testing.T) {
	testLongUrl := "http://test.long"
	testShortUrl := "short"
	testErr := errors.New("test error")

	testCases := []struct {
		name            string
		buildUrlService func() service.URLService
		request         *url.LongUrlRequest
		expectedResp    *url.UrlDataResponse
		isErrExpected   bool
		expectedCode    codes.Code
	}{
		{
			name: "short url without error. 0 OK",
			buildUrlService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("SaveURL", mock.Anything, mock.Anything).
					Return(testShortUrl, nil)

				return mockService
			},
			request: &url.LongUrlRequest{
				LongUrl: testLongUrl,
			},
			expectedResp: &url.UrlDataResponse{
				LongUrl:  testLongUrl,
				ShortUrl: testShortUrl,
			},
			isErrExpected: false,
			expectedCode:  codes.OK,
		},
		{
			name: "short url is empty should return error. 3 InvalidArgument",
			buildUrlService: func() service.URLService {
				mockService := mocks.NewURLService(t)

				return mockService
			},
			request: &url.LongUrlRequest{
				LongUrl: "",
			},
			expectedResp:  &url.UrlDataResponse{},
			isErrExpected: true,
			expectedCode:  codes.InvalidArgument,
		},
		{
			name: "shorten url with internal error while save url. 13 Internal",
			buildUrlService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("SaveURL", mock.Anything, mock.Anything).
					Return("", testErr)

				return mockService
			},
			request: &url.LongUrlRequest{
				LongUrl: testLongUrl,
			},
			expectedResp:  &url.UrlDataResponse{},
			isErrExpected: true,
			expectedCode:  codes.Internal,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)

			urlClient, cancel := initUrlClient(logger, tc.buildUrlService())
			defer cancel()

			resp, err := urlClient.ShortenUrl(context.Background(), tc.request)
			isErrorHappened := err != nil

			assert.Equal(t, tc.isErrExpected, isErrorHappened)
			if tc.isErrExpected {
				st, ok := status.FromError(err)

				assert.Equal(t, ok, true)
				assert.Equal(t, tc.expectedCode, st.Code())
				return
			}

			assert.Equal(t, tc.expectedResp.LongUrl, resp.LongUrl)
			assert.Equal(t, tc.expectedResp.ShortUrl, resp.ShortUrl)

		})
	}
}

func TestFollowUrl(t *testing.T) {
	testLongUrl := "http://test.long"
	testShortUrl := "short"
	testErr := errors.New("test error")

	testCases := []struct {
		name            string
		buildUrlService func() service.URLService
		request         *url.ShortUrlRequest
		expectedResp    *url.LongUrlResponse
		isErrExpected   bool
		expectedCode    codes.Code
	}{
		{
			name: "get long url without error. 0 OK",
			buildUrlService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("GetLongURL", mock.Anything, mock.Anything).
					Return(testLongUrl, nil)

				return mockService
			},
			request:       &url.ShortUrlRequest{ShortUrl: testShortUrl},
			expectedResp:  &url.LongUrlResponse{LongUrl: testLongUrl},
			isErrExpected: false,
			expectedCode:  codes.OK,
		},
		{
			name: "url not found . 5 Not found",
			buildUrlService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("GetLongURL", mock.Anything, mock.Anything).
					Return("", errs.ErrNoURL)

				return mockService
			},
			request:       &url.ShortUrlRequest{ShortUrl: testShortUrl},
			expectedResp:  &url.LongUrlResponse{},
			isErrExpected: true,
			expectedCode:  codes.NotFound,
		},
		{
			name: "get long url while internal error. 13 Internal",
			buildUrlService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("GetLongURL", mock.Anything, mock.Anything).
					Return("", testErr)

				return mockService
			},
			request:       &url.ShortUrlRequest{ShortUrl: testShortUrl},
			expectedResp:  &url.LongUrlResponse{},
			isErrExpected: true,
			expectedCode:  codes.Internal,
		},
		{
			name: "pass empty short url should be error. 3 InvalidArgument",
			buildUrlService: func() service.URLService {
				mockService := mocks.NewURLService(t)

				return mockService
			},
			request:       &url.ShortUrlRequest{ShortUrl: ""},
			expectedResp:  &url.LongUrlResponse{},
			isErrExpected: true,
			expectedCode:  codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger := slog.New(
				slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
			)

			urlClient, cancel := initUrlClient(logger, tc.buildUrlService())
			defer cancel()

			resp, err := urlClient.FollowUrl(context.Background(), tc.request)
			isErrorHappened := err != nil

			assert.Equal(t, tc.isErrExpected, isErrorHappened)
			if tc.isErrExpected {
				st, ok := status.FromError(err)

				assert.Equal(t, ok, true)
				assert.Equal(t, tc.expectedCode, st.Code())
				return
			}

			assert.Equal(t, tc.expectedResp.LongUrl, resp.LongUrl)
		})
	}
}
