package http

import (
	"encoding/json"
	"flights/internal/domain"
	"flights/internal/service"
	"net/http"
	"time"
)

func (s *Server) HandleGetFlights() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()

		q := FlightQueryParams{
			FlightNumber: queryParams.Get("flight_number"),
			Origin:       queryParams.Get("origin"),
			Destination:  queryParams.Get("destination"),
			Status:       queryParams.Get("status"),
		}

		if dateFrom := queryParams.Get("date_from"); dateFrom != "" {
			t, err := time.Parse(time.RFC3339, dateFrom)
			if err != nil {
				http.Error(w, "invalid date_from", http.StatusBadRequest)
				return
			}

			q.DateFrom = &t
		}

		if dateTo := queryParams.Get("date_to"); dateTo != "" {
			t, err := time.Parse(time.RFC3339, dateTo)
			if err != nil {
				http.Error(w, "invalid date_to", http.StatusBadRequest)
				return
			}

			q.DateTo = &t
		}

		err := s.validator.Struct(q)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		flights, err := s.service.GetFlights(service.FlightQuery{
			FlightNumber: q.FlightNumber,
			Origin:       q.Origin,
			Destination:  q.Destination,
			Status:       domain.Status(q.Status),
			DateFrom:     q.DateFrom,
			DateTo:       q.DateTo,
		})

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Flights not found"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		flightsResponses := make([]FlightResponse, 0)

		for _, val := range flights {
			flightsResponses = append(flightsResponses, FlightResponse{
				ID:           val.ID,
				FlightNumber: val.FlightNumber,
				Origin:       val.Origin,
				Destination:  val.Destination,
				Date:         val.Date,
				Status:       string(val.Status),
				Aircraft:     val.Aircraft,
			})
		}

		flightsList := FlightListResponse{
			Data: flightsResponses,
		}
		json.NewEncoder(w).Encode(flightsList)
	}
}

func (s *Server) HandleCreateFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flight := CreateFlightRequest{}
		if err := json.NewDecoder(r.Body).Decode(&flight); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "invalid JSON body	", http.StatusBadRequest)
			return
		}

		if err := s.validator.Struct(flight); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := s.service.CreateFlight(&service.CreateFlightModel{
			FlightNumber: flight.FlightNumber,
			Origin:       flight.Origin,
			Destination:  flight.Destination,
			Date:         flight.Date,
			Status:       flight.Status,
			Aircraft:     flight.Aircraft,
		})

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Error("Error while creating flight", "err", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
