package auth

import (
	"net/http"

	"github.com/nazarkozynets/log-viewer-api/internal/response"
)

func (h *Handler) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		me, statusCode, err := h.client.GetMe(r)
		if err != nil {
			response.Error(w, statusCode, err.Error())
			return
		}

		if statusCode == http.StatusUnauthorized {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		if statusCode == http.StatusForbidden {
			response.Error(w, http.StatusForbidden, "forbidden")
			return
		}

		if !isAdminRole(me.Role) {
			response.Error(w, http.StatusForbidden, "admin access required")
			return
		}

		next.ServeHTTP(w, r)
	})
}
