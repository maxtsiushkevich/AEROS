package http

import (
	"encoding/json"
	"net/http"
)

// Hhandler for getting flights
func (s *Server) HandleGetFlights() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request
		request, err := ParseGetFlightsQuery(r.URL.Query())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "Invalid request parameters", http.StatusBadRequest)
			return
		}

		// There is normalazing data

		// Validate
		if err := s.validator.Struct(request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call service
		flights, err := s.service.GetFlights(request.ToServiceQuery())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Flights not found"})
			return
		}

		response := FlightListResponse{
			Data: FlightsToResponses(flights),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// Handler for creating flight
func (s *Server) HandleCreateFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		request := &CreateFlightRequest{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		// There is normalazing data

		// Validate
		if err := s.validator.Struct(request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Convert DTO to domain model
		flight := request.ToServiceFlight()

		err := s.service.CreateFlight(flight)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Error("Error while creating flight", "err", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create flight"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Flight created successfully"})
	}
}
