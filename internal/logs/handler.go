package logs

import (
	"net/http"
	"strconv"

	"github.com/nazarkozynets/log-viewer-api/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetSources(w http.ResponseWriter, r *http.Request) {
	result := h.service.GetSources()
	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	query := GetLogsQuery{
		Level:     queryParams.Get("level"),
		Context:   queryParams.Get("context"),
		RequestID: queryParams.Get("requestId"),
		UserID:    queryParams.Get("userId"),
		Query:     queryParams.Get("q"),
		Source:    queryParams.Get("source"),
		Limit:     parseInt(queryParams.Get("limit")),
	}

	result := h.service.GetLogs(query)
	response.JSON(w, http.StatusOK, result)
}

func parseInt(value string) int {
	if value == "" {
		return 0
	}

	result, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return result
}
