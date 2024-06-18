package rest

import (
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//
//import (
//	"bytes"
//	"encoding/json"
//	"errors"
//	"fmt"
//	"log/slog"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"testing"
//
//	"CoolUrlShortener/internal/errs"
//	"CoolUrlShortener/internal/service"
//	"CoolUrlShortener/internal/service/mocks"
//	"CoolUrlShortener/internal/transport/rest/dto"
//	"github.com/go-playground/validator/v10"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//)

func TestFollowUrl(t *testing.T) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)
	serverDomain := "test"
	serverPort := "8000"
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
				serverPort,
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

//func TestFollowUrl(t *testing.T) {
//	logger := slog.New(
//		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//	)
//	validate := validator.New(validator.WithRequiredStructEnabled())
//	testLongURL := "https://test.long"
//	testShortURL := "short"
//	basePath := ""
//	unexpectedErr := errors.New("unexpected error")
//	serverDomain := "test"
//
//	testCases := []struct {
//		name            string
//		buildURLService func() service.URLService
//		shortURL        string
//		expectedCode    int
//	}{
//		{
//			name: "short url is empty. 404 Not found",
//			buildURLService: func() service.URLService {
//				mockService := mocks.NewURLService(t)
//				return mockService
//			},
//			shortURL:     "",
//			expectedCode: http.StatusNotFound,
//		},
//		{
//			name: "redirect by short url. 302 Status found",
//			buildURLService: func() service.URLService {
//				mockService := mocks.NewURLService(t)
//				mockService.On("GetLongURL", mock.Anything, testShortURL).
//					Return(testLongURL, nil)
//
//				return mockService
//			},
//			shortURL:     testShortURL,
//			expectedCode: http.StatusFound,
//		},
//		{
//			name: "short url not found. 404 Not found",
//			buildURLService: func() service.URLService {
//				mockService := mocks.NewURLService(t)
//				mockService.On("GetLongURL", mock.Anything, testShortURL).
//					Return("", errs.ErrNoURL)
//
//				return mockService
//			},
//			shortURL:     testShortURL,
//			expectedCode: http.StatusNotFound,
//		},
//		{
//			name: "unexpected error. 500 Internal Server Error",
//			buildURLService: func() service.URLService {
//				mockService := mocks.NewURLService(t)
//				mockService.On("GetLongURL", mock.Anything, testShortURL).
//					Return("", unexpectedErr)
//
//				return mockService
//			},
//			shortURL:     testShortURL,
//			expectedCode: http.StatusInternalServerError,
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			handler := NewURLHandler(
//				logger,
//				tc.buildURLService(),
//				validate,
//				serverDomain,
//			)
//
//			path := fmt.Sprintf("%s/%s", basePath, tc.shortURL)
//			req := httptest.NewRequest(http.MethodGet, path, nil)
//			rec := httptest.NewRecorder()
//
//			mux := http.NewServeMux()
//			mux.HandleFunc("GET /{short_url}", handler.FollowUrl)
//
//			mux.ServeHTTP(rec, req)
//
//			assert.Equal(t, tc.expectedCode, rec.Code)
//		})
//	}
//}
//
//func TestSaveURL(t *testing.T) {
//	logger := slog.New(
//		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//	)
//	validate := validator.New(validator.WithRequiredStructEnabled())
//	basePath := "/api/save_url"
//
//	testLongURL := "https://test.long"
//	testShortURL := "short"
//	unexpectedErr := errors.New("unexpected error")
//
//	//serverProtocol := "http"
//	serverDomain := "test"
//
//	testCases := []struct {
//		name             string
//		buildURLService  func() service.URLService
//		buildLongURLData dto.LongURLData
//		expectedCode     int
//		expectedLongURL  string
//		expectedShortURL string
//	}{
//		{
//			name: "Empty long url. 400 Bad Request",
//			buildURLService: func() service.URLService {
//				mockService := mocks.NewURLService(t)
//				return mockService
//			},
//			buildLongURLData: dto.LongURLData{
//				LongURL: "",
//			},
//			expectedCode:     http.StatusBadRequest,
//			expectedLongURL:  "",
//			expectedShortURL: "",
//		},
//		{
//			name: "Create short url without error. 200 Status OK",
//			buildURLService: func() service.URLService {
//				mockService := mocks.NewURLService(t)
//				mockService.On("SaveURL", mock.Anything, testLongURL).
//					Return(testShortURL, nil)
//
//				return mockService
//			},
//			buildLongURLData: dto.LongURLData{
//				LongURL: testLongURL,
//			},
//			expectedCode:     http.StatusOK,
//			expectedLongURL:  testLongURL,
//			expectedShortURL: fmt.Sprintf("%s://%s/%s", serverProtocol, serverDomain, testShortURL),
//		},
//		{
//			name: "Unexpected error while saving url. 500 Internal Server Error",
//			buildURLService: func() service.URLService {
//				mockService := mocks.NewURLService(t)
//				mockService.On("SaveURL", mock.Anything, testLongURL).
//					Return("", unexpectedErr)
//
//				return mockService
//			},
//			buildLongURLData: dto.LongURLData{
//				LongURL: testLongURL,
//			},
//			expectedCode:     http.StatusInternalServerError,
//			expectedLongURL:  "",
//			expectedShortURL: "",
//		},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			handler := NewURLHandler(
//				logger,
//				tc.buildURLService(),
//				validate,
//				serverDomain,
//			)
//
//			var buf bytes.Buffer
//			err := json.NewEncoder(&buf).Encode(tc.buildLongURLData)
//			assert.NoError(t, err)
//
//			req := httptest.NewRequest(http.MethodPost, basePath, &buf)
//			rec := httptest.NewRecorder()
//
//			handler.SaveURL(rec, req)
//			assert.Equal(t, tc.expectedCode, rec.Code)
//
//			if rec.Code == http.StatusOK {
//				urlData := dto.URlData{}
//				err = json.NewDecoder(rec.Body).Decode(&urlData)
//				assert.NoError(t, err)
//
//				assert.Equal(t, tc.expectedLongURL, urlData.LongURL)
//				assert.Equal(t, tc.expectedShortURL, urlData.ShortURL)
//			}
//		})
//	}
//}
//
//func FuzzSaveURL(f *testing.F) {
//	logger := slog.New(
//		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//	)
//	validate := validator.New(validator.WithRequiredStructEnabled())
//	basePath := "/api/save_url"
//
//	testShortURL := "http://test/short"
//	serverDomain := "test"
//
//	mockService := mocks.NewURLService(f)
//
//	handler := NewURLHandler(
//		logger,
//		mockService,
//		validate,
//		serverDomain,
//	)
//
//	args := []dto.LongURLData{
//		{LongURL: "https://test.long1"},
//		{LongURL: "https://test.long2"},
//		{LongURL: "https://test.long3"},
//	}
//
//	for _, arg := range args {
//		data, _ := json.Marshal(arg)
//		f.Add(data)
//	}
//
//	f.Fuzz(func(t *testing.T, data []byte) {
//		mockService.On("SaveURL", mock.Anything, mock.AnythingOfType("string")).
//			Return(testShortURL, nil)
//
//		req := httptest.NewRequest(http.MethodPost, basePath, bytes.NewBuffer(data))
//		rec := httptest.NewRecorder()
//
//		handler.SaveURL(rec, req)
//
//		var longUrlData dto.LongURLData
//		err := json.NewDecoder(bytes.NewReader(data)).Decode(&longUrlData)
//
//		if err != nil {
//			assert.Equal(t, http.StatusBadRequest, rec.Code)
//			return
//		}
//
//		err = validate.Struct(longUrlData)
//
//		if err != nil {
//			assert.Equal(t, http.StatusBadRequest, rec.Code)
//			return
//		}
//
//		assert.Equal(t, http.StatusOK, rec.Code)
//	})
//}
