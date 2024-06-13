package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"analytics_service/internal/converter"
	"analytics_service/internal/repository/clickhouserepo"
	"analytics_service/internal/service"
	"analytics_service/internal/transport/rest"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
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
	paginationConverter := converter.NewPaginationConverter()

	clickhouseConn, err := setupClickhouseConn()
	if err != nil {
		panic(err)
	}

	paginationRepo := clickhouserepo.NewPaginationRepoClickhouse(clickhouseConn)
	paginationService := service.NewPaginationService(paginationRepo)

	analyticsRepo, err := clickhouserepo.NewAnalyticsRepoClickhouse(logger, clickhouseConn)
	if err != nil {
		panic(err)
	}
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	analyticsHandler := rest.NewAnalyticsHandler(
		logger, analyticsService, paginationService, topURLConverter, paginationConverter,
	)

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

func setupClickhouseConn() (driver.Conn, error) {
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

	addr := fmt.Sprintf("%s:%s", host, port)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Protocol: clickhouse.Native,
		Addr:     []string{addr},
		Auth: clickhouse.Auth{
			Database: database,
			Username: username,
			Password: password,
		},
		Debug:           true,
		DialTimeout:     30 * time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})

	if err != nil {
		return nil, err
	}
	return conn, err
}
