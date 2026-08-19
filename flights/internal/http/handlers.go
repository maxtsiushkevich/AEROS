package http

import (
	"encoding/json"
	"io"
	"net/http"
)

// Hhandler for getting flights
func (s *Server) HandleGetFlights() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Parse request
		request, err := ParseGetFlightsQuery(r.URL.Query())
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "Invalid request parameters", http.StatusBadRequest)
			return
		}

		// Normalazing data
		request.Normalize()

		// Validate
		if err := s.validator.Struct(request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Call service
		flights, err := s.service.GetFlights(ctx, request.ToServiceQuery())
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
		ctx := r.Context()
		// Parse request body
		request := &CreateFlightRequest{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		// Normalazing data
		request.Normalize()

		// Validate
		if err := s.validator.Struct(request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Convert DTO to domain model
		flight := request.ToServiceFlight()

		err := s.service.CreateFlight(ctx, flight)
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

func (s *Server) HandlePatchFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Unmarshall JSON
		request := &PatchFlightRequest{}
		if err := json.Unmarshal(body, request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}

		request.Normalize()

		// Validate
		if err := s.validator.Struct(request); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		flight := request.ToFlightUpdate()

		updatedFlight, err := s.service.UpdateFlight(ctx, flight)

		// Add custrom error handling for not found case
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Error("Error while update flight", "err", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update flight"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := UpdatedFlightResponse{
			Data: FlightToResponse(updatedFlight),
		}
		json.NewEncoder(w).Encode(response)
	}
}
