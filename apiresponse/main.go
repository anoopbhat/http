// Package apiresponse ...
package apiresponse

import (
	"encoding/json"
	"net/http"
)

type responseBody struct {
	Message string `json:"message"`
}

// Error ...
func Error(w http.ResponseWriter, message string, statusCode int) {
	if message == "" {
		message = http.StatusText(statusCode)
	}

	JSON(w, statusCode, &responseBody{Message: message})
}

// JSON ...
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	if data == nil {
		data = &responseBody{Message: http.StatusText(statusCode)}
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
}
