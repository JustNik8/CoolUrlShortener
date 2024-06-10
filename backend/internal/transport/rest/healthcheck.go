package rest

import (
	"log/slog"
	"net/http"
)

type HealthCheckHandler struct {
	logger *slog.Logger
}

func NewHealthCheckHandler(logger *slog.Logger) *HealthCheckHandler {
	return &HealthCheckHandler{
		logger: logger,
	}
}

func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Healthcheck for the service")
	w.WriteHeader(http.StatusOK)
}
