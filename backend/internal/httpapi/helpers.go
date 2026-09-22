package httpapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to encode response: %v", err)
		status = http.StatusInternalServerError
		encoded = []byte(`{"error":{"code":"internal_error","message":"internal server error"}}`)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(append(encoded, '\n')); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
		},
	})
}

func writeDecodeError(w http.ResponseWriter, err error) {
	var sizeError *http.MaxBytesError
	if errors.As(err, &sizeError) {
		writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large")
		return
	}
	writeError(w, http.StatusBadRequest, "invalid_request", "invalid JSON request body")
}
