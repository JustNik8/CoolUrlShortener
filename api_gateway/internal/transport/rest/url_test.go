package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"api_gateway/errs"
	"api_gateway/internal/client"
	"api_gateway/internal/client/mocks"
	"api_gateway/internal/transport/rest/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFollowUrl(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	serverDomain := "test"
	basePath := ""

	testErr := errors.New("test error")

	testCases := []struct {
		name           string
		buildUrlClient func() client.UrlClient
		shortURL       string
		expectedCode   int
	}{
		{
			name: "redirect by short url. 302 Status found",
			buildUrlClient: func() client.UrlClient {
				mockClient := mocks.NewUrlClient(t)
				mockClient.On("FollowUrl", mock.Anything, mock.Anything).
					Return("http://test.long", nil)

				return mockClient
			},
			shortURL:     "short",
			expectedCode: http.StatusFound,
		},
		{
			name: "short url is empty. 404 Not found",
			buildUrlClient: func() client.UrlClient {
				mockClient := mocks.NewUrlClient(t)
				return mockClient
			},
			shortURL:     "",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "short url not found. 404 Not found",
			buildUrlClient: func() client.UrlClient {
				mockClient := mocks.NewUrlClient(t)
				mockClient.On("FollowUrl", mock.Anything, mock.Anything).
					Return("", errs.ErrNotFound)

				return mockClient
			},
			shortURL:     "test",
			expectedCode: http.StatusNotFound,
		},
		{
			name: "unexpected error. 500 Internal Server Error",
			buildUrlClient: func() client.UrlClient {
				mockClient := mocks.NewUrlClient(t)
				mockClient.On("FollowUrl", mock.Anything, mock.Anything).
					Return("", testErr)

				return mockClient
			},
			shortURL:     "test",
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewURLHandler(
				logger,
				tc.buildUrlClient(),
				serverDomain,
			)

			path := fmt.Sprintf("%s/%s", basePath, tc.shortURL)
			req := httptest.NewRequest(http.MethodGet, path, nil)
			rec := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("GET /{short_url}", handler.FollowUrl)

			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestSaveURL(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	basePath := "/api/save_url"

	testErr := errors.New("test error")
	serverDomain := "test:8000"

	testCases := []struct {
		name             string
		buildUrlClient   func() client.UrlClient
		longUrlRequest   dto.LongURLData
		expectedCode     int
		expectedLongURL  string
		expectedShortURL string
	}{
		{
			name: "Empty long url. 400 Bad Request",
			buildUrlClient: func() client.UrlClient {
				mockClient := mocks.NewUrlClient(t)
				mockClient.On("ShortenUrl", mock.Anything, mock.Anything).
					Return("", errs.ErrInvalidArgument)

				return mockClient
			},
			longUrlRequest: dto.LongURLData{
				LongURL: "",
			},
			expectedCode:     http.StatusBadRequest,
			expectedLongURL:  "",
			expectedShortURL: "",
		},
		{
			name: "Create short url without error. 200 Status OK",
			buildUrlClient: func() client.UrlClient {
				mockClient := mocks.NewUrlClient(t)
				mockClient.On("ShortenUrl", mock.Anything, mock.Anything).
					Return("short", nil)

				return mockClient
			},
			longUrlRequest: dto.LongURLData{
				LongURL: "http://test.long",
			},
			expectedCode:     http.StatusOK,
			expectedLongURL:  "http://test.long",
			expectedShortURL: fmt.Sprintf("%s://%s/%s", serverProtocol, serverDomain, "short"),
		},
		{
			name: "Unexpected error while saving url. 500 Internal Server Error",
			buildUrlClient: func() client.UrlClient {
				mockClient := mocks.NewUrlClient(t)
				mockClient.On("ShortenUrl", mock.Anything, mock.Anything).
					Return("", testErr)

				return mockClient
			},
			longUrlRequest: dto.LongURLData{
				LongURL: "http://test.long",
			},
			expectedCode:     http.StatusInternalServerError,
			expectedLongURL:  "",
			expectedShortURL: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewURLHandler(
				logger,
				tc.buildUrlClient(),
				serverDomain,
			)

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(tc.longUrlRequest)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, basePath, &buf)
			rec := httptest.NewRecorder()

			handler.SaveURL(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)

			if rec.Code == http.StatusOK {
				urlData := dto.URlData{}
				err = json.NewDecoder(rec.Body).Decode(&urlData)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedLongURL, urlData.LongURL)
				assert.Equal(t, tc.expectedShortURL, urlData.ShortURL)
			}
		})
	}
}

func FuzzSaveURL(f *testing.F) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	basePath := "/api/save_url"

	testShortURL := "http://test/short"
	serverDomain := "test"

	mockClient := mocks.NewUrlClient(f)

	handler := NewURLHandler(
		logger,
		mockClient,
		serverDomain,
	)

	args := []dto.LongURLData{
		{LongURL: "https://test.long1"},
		{LongURL: "https://test.long2"},
		{LongURL: "https://test.long3"},
	}

	for _, arg := range args {
		data, _ := json.Marshal(arg)
		f.Add(data)
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		mockClient.On("ShortenUrl", mock.Anything, mock.AnythingOfType("string")).
			Return(testShortURL, nil)

		req := httptest.NewRequest(http.MethodPost, basePath, bytes.NewBuffer(data))
		rec := httptest.NewRecorder()

		handler.SaveURL(rec, req)

		var longUrlData dto.LongURLData
		err := json.NewDecoder(bytes.NewReader(data)).Decode(&longUrlData)

		if err != nil {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			return
		}

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
