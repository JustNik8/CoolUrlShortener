package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"CoolUrlShortener/internal/domain"
	"CoolUrlShortener/internal/repository"
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
	serverPortKey    = "SERVER_PORT"
	databaseUsername = "DATABASE_USERNAME"
	databasePassword = "DATABASE_PASSWORD"
	databaseHost     = "DATABASE_HOST"
	databasePort     = "DATABASE_PORT"
	databaseName     = "DATABASE_NAME"

	redisHost     = "REDIS_HOST"
	redisPort     = "REDIS_PORT"
	redisPassword = "REDIS_PASSWORD"

	clickhouseUsername = "CLICKHOUSE_USERNAME"
	clickhousePassword = "CLICKHOUSE_PASSWORD"
	clickhouseHost     = "CLICKHOUSE_HOST"
	clickhousePort     = "CLICKHOUSE_PORT"
	clickhouseDatabase = "CLICKHOUSE_DATABASE"
	batchTimePeriod    = "BATCH_TIME_PERIOD"
)

func Run() {
	eventsCh := make(chan domain.URLEvent)
	doneCh := make(chan struct{})

	defer func() {
		close(eventsCh)
		doneCh <- struct{}{}
		close(doneCh)
	}()

	serverPort := os.Getenv(serverPortKey)
	if serverPort == "" {
		msg := fmt.Sprintf("You did not provide env: %s", serverPortKey)
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

	eventsWriter, err := setupEventsWriter(logger)
	if err != nil {
		panic(err)
	}

	eventsServiceConsumer, err := setupEventsServiceConsumer(logger, eventsCh, doneCh, eventsWriter)
	if err != nil {
		panic(err)
	}
	eventsServiceProducer := service.NewEventsServiceProducer(eventsCh)

	redisClient, err := setupRedisClient()
	if err != nil {
		panic(err.Error())
	}

	base62URLShortener := shortener.NewBase62UrlShortener()

	urlCache := rediscache.NewURLCacheRedis(redisClient)
	urlRepo := postgresql.NewUrlRepoPostgres(dbPool)
	urlService := service.NewURLService(logger, urlRepo, urlCache, eventsServiceProducer, base62URLShortener)
	urlHandler := rest.NewURLHandler(logger, urlService, validate)

	healthCheckHandler := rest.NewHealthCheckHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/healthcheck", healthCheckHandler.HealthCheck)
	mux.HandleFunc("POST /api/save_url", urlHandler.SaveURL)
	mux.HandleFunc("GET /{short_url}", urlHandler.FollowUrl)

	addr := fmt.Sprintf(":%s", serverPort)
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	logger.Info("Start consume url events")
	eventsServiceConsumer.ConsumeEvents()

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

func setupEventsWriter(logger *slog.Logger) (repository.EventsWriter, error) {
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
	eventsWriter, err := events.NewEventsWriterClickhouse(
		logger,
		database,
		username,
		password,
		host,
		port,
	)

	return eventsWriter, err
}

func setupEventsServiceConsumer(
	logger *slog.Logger,
	eventsCh <-chan domain.URLEvent,
	doneCh <-chan struct{},
	eventsWriter repository.EventsWriter,
) (service.EventsServiceConsumer, error) {
	periodTime := os.Getenv(batchTimePeriod)
	if periodTime == "" {
		return nil, fmt.Errorf("you did not provice env: %s", batchTimePeriod)
	}

	duration, err := time.ParseDuration(periodTime)
	if err != nil {
		return nil, err
	}
	periodCh := time.NewTicker(duration).C

	eventsServiceConsumer := service.NewEventsServiceConsumer(logger, eventsCh, periodCh, doneCh, eventsWriter)
	return eventsServiceConsumer, nil
}
