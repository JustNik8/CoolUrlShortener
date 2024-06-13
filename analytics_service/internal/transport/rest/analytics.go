package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"analytics_service/internal/converter"
	"analytics_service/internal/service"
	"analytics_service/internal/transport/rest/response"
)

const (
	limitQueryParam = "limit"
)

type AnalyticsHandler struct {
	logger           *slog.Logger
	analyticsService service.AnalyticsService
	topURLConverter  *converter.TopURLConverter
}

func NewAnalyticsHandler(
	logger *slog.Logger,
	analyticsService service.AnalyticsService,
	topURLConverter *converter.TopURLConverter,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		logger:           logger,
		analyticsService: analyticsService,
		topURLConverter:  topURLConverter,
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

	topURLs, err := h.analyticsService.GetTopUrls(context.Background(), limit)
	if err != nil {
		response.InternalServerError(w)
		return
	}

	topURLsDto := h.topURLConverter.ConvertSliceDomainToDto(topURLs)
	respBytes, err := json.Marshal(topURLsDto)
	if err != nil {
		h.logger.Error(err.Error())
		response.InternalServerError(w)
		return
	}

	response.WriteResponse(w, http.StatusOK, respBytes)
}
