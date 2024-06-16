package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"api_gateway/internal/converter"
	"api_gateway/internal/transport/rest/dto"
	"api_gateway/internal/transport/rest/response"
	"api_gateway/pkg/proto/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	limitQueryParam = "limit"
	pageQueryParam  = "page"
	defaultPage     = 1
	defaultLimit    = 10
)

type AnalyticsHandler struct {
	logger              *slog.Logger
	analyticsGrpcClient analytics.AnalyticsClient
	topUrlConverter     converter.TopURLConverter
	paginationConverter converter.PaginationConverter
}

func NewAnalyticsHandler(
	logger *slog.Logger,
	grpcClient analytics.AnalyticsClient,
	topUrlConverter converter.TopURLConverter,
	paginationConverter converter.PaginationConverter,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		logger:              logger,
		analyticsGrpcClient: grpcClient,
		topUrlConverter:     topUrlConverter,
		paginationConverter: paginationConverter,
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

	topUrlsGrpcResp, err := h.analyticsGrpcClient.GetTopUrls(context.Background(), &analytics.TopUrlsRequest{
		Page:  int64(page),
		Limit: int64(limit),
	})

	if err != nil {
		h.logger.Error(err.Error())

		st, ok := status.FromError(err)
		if !ok || st.Code() == codes.Internal {
			response.InternalServerError(w)
			return
		}

		if st.Code() == codes.InvalidArgument {
			response.BadRequest(w, st.Message())
			return
		}

		response.InternalServerError(w)
		return
	}

	topUrlsResp := dto.TopURLDataResponse{
		TopURLData: h.topUrlConverter.MapSlicePbToDto(topUrlsGrpcResp.TopUrlData),
		Pagination: h.paginationConverter.MapPbToDto(topUrlsGrpcResp.Pagination),
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
