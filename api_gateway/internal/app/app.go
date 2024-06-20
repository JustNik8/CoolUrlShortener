package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"api_gateway/internal/client"
	"api_gateway/internal/config"
	"api_gateway/internal/converter"
	"api_gateway/internal/transport/rest"
	"api_gateway/internal/transport/rest/middlewares"
	"api_gateway/pkg/proto/analytics"
	"api_gateway/pkg/proto/url"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/time/rate"
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

	runHttpServer(logger, cfg)
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

func runHttpServer(logger *slog.Logger, cfg config.Config) {
	topUrlConverter := converter.NewTopURLConverter()
	paginationConverter := converter.NewPaginationConverter()

	urlTarget := fmt.Sprintf("%s:%s", cfg.UrlServiceConfig.Host, cfg.UrlServiceConfig.Port)
	urlTransportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	urlConn, err := grpc.NewClient(urlTarget, urlTransportOpt)
	if err != nil {
		panic(err)
	}
	grpcUrlClient := url.NewUrlClient(urlConn)

	analyticsTarget := fmt.Sprintf("%s:%s", cfg.AnalyticsServiceConfig.Host, cfg.AnalyticsServiceConfig.Port)
	analyticsTransportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	analyticsConn, err := grpc.NewClient(analyticsTarget, analyticsTransportOpt)
	if err != nil {
		panic(err)
	}

	analyticsGrpcClient := analytics.NewAnalyticsClient(analyticsConn)
	analyticsClient := client.NewGrpcAnalyticsClient(logger, analyticsGrpcClient, topUrlConverter, paginationConverter)

	limiter := rate.NewLimiter(rate.Limit(cfg.RateLimitConfig.TokensPerSecond), cfg.RateLimitConfig.BurstSize)
	rateLimitMiddleware := middlewares.NewRateLimiterMiddleware(
		logger, limiter,
	)

	urlClient := client.NewGrpcUrlClient(logger, grpcUrlClient)
	urlHandler := rest.NewURLHandler(logger, urlClient, cfg.ServerDomain)
	analyticsHandler := rest.NewAnalyticsHandler(logger, analyticsClient)

	mux := http.NewServeMux()
	mux.Handle("GET /api/top_urls", rateLimitMiddleware.RateLimit(
		http.HandlerFunc(analyticsHandler.GetTopURLs),
	))
	mux.Handle("POST /api/save_url", rateLimitMiddleware.RateLimit(
		http.HandlerFunc(urlHandler.SaveURL),
	))
	mux.HandleFunc("OPTIONS /api/save_url", urlHandler.SaveURLOptions)
	mux.Handle("GET /{short_url}", rateLimitMiddleware.RateLimit(
		http.HandlerFunc(urlHandler.FollowUrl),
	))
	mux.Handle("GET /api/docs/", httpSwagger.WrapHandler)

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
