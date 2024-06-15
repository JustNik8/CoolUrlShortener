package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"api_gateway/internal/converter"
	"api_gateway/internal/transport/rest/dto"
	"api_gateway/internal/transport/rest/response"
	"api_gateway/pkg/proto/analytics"
	"google.golang.org/grpc/status"
)

const (
	limitQueryParam = "limit"
	pageQueryParam  = "page"
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
//	@Param			page	query		int	true	"Страница"
//	@Param			limit	query		int	true	"Максимальное количество url на странице"
//	@Success		200		{object}	dto.TopURLDataResponse
//	@Failure		400		{object}	response.Body
//	@Failure		500		{object}	response.Body
//	@Router			/api/top_urls [get]
func (h *AnalyticsHandler) GetTopURLs(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Credentials", "true")

	pageRaw := r.URL.Query().Get(pageQueryParam)
	if pageRaw == "" {
		msg := fmt.Sprintf("Query param '%s' is empty", pageQueryParam)
		response.BadRequest(w, msg)
		return
	}

	page, err := strconv.Atoi(pageRaw)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	limitRaw := r.URL.Query().Get(limitQueryParam)
	if limitRaw == "" {
		msg := fmt.Sprintf("Query param '%s' is empty", limitQueryParam)
		response.BadRequest(w, msg)
		return
	}

	limit, err := strconv.Atoi(limitRaw)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	topUrlsGrpcResp, err := h.analyticsGrpcClient.GetTopUrls(context.Background(), &analytics.TopUrlsRequest{
		Page:  int64(page),
		Limit: int64(limit),
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			h.logger.Error(st.Code().String())
			h.logger.Error(st.Message())
		}
		h.logger.Error(err.Error())

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
