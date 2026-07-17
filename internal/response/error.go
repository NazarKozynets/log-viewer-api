package response

import "net/http"

type ErrorResponse struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	JSON(w, statusCode, ErrorResponse{
		Error: message,
	})
}
