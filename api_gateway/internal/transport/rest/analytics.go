package rest

//import (
//	"fmt"
//	"log/slog"
//	"net/http"
//	"strconv"
//
//	"api_gateway/internal/transport/rest/response"
//)
//
//const (
//	limitQueryParam = "limit"
//	pageQueryParam  = "page"
//)
//
//type AnalyticsHandler struct {
//	logger *slog.Logger
//}
//
//func NewAnalyticsHandler(
//	logger *slog.Logger,
//) *AnalyticsHandler {
//	return &AnalyticsHandler{
//		logger: logger,
//	}
//}
//
//func (h *AnalyticsHandler) GetTopURLs(w http.ResponseWriter, r *http.Request) {
//	limitRaw := r.URL.Query().Get(limitQueryParam)
//	if limitRaw == "" {
//		msg := fmt.Sprintf("Query param '%s' is empty", limitQueryParam)
//		response.BadRequest(w, msg)
//		return
//	}
//
//	limit, err := strconv.Atoi(limitRaw)
//	if err != nil {
//		response.BadRequest(w, err.Error())
//		return
//	}
//
//	pageRaw := r.URL.Query().Get(pageQueryParam)
//	if pageRaw == "" {
//		msg := fmt.Sprintf("Query param '%s' is empty", pageQueryParam)
//		response.BadRequest(w, msg)
//		return
//	}
//
//	page, err := strconv.Atoi(pageRaw)
//	if err != nil {
//		response.BadRequest(w, err.Error())
//		return
//	}
//
//	response.WriteResponse(w, http.StatusOK, respBytes)
//}
