package util

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// StringPtr returns a pointer to the given string.
func StringPtr(s string) *string {
	return &s
}

// WriteJSONResponse marshals data to JSON and writes it to the response writer.
// It handles JSON encoding errors by logging and returning a 500 error.
func WriteJSONResponse(w http.ResponseWriter, logger *slog.Logger, statusCode int, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal JSON response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		logger.Error("Failed to write response", "error", err)
	}
}
