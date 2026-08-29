package httpresp

import (
	"encoding/json"
	"net/http"
)

type envelope[T any] struct {
	Data T `json:"data"`
}

func OK[T any](w http.ResponseWriter, data T) {
	writeJSON(w, http.StatusOK, envelope[T]{Data: data})
}

func Created[T any](w http.ResponseWriter, data T) {
	writeJSON(w, http.StatusCreated, envelope[T]{Data: data})
}

func writeJSON[T any](w http.ResponseWriter, status int, data T) {
	if w == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}
