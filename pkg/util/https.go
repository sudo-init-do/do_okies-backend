package util

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func Err(w http.ResponseWriter, code int, key string) {
	JSON(w, code, map[string]string{"error": key})
}
