package handlers

import (
	"encoding/json"
	"errors"
	"flights/internal/dto"
	"flights/internal/service"
	"flights/internal/storage"
	"flights/internal/utils"
	"fmt"
	"io"
	"net/http"

	"github.com/maxtsiushkevich/AEROS/pkg/httperr"
	"github.com/maxtsiushkevich/AEROS/pkg/httpresp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type FlightHandler struct {
	storage   storage.FlightsStorage
	service   *service.FlightService
	validator *validator.Validate
}

func NewFlightHandler(
	storage storage.FlightsStorage,
	validator *validator.Validate) (*FlightHandler, error) {

	if err := storage.Open(); err != nil {
		return nil, err
	}

	return &FlightHandler{
		storage:   storage,
		service:   service.CreateFlightService(storage),
		validator: validator,
	}, nil
}

// Handler for getting flights
func (h *FlightHandler) HandleGetFlights() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Parse request
		request, err := utils.ParseGetFlightsQuery(r.URL.Query())
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid request parameters")
			return
		}

		// Normalazing data
		request.Normalize()

		// Validate
		if err := h.validator.Struct(request); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		// Call service
		flights, err := h.service.GetFlights(ctx, request.ToServiceQuery())
		if err != nil {
			httperr.Write(w, http.StatusNotFound, "Flights not found")
			return
		}

		data := dto.FlightsToResponses(flights)
		httpresp.OK(w, data)
	}
}

// Handler for creating flight
func (h *FlightHandler) HandleCreateFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Parse request body
		request := &dto.CreateFlightRequest{}
		if err := json.NewDecoder(r.Body).Decode(request); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		// Normalazing data
		request.Normalize()

		// Validate
		if err := h.validator.Struct(request); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		// Convert DTO to domain model
		flight := request.ToServiceFlight()

		created, err := h.service.CreateFlight(ctx, flight)
		if err != nil {
			httperr.Write(w, http.StatusInternalServerError, "Failed to create flight")
			return
		}

		data := dto.FlightToResponse(created)
		httpresp.Created(w, data)
	}
}

func (h *FlightHandler) HandlePatchFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			httperr.Write(w, http.StatusInternalServerError, "Failed to read request body")
			return
		}
		defer r.Body.Close()

		// Unmarshall JSON
		request := &dto.PatchFlightRequest{}
		if err := json.Unmarshal(body, request); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		request.Normalize()

		// Validate
		if err := h.validator.Struct(request); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		flight := request.ToFlightUpdate()

		updatedFlight, err := h.service.UpdateFlight(ctx, flight)

		// Add custrom error handling for not found case
		if err != nil {
			httperr.Write(w, http.StatusInternalServerError, "Failed to update flight")
			return
		}

		data := dto.FlightToResponse(updatedFlight)
		httpresp.OK(w, data)
	}
}

func (h *FlightHandler) HandleDeleteFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.URL.Query().Get("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Failed to parse query param `id`. Should be UUID")
			return
		}

		err = h.service.DeleteFlight(ctx, id)
		if err != nil {
			if errors.Is(err, storage.ErrFlightNotFound) {
				httperr.Write(w, http.StatusNotFound, fmt.Sprintf("Flight with id=%s not found", id))
				return
			}

			httperr.Write(w, http.StatusInternalServerError, "Failed to delete flight")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
