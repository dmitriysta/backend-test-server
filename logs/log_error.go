package logs

import (
	"back-end/internal/handlers"
	"encoding/json"
	"log"
	"net/http"
)

func LogError(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		LogError("Failed to encode JSON: %v", err)
	}
}

func SendError(w http.ResponseWriter, statusCode int, message string) {
	SendJSONResponse(w, statusCode, handlers.ErrorResponse{Status: statusCode, Message: message})
}
