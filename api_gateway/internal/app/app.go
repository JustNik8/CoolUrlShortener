package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"api_gateway/internal/converter"
	"api_gateway/internal/transport/rest"
	"api_gateway/pkg/proto/analytics"
	"api_gateway/pkg/proto/url"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	envLocal = "local"
	envProd  = "prod"

	envKey = "ENV"
)

func Run() {
	env := os.Getenv(envKey)
	if env == "" {
		msg := fmt.Sprintf("You did not provide env: %s", envKey)
		panic(msg)
	}

	logger, err := setupLogger(env)
	if err != nil {
		panic(err)
	}

	topUrlConverter := converter.NewTopURLConverter()
	paginationConverter := converter.NewPaginationConverter()

	analyticsTarget := fmt.Sprintf("%s:%s", "analytics_service", "8101")
	analyticsTransportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	analyticsConn, err := grpc.NewClient(analyticsTarget, analyticsTransportOpt)
	if err != nil {
		panic(err)
	}
	analyticsClient := analytics.NewAnalyticsClient(analyticsConn)

	urlTarget := fmt.Sprintf("%s:%s", "url_shortener_service", "8001")
	urlTransportOpr := grpc.WithTransportCredentials(insecure.NewCredentials())

	urlConn, err := grpc.NewClient(urlTarget, urlTransportOpr)
	if err != nil {
		panic(err)
	}
	urlClient := url.NewUrlClient(urlConn)

	analyticsHandler := rest.NewAnalyticsHandler(logger, analyticsClient, topUrlConverter, paginationConverter)
	urlHandler := rest.NewURLHandler(logger, urlClient, validator.New(), "localhost:8200")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/top_urls", analyticsHandler.GetTopURLs)
	mux.HandleFunc("POST /api/save_url", urlHandler.SaveURL)
	mux.HandleFunc("OPTIONS /api/save_url", urlHandler.SaveURLOptions)
	mux.HandleFunc("GET /{short_url}", urlHandler.FollowUrl)

	addr := fmt.Sprintf(":%s", "8200")
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	logger.Info(fmt.Sprintf("Run server on %s", addr))
	err = server.ListenAndServe()
	if err != nil {
		logger.Info(err.Error())
	}
}

func setupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return nil, fmt.Errorf("incorrect env %s", env)
	}

	return logger, nil
}
