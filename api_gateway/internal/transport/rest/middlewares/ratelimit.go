package middlewares

import (
	"log/slog"
	"net/http"

	"api_gateway/internal/transport/rest/response"
	"golang.org/x/time/rate"
)

type RateLimiterMiddleware struct {
	logger  *slog.Logger
	limiter *rate.Limiter
}

func NewRateLimiterMiddleware(
	logger *slog.Logger,
	limiter *rate.Limiter,
) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		logger:  logger,
		limiter: limiter,
	}
}

func (m *RateLimiterMiddleware) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.limiter.Allow() {
			response.TooManyRequests(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
