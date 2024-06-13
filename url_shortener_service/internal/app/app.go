package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"CoolUrlShortener/internal/repository/events"
	"CoolUrlShortener/internal/repository/postgresql"
	"CoolUrlShortener/internal/repository/rediscache"
	"CoolUrlShortener/internal/service"
	"CoolUrlShortener/internal/transport/rest"
	"CoolUrlShortener/pkg/shortener"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	httpSwagger "github.com/swaggo/http-swagger"
)

const (
	envLocal = "local"
	envProd  = "prod"

	envKey           = "ENV"
	databaseUsername = "DATABASE_USERNAME"
	databasePassword = "DATABASE_PASSWORD"
	databaseHost     = "DATABASE_HOST"
	databasePort     = "DATABASE_PORT"
	databaseName     = "DATABASE_NAME"

	redisHost     = "REDIS_HOST"
	redisPort     = "REDIS_PORT"
	redisPassword = "REDIS_PASSWORD"

	serverDomainKey = "SERVER_DOMAIN"

	serverPort = "8000"
)

func Run() {
	doneCh := make(chan struct{})

	defer func() {
		doneCh <- struct{}{}
		close(doneCh)
	}()

	serverDomain := os.Getenv(serverDomainKey)
	if serverDomain == "" {
		msg := fmt.Sprintf("You did not provide env: %s", serverDomainKey)
		panic(msg)
	}

	env := os.Getenv(envKey)
	if env == "" {
		msg := fmt.Sprintf("You did not provide env: %s", envKey)
		panic(msg)
	}

	logger, err := setupLogger(env)
	if err != nil {
		panic(err)
	}

	connString, err := setupConnString()
	if err != nil {
		panic(err)
	}
	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	validate := validator.New(validator.WithRequiredStructEnabled())

	eventsServiceProducer, err := events.NewKafkaEventProducer(logger, []string{"kafka1:9092"}, nil, doneCh)
	if err != nil {
		panic(err)
	}

	redisClient, err := setupRedisClient()
	if err != nil {
		panic(err.Error())
	}

	base62URLShortener := shortener.NewBase62UrlShortener()

	urlCache := rediscache.NewURLCacheRedis(redisClient)
	urlRepo := postgresql.NewUrlRepoPostgres(dbPool)
	urlService := service.NewURLService(logger, urlRepo, urlCache, eventsServiceProducer, base62URLShortener)
	urlHandler := rest.NewURLHandler(logger, urlService, validate, serverDomain)

	healthCheckHandler := rest.NewHealthCheckHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/healthcheck", healthCheckHandler.HealthCheck)
	mux.HandleFunc("POST /api/save_url", urlHandler.SaveURL)
	mux.HandleFunc("OPTIONS /api/save_url", urlHandler.SaveURLOptions)
	mux.HandleFunc("GET /{short_url}", urlHandler.FollowUrl)

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

func setupConnString() (string, error) {
	username := os.Getenv(databaseUsername)
	if username == "" {
		return "", fmt.Errorf("you did not provide env: %s", databaseUsername)
	}

	password := os.Getenv(databasePassword)
	if password == "" {
		return "", fmt.Errorf("you did not provide env: %s", databasePassword)
	}

	host := os.Getenv(databaseHost)
	if host == "" {
		return "", fmt.Errorf("you did not provide env: %s", databaseHost)
	}

	port := os.Getenv(databasePort)
	if port == "" {
		return "", fmt.Errorf("you did not provide env: %s", databasePort)
	}

	dbName := os.Getenv(databaseName)
	if dbName == "" {
		return "", fmt.Errorf("you did not provide env: %s", databaseName)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)
	return connString, nil
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

func setupRedisClient() (*redis.Client, error) {
	host := os.Getenv(redisHost)
	if host == "" {
		return nil, fmt.Errorf("you did not provide env: %s", redisHost)
	}

	port := os.Getenv(redisPort)
	if port == "" {
		return nil, fmt.Errorf("you did not provide env: %s", redisPort)
	}

	password := os.Getenv(redisPassword)
	if password == "" {
		return nil, fmt.Errorf("you did not provide env: %s", redisPassword)
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return redisClient, nil
}
