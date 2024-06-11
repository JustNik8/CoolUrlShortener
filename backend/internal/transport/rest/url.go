package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"CoolUrlShortener/internal/errs"
	"CoolUrlShortener/internal/service"
	"CoolUrlShortener/internal/transport/rest/dto"
	"CoolUrlShortener/internal/transport/rest/response"
	"github.com/go-playground/validator/v10"
)

const (
	shortUrlPathValue = "short_url"
)

type URLHandler struct {
	logger     *slog.Logger
	urlService service.URLService
	validate   *validator.Validate
}

func NewURLHandler(
	logger *slog.Logger,
	urlService service.URLService,
	validate *validator.Validate,
) *URLHandler {
	return &URLHandler{
		logger:     logger,
		urlService: urlService,
		validate:   validate,
	}
}

func (h *URLHandler) FollowUrl(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.PathValue(shortUrlPathValue)
	msg := fmt.Sprintf("got url: %s", shortUrl)
	h.logger.Info(msg)

	if shortUrl == "" {
		msg := "short url is empty"
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

func (h *URLHandler) SaveURL(w http.ResponseWriter, r *http.Request) {
	var longURLData dto.LongURLData
	err := json.NewDecoder(r.Body).Decode(&longURLData)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	err = h.validate.Struct(longURLData)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	shortURL, err := h.urlService.SaveURL(context.Background(), longURLData.LongURL)
	if err != nil {
		h.logger.Error(err.Error())
		response.InternalServerError(w)
		return
	}

	urlData := dto.URlData{
		LongURL:  longURLData.LongURL,
		ShortURL: shortURL,
	}
	urlBody, err := json.Marshal(urlData)
	if err != nil {
		h.logger.Error(err.Error())
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, urlBody)
}
