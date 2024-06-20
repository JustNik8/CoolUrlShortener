package app

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"analytics_service/internal/config"
	"analytics_service/internal/converter"
	"analytics_service/internal/repository/clickhouserepo"
	"analytics_service/internal/service"
	analytics_grpc "analytics_service/internal/transport/grpc"
	"analytics_service/internal/transport/rest"
	analytics "analytics_service/pkg/proto"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"google.golang.org/grpc"
)

const (
	envLocal = "local"
	envProd  = "prod"

	httpServerPort    = "8002"
	grpcServerPort    = "8102"
	grpcServerNetwork = "tcp"
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
	runGrpcServer(logger, cfg.ClickhouseConfig)
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

func setupClickhouseConn(clickhouseCfg config.ClickhouseConfig) (driver.Conn, error) {

	addr := fmt.Sprintf("%s:%s", clickhouseCfg.Host, clickhouseCfg.Port)

	conn, err := clickhouse.Open(&clickhouse.Options{
		Protocol: clickhouse.Native,
		Addr:     []string{addr},
		Auth: clickhouse.Auth{
			Database: clickhouseCfg.Database,
			Username: clickhouseCfg.Username,
			Password: clickhouseCfg.Password,
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

func runGrpcServer(logger *slog.Logger, clickhouseCfg config.ClickhouseConfig) {
	topURLConverter := converter.NewTopURLConverter()
	paginationConverter := converter.NewPaginationConverter()

	clickhouseConn, err := setupClickhouseConn(clickhouseCfg)
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

	go func() {
		s := grpc.NewServer()
		analyticsServer := analytics_grpc.NewAnalyticsServer(
			logger,
			analyticsService,
			paginationService,
			topURLConverter,
			paginationConverter,
		)

		analytics.RegisterAnalyticsServer(s, analyticsServer)
		port := fmt.Sprintf(":%s", grpcServerPort)
		listener, err := net.Listen(grpcServerNetwork, port)
		if err != nil {
			panic(err)
		}

		err = s.Serve(listener)
		if err != nil {
			logger.Info(err.Error())
		}
	}()
}

func runHttpServer(logger *slog.Logger) {
	healthCheckHandler := rest.NewHealthCheckHandler(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthcheck", healthCheckHandler.HealthCheck)

	go func() {
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
