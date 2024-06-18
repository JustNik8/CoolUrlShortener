package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"api_gateway/errs"
	"api_gateway/internal/client"
	"api_gateway/internal/transport/rest/response"
)

const (
	limitQueryParam = "limit"
	pageQueryParam  = "page"
	defaultPage     = 1
	defaultLimit    = 10
)

type AnalyticsHandler struct {
	logger          *slog.Logger
	analyticsClient client.AnalyticsClient
}

func NewAnalyticsHandler(
	logger *slog.Logger,
	analyticsClient client.AnalyticsClient,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		logger:          logger,
		analyticsClient: analyticsClient,
	}
}

// GetTopURLs docs
//
//	@Summary		Получение списка популярных url
//	@Tags			url
//	@Description	Принимает page и limit. Возвращает список популярных url. Поддерживает пагинацию
//	@ID				get-top-urls
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Страница"
//	@Param			limit	query		int	false	"Максимальное количество url на странице"
//	@Success		200		{object}	dto.TopURLDataResponse
//	@Failure		400		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/api/top_urls [get]
func (h *AnalyticsHandler) GetTopURLs(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}

	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	page, err := h.parseQueryParam(r, pageQueryParam, defaultPage)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	limit, err := h.parseQueryParam(r, limitQueryParam, defaultLimit)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	topUrlsResp, err := h.analyticsClient.GetTopUrls(context.Background(), int64(page), int64(limit))

	if err != nil {
		if errors.Is(err, errs.ErrInvalidArgument) {
			response.BadRequest(w, "bad params")
		}
		response.InternalServerError(w)
		return
	}

	respBytes, err := json.Marshal(topUrlsResp)
	if err != nil {
		h.logger.Error(err.Error())
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, respBytes)
}

func (h *AnalyticsHandler) parseQueryParam(r *http.Request, key string, defaultValue int) (int, error) {
	queryParam := r.URL.Query().Get(key)

	if queryParam == "" {
		return defaultValue, nil
	}

	param, err := strconv.Atoi(queryParam)
	if err != nil {
		return 0, err
	}

	if param == 0 {
		return defaultValue, nil
	}
	return param, nil

}
