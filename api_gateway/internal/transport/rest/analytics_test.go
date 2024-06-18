package rest

//import (
//	"fmt"
//	"log/slog"
//	"net/http"
//	"net/http/httptest"
//	"os"
//	"testing"
//
//	"api_gateway/internal/client"
//)
//
//func TestGetTopURLs(t *testing.T) {
//	logger := slog.New(
//		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
//	)
//	serverDomain := "test"
//	serverPort := "8000"
//	basePath := "/api/top_urls"
//
//	testCases := []struct {
//		name           string
//		buildUrlClient func() client.UrlClient
//		shortURL       string
//		expectedCode   int
//	}{
//		{},
//	}
//
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			handler := NewAnalyticsHandler(
//				logger,
//				tc.buildUrlClient(),
//				serverDomain,
//				serverPort,
//			)
//
//			path := fmt.Sprintf("%s/%s", basePath, tc.shortURL)
//			req := httptest.NewRequest(http.MethodGet, path, nil)
//			rec := httptest.NewRecorder()
//
//			handler.GetTopURLs(rec, req)
//
//			assert.Equal(t, tc.expectedCode, rec.Code)
//		})
//	}
//}
