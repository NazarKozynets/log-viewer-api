package server

import (
	"net/http"

	"github.com/nazarkozynets/log-viewer-api/internal/auth"
	"github.com/nazarkozynets/log-viewer-api/internal/logs"
	"github.com/nazarkozynets/log-viewer-api/internal/response"
)

func NewRouter(logsHandler *logs.Handler, authHandler *auth.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)

	mux.Handle("GET /api/auth/me", authHandler.RequireAdmin(
		http.HandlerFunc(authHandler.Me),
	))

	mux.Handle("GET /api/sources", authHandler.RequireAdmin(
		http.HandlerFunc(logsHandler.GetSources),
	))

	mux.Handle("GET /api/logs", authHandler.RequireAdmin(
		http.HandlerFunc(logsHandler.GetLogs),
	))

	return mux
}
