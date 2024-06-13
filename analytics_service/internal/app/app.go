package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"analytics_service/internal/converter"
	"analytics_service/internal/repository"
	"analytics_service/internal/repository/analytics"
	"analytics_service/internal/service"
	"analytics_service/internal/transport/rest"
)

const (
	envLocal = "local"
	envProd  = "prod"

	envKey = "ENV"

	clickhouseUsername = "CLICKHOUSE_USERNAME"
	clickhousePassword = "CLICKHOUSE_PASSWORD"
	clickhouseHost     = "CLICKHOUSE_HOST"
	clickhousePort     = "CLICKHOUSE_PORT"
	clickhouseDatabase = "CLICKHOUSE_DATABASE"

	serverPort = "8001"
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

	topURLConverter := converter.NewTopURLConverter()

	analyticsRepo, err := setupAnalyticsRepo(logger)
	if err != nil {
		panic(err)
	}
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	analyticsHandler := rest.NewAnalyticsHandler(logger, analyticsService, topURLConverter)

	healthCheckHandler := rest.NewHealthCheckHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/top_urls", analyticsHandler.GetTopURLs)
	mux.HandleFunc("GET /api/healthcheck", healthCheckHandler.HealthCheck)

	addr := fmt.Sprintf(":%s", serverPort)

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

func setupAnalyticsRepo(
	logger *slog.Logger,
) (repository.AnalyticsRepo, error) {
	username := os.Getenv(clickhouseUsername)
	if username == "" {
		return nil, fmt.Errorf("you did not provide env: %s", clickhouseUsername)
	}
	password := os.Getenv(clickhousePassword)
	if password == "" {
		return nil, fmt.Errorf("you did not provice env: %s", clickhousePassword)
	}
	host := os.Getenv(clickhouseHost)
	if host == "" {
		return nil, fmt.Errorf("you did not provide env: %s", clickhouseHost)
	}
	port := os.Getenv(clickhousePort)
	if port == "" {
		return nil, fmt.Errorf("you did not provide env: %s", clickhousePort)
	}

	database := os.Getenv(clickhouseDatabase)
	if database == "" {
		return nil, fmt.Errorf("you did not provice env: %s", clickhouseDatabase)
	}

	analyticsRepo, err := analytics.NewAnalyticsRepoClickhouse(
		logger,
		database,
		username,
		password,
		host,
		port,
	)
	return analyticsRepo, err
}
