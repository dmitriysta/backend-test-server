package handlers

import (
	"back-end/logs"
	"net/http"
)

func handlePanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		logs.LogError("Internal server error: %v", r)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
