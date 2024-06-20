package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"CoolUrlShortener/internal/config"
	"CoolUrlShortener/internal/repository/events"
	"CoolUrlShortener/internal/repository/postgresql"
	"CoolUrlShortener/internal/repository/rediscache"
	"CoolUrlShortener/internal/service"
	url_grpc "CoolUrlShortener/internal/transport/grpc"
	"CoolUrlShortener/internal/transport/rest"
	url "CoolUrlShortener/pkg/proto"
	"CoolUrlShortener/pkg/shortener"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

const (
	envLocal = "local"
	envProd  = "prod"

	httpServerPort    = "8001"
	grpcServerPort    = "8101"
	grpcServerNetwork = "tcp"
)

func Run() {
	doneCh := make(chan struct{})

	defer func() {
		doneCh <- struct{}{}
		close(doneCh)
	}()

	cfg, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	logger, err := setupLogger(cfg.Env)
	if err != nil {
		panic(err)
	}

	dbPool := createDBPool(cfg.DatabaseConfig)
	defer dbPool.Close()

	runGrpcServer(logger, cfg, dbPool, doneCh)
	runHttpServer(logger)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
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

func createDBPool(dbCfg config.DatabaseConfig) *pgxpool.Pool {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
	)

	dbPool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		panic(err)
	}

	return dbPool
}

func setupRedisClient(redisCfg config.RedisConfig) (*redis.Client, error) {

	addr := fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisCfg.Password,
		DB:       0,
	})

	return redisClient, nil
}

func runGrpcServer(
	logger *slog.Logger,
	cfg config.Config,
	dbPool *pgxpool.Pool,
	doneCh <-chan struct{},
) {

	eventsServiceProducer, err := events.NewKafkaEventProducer(logger, cfg.KafkaConfig.Addrs, nil, doneCh)
	if err != nil {
		panic(err)
	}

	redisClient, err := setupRedisClient(cfg.RedisConfig)
	if err != nil {
		panic(err.Error())
	}

	base62URLShortener := shortener.NewBase62UrlShortener()

	urlCache := rediscache.NewURLCacheRedis(redisClient)
	urlRepo := postgresql.NewUrlRepoPostgres(dbPool)
	urlService := service.NewURLService(logger, urlRepo, urlCache, eventsServiceProducer, base62URLShortener)

	go func() {
		s := grpc.NewServer()
		urlServer := url_grpc.NewUrlServer(
			logger,
			urlService,
		)

		url.RegisterUrlServer(s, urlServer)
		port := fmt.Sprintf(":%s", grpcServerPort)
		listener, err := net.Listen(grpcServerNetwork, port)
		if err != nil {
			panic(err)
		}

		err = s.Serve(listener)
		if err != nil {
			logger.Info(err.Error())
			return
		}
	}()
}

func runHttpServer(logger *slog.Logger) {
	healthCheckHandler := rest.NewHealthCheckHandler(logger)

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("GET /api/healthcheck", healthCheckHandler.HealthCheck)

		addr := fmt.Sprintf(":%s", httpServerPort)
		server := http.Server{
			Addr:    addr,
			Handler: mux,
		}

		logger.Info(fmt.Sprintf("Run server on %s", addr))
		err := server.ListenAndServe()
		if err != nil {
			logger.Info(err.Error())
		}
	}()
}
