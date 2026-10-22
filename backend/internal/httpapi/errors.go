package httpapi

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details"`
}

func writeError(w http.ResponseWriter, status int, code, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(APIError{
		Error: ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WriteBadRequest(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func WriteNotFound(w http.ResponseWriter, message string) {
	writeError(w, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func WriteMethodNotAllowed(w http.ResponseWriter) {
	writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed", nil)
}

func WriteUnsupportedMediaType(w http.ResponseWriter) {
	writeError(w, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE", "Content-Type must be application/json", nil)
}

func WritePayloadTooLarge(w http.ResponseWriter) {
	writeError(w, http.StatusRequestEntityTooLarge, "PAYLOAD_TOO_LARGE", "Request body exceeds the 1 MiB limit", nil)
}

func WriteInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal server error occurred", nil)
}
