package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"api_gateway/internal/transport/rest/dto"
	"api_gateway/internal/transport/rest/response"
	"api_gateway/pkg/proto/url"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	shortUrlPathValue = "short_url"
	serverProtocol    = "http"
)

type URLHandler struct {
	logger        *slog.Logger
	urlGrpcClient url.UrlClient
	validate      *validator.Validate
	serverDomain  string
}

func NewURLHandler(
	logger *slog.Logger,
	urlGrpcClient url.UrlClient,
	validate *validator.Validate,
	serverDomain string,
) *URLHandler {
	return &URLHandler{
		logger:        logger,
		urlGrpcClient: urlGrpcClient,
		validate:      validate,
		serverDomain:  serverDomain,
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
	msg := fmt.Sprintf("got url: %s", shortUrl)
	h.logger.Info(msg)

	if shortUrl == "" {
		msg := "short url is empty"
		response.BadRequest(w, msg)
		return
	}

	longURLResp, err := h.urlGrpcClient.FollowUrl(context.Background(), &url.ShortUrlRequest{
		ShortUrl: shortUrl,
	})
	if err != nil {
		h.logger.Error(err.Error())
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.NotFound {
				response.NotFound(w, err.Error())
				return
			}
		}

		response.InternalServerError(w)

	}

	http.Redirect(w, r, longURLResp.LongUrl, http.StatusFound)
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

	err = h.validate.Struct(longURLData)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	shortURLRaw, err := h.urlGrpcClient.ShortenUrl(context.Background(), &url.LongUrlRequest{
		LongUrl: longURLData.LongURL,
	})
	if err != nil {
		h.logger.Error(err.Error())
		response.InternalServerError(w)
		return
	}

	shortURL := fmt.Sprintf("%s://%s/%s", serverProtocol, h.serverDomain, shortURLRaw.ShortUrl)
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

func (h *URLHandler) SaveURLOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Request-Method", "POST")
	w.Header().Add("Access-Control-Request-Headers", "x-requested-with")
	w.Header().Add("Origin", "*")
}