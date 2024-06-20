package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"api_gateway/errs"
	"api_gateway/internal/client"
	"api_gateway/internal/transport/rest/dto"
	"api_gateway/internal/transport/rest/response"
)

const (
	shortUrlPathValue = "short_url"
	serverProtocol    = "http"
)

type URLHandler struct {
	logger       *slog.Logger
	urlClient    client.UrlClient
	serverDomain string
}

func NewURLHandler(
	logger *slog.Logger,
	urlClient client.UrlClient,
	serverDomain string,
) *URLHandler {
	return &URLHandler{
		logger:       logger,
		urlClient:    urlClient,
		serverDomain: serverDomain,
	}
}

// FollowUrl docs
//
//	@Summary		Редирект с короткой ссылки на исходную ссылку
//	@Tags			url
//	@Description	Принимает короткую ссылку в path параметрах и производит редирект на исходную ссылку
//	@ID				follow-url
//	@Param			id	query	string	true	"короткая ссылка"
//	@Success		302
//	@Failure		400,404	{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/{short_url} [get]
func (h *URLHandler) FollowUrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	shortUrl := r.PathValue(shortUrlPathValue)

	longUrl, err := h.urlClient.FollowUrl(context.Background(), shortUrl)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			response.NotFound(w, "short url not found")
			return
		}
		if errors.Is(err, errs.ErrInvalidArgument) {
			response.BadRequest(w, "bad short url")
			return
		}

		response.InternalServerError(w)
		return
	}

	http.Redirect(w, r, longUrl, http.StatusFound)
}

// SaveURL docs
//
//	@Summary		Создание и сохранение короткой ссылки по исходной ссылки
//	@Tags			url
//	@Description	Принимает исходную ссылку, создает короткую ссылку и возвращает короткую ссылку
//	@ID				save-url
//	@Accept			json
//	@Produce		json
//	@Param			input	body		dto.LongURLData	true	"Длинная ссылка"
//	@Success		200		{object}	dto.URlData
//	@Failure		400		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/api/save_url [post]
func (h *URLHandler) SaveURL(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}

	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	var longURLData dto.LongURLData
	err := json.NewDecoder(r.Body).Decode(&longURLData)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	shortURLRaw, err := h.urlClient.ShortenUrl(context.Background(), longURLData.LongURL)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidArgument) {
			response.BadRequest(w, err.Error())
			return
		}
		response.InternalServerError(w)
		return
	}

	shortURL := fmt.Sprintf("%s://%s/%s", serverProtocol, h.serverDomain, shortURLRaw)
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

// SaveURLOptions docs
//
//	@Summary		Получение описания параметров соединения с сервером
//	@Tags			options
//	@Description	Возвращает информацию по хедерам Access-Control-Request-Method, Access-Control-Request-Headers, Origin
//	@ID				options-save-url
//	@Success		200	""
//	@Router			/api/save_url [options]
func (h *URLHandler) SaveURLOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Request-Method", "POST")
	w.Header().Add("Access-Control-Request-Headers", "x-requested-with")
	w.Header().Add("Origin", "*")
}
