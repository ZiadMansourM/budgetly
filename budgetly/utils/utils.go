package utils

import (
	"encoding/json"
	"net/http"
)

// WriteJson writes a JSON response with the given status code and data
func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
