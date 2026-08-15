package http

import (
	"encoding/json"
	"net/http"
)

func (s *Server) HandleFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flightNumber := r.PathValue("flightNumber")

		flight, err := s.service.GetFlightsByNumber(flightNumber)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Flight not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(flight)
	}
}
