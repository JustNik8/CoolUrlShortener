package rest

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/service"
	"CoolUrlShortener/internal/transport/rest/response"
)

const (
	shortUrlPathValue = "short_url"
)

type URLHandler struct {
	logger     *slog.Logger
	urlService service.URLService
}

func NewURLHandler(
	logger *slog.Logger,
	urlService service.URLService,
) *URLHandler {
	return &URLHandler{
		logger:     logger,
		urlService: urlService,
	}
}

func (h *URLHandler) FollowLink(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue(shortUrlPathValue)

	msg := fmt.Sprintf("got url: %s", shortUrl)
	h.logger.Info(msg)
	if shortUrl == "" {
		msg := "short url is empty"
		h.logger.Info(msg)
		response.BadRequest(w, msg)
		return
	}

	longURL, err := h.urlService.GetLongURL(context.Background(), shortUrl)
	if err != nil {
		h.logger.Error(err.Error())
		if errors.Is(err, errs.ErrNoURL) {
			msg := "short url not found"
			h.logger.Info(msg)
			response.NotFound(w, msg)
			return
		}

		response.InternalServerError(w)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
