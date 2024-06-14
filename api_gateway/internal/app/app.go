package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"api_gateway/internal/converter"
	"api_gateway/internal/transport/rest"
	analytics "api_gateway/pkg/proto"
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

	target := fmt.Sprintf("%s:%s", "analytics_service", "8101")
	transportOpt := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.NewClient(target, transportOpt)
	if err != nil {
		panic(err)
	}
	analyticsClient := analytics.NewAnalyticsClient(conn)

	analyticsHandler := rest.NewAnalyticsHandler(logger, analyticsClient, topUrlConverter, paginationConverter)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/top_urls", analyticsHandler.GetTopURLs)

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
