package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"api_gateway/internal/config"
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

	httpServerPort = "8000"
)

func Run() {
	cfg, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	logger, err := setupLogger(cfg.Env)
	if err != nil {
		panic(err)
	}

	topUrlConverter := converter.NewTopURLConverter()
	paginationConverter := converter.NewPaginationConverter()

	urlTarget := fmt.Sprintf("%s:%s", cfg.UrlServiceConfig.Host, cfg.UrlServiceConfig.Port)
	urlTransportOpr := grpc.WithTransportCredentials(insecure.NewCredentials())

	urlConn, err := grpc.NewClient(urlTarget, urlTransportOpr)
	if err != nil {
		panic(err)
	}
	urlClient := url.NewUrlClient(urlConn)

	analyticsTarget := fmt.Sprintf("%s:%s", cfg.AnalyticsServiceConfig.Host, cfg.AnalyticsServiceConfig.Port)
	analyticsTransportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	analyticsConn, err := grpc.NewClient(analyticsTarget, analyticsTransportOpt)
	if err != nil {
		panic(err)
	}
	analyticsClient := analytics.NewAnalyticsClient(analyticsConn)

	urlHandler := rest.NewURLHandler(logger, urlClient, validator.New(), cfg.ServerDomain, httpServerPort)

	analyticsHandler := rest.NewAnalyticsHandler(logger, analyticsClient, topUrlConverter, paginationConverter)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/top_urls", analyticsHandler.GetTopURLs)
	mux.HandleFunc("POST /api/save_url", urlHandler.SaveURL)
	mux.HandleFunc("OPTIONS /api/save_url", urlHandler.SaveURLOptions)
	mux.HandleFunc("GET /{short_url}", urlHandler.FollowUrl)

	addr := fmt.Sprintf(":%s", httpServerPort)
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
