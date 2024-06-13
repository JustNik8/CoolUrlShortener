package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"analytics_service/internal/converter"
	"analytics_service/internal/domain"
	"analytics_service/internal/service"
	"analytics_service/internal/transport/rest/dto"
	"analytics_service/internal/transport/rest/response"
)

const (
	urlEventsCounterTableName = "url_events_counter"
	limitQueryParam           = "limit"
	pageQueryParam            = "page"
)

type AnalyticsHandler struct {
	logger              *slog.Logger
	analyticsService    service.AnalyticsService
	paginationService   service.PaginationService
	topURLConverter     converter.TopURLConverter
	paginationConverter converter.PaginationConverter
}

func NewAnalyticsHandler(
	logger *slog.Logger,
	analyticsService service.AnalyticsService,
	paginationService service.PaginationService,
	topURLConverter converter.TopURLConverter,
	paginationConverter converter.PaginationConverter,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		logger:              logger,
		analyticsService:    analyticsService,
		paginationService:   paginationService,
		topURLConverter:     topURLConverter,
		paginationConverter: paginationConverter,
	}
}

func (h *AnalyticsHandler) GetTopURLs(w http.ResponseWriter, r *http.Request) {
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

	paginationParams := domain.PaginationParams{
		Page:  page,
		Limit: limit,
	}

	topURLs, err := h.analyticsService.GetTopUrls(context.Background(), paginationParams)
	if err != nil {
		response.InternalServerError(w)
		return
	}
	pagination, err := h.paginationService.GetPaginationInfo(urlEventsCounterTableName, paginationParams)
	if err != nil {
		h.logger.Error(err.Error())
		response.InternalServerError(w)
		return
	}

	topURLDataResp := dto.TopURLDataResponse{
		TopURLData: h.topURLConverter.MapSliceDomainToDto(topURLs),
		Pagination: h.paginationConverter.MapDomainToDto(pagination),
	}

	respBytes, err := json.Marshal(topURLDataResp)
	if err != nil {
		h.logger.Error(err.Error())
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, respBytes)
}
