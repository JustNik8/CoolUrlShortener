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

// HealthCheck docs
//
//	@Summary		Проверка работоспособности сервиса
//	@Tags			healthcheck
//	@Description	Возвращает код 200, когда сервис работоспособен
//	@ID				healthcheck
//	@Success		200
//	@Router			/api/healthcheck [get]
func (h *HealthCheckHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Healthcheck for the service")
	w.WriteHeader(http.StatusOK)
}
