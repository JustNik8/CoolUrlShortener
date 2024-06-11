package rest

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/service"
	"CoolUrlShortener/internal/service/mocks"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFollowUrl(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	validate := validator.New(validator.WithRequiredStructEnabled())
	testLongURL := "https://test.long"
	testShortURL := "short"
	baseURL := ""
	unexpectedErr := errors.New("unexpected error")

	testCases := []struct {
		name            string
		buildURLService func() service.URLService
		shortURL        string
		expectedCode    int
	}{
		{
			name: "short url is empty. 404 Not found",
			buildURLService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				return mockService
			},
			shortURL:     "",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "redirect by short url. 302 Status found",
			buildURLService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("GetLongURL", mock.Anything, testShortURL).
					Return(testLongURL, nil)

				return mockService
			},
			shortURL:     testShortURL,
			expectedCode: http.StatusFound,
		},
		{
			name: "short url not found. 404 Not found",
			buildURLService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("GetLongURL", mock.Anything, testShortURL).
					Return("", errs.ErrNoURL)

				return mockService
			},
			shortURL:     testShortURL,
			expectedCode: http.StatusNotFound,
		},
		{
			name: "unexpected error. 500 Internal Server Error",
			buildURLService: func() service.URLService {
				mockService := mocks.NewURLService(t)
				mockService.On("GetLongURL", mock.Anything, testShortURL).
					Return("", unexpectedErr)

				return mockService
			},
			shortURL:     testShortURL,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewURLHandler(
				logger,
				tc.buildURLService(),
				validate,
			)

			path := fmt.Sprintf("%s/%s", baseURL, tc.shortURL)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /{short_url}", handler.FollowUrl)

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
