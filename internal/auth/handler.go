package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nazarkozynets/log-viewer-api/internal/response"
)

type Handler struct {
	client *Client
}

func NewHandler(client *Client) *Handler {
	return &Handler{
		client: client,
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid login request")
		return
	}

	result, statusCode, err := h.client.Login(loginRequest)
	if err != nil {
		response.Error(w, statusCode, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, result)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
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

	response.JSON(w, http.StatusOK, me)
}

func isAdminRole(role string) bool {
	return strings.EqualFold(role, "admin")
}
