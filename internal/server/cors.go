package server

import (
	"net/http"
	"strings"
)

func WithCORS(next http.Handler, allowedOrigin string) http.Handler {
	allowedOrigin = strings.TrimSpace(allowedOrigin)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if isAllowedOrigin(origin, allowedOrigin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string, allowedOrigin string) bool {
	if origin == "" || allowedOrigin == "" {
		return false
	}

	if allowedOrigin == "*" {
		return true
	}

	return origin == allowedOrigin
}
