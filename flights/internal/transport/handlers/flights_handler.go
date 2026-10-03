package handlers

import (
	"encoding/json"
	"errors"
	domain_err "flights/internal/domain/errors"
	"flights/internal/domain/models"
	"flights/internal/domain/repository"
	"flights/internal/domain/service"
	"flights/internal/transport/dto"
	"flights/internal/utils"
	"fmt"
	"net/http"
	"pkg/httperr"
	"pkg/httpresp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type FlightHandler struct {
	storage    repository.FlightRepository
	flight_srv service.FlightService
	validator  *validator.Validate
}

func NewFlightHandler(storage repository.FlightRepository, service service.FlightService) (*FlightHandler, error) {
	return &FlightHandler{
		storage:    storage,
		flight_srv: service,
		validator:  validator.New(),
	}, nil
}

// Handler for getting flights
func (h *FlightHandler) HandleGetFlights() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := utils.ParseGetFlightsQuery(r.URL.Query())
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid request parameters")
			return
		}

		req.Normalize()
		if err := h.validator.Struct(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		flights, err := h.flight_srv.GetFlights(ctx, req.ToFlightFilter())
		if err != nil {
			switch {
			case errors.Is(err, domain_err.ErrFlightNotFound):
				httperr.Write(w, http.StatusNotFound, "Flights not found")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		data := dto.FlightListResponse{
			Data: dto.FlightsToResponses(flights),
			Pagination: dto.PaginationResponse{
				Page:  req.Page,
				Limit: req.Limit,
			},
		}
		httpresp.OK(w, data)
	}
}

// Handler for creating flight
func (h *FlightHandler) HandleCreateFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := &dto.CreateFlightRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		req.Normalize()
		if err := h.validator.Struct(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		created, err := h.flight_srv.CreateFlight(ctx, req.FlightNumber, req.Origin, req.Destination, req.Date, models.FlightStatus(req.Status), req.Aircraft)
		if err != nil {
			httperr.Write(w, http.StatusInternalServerError, "Failed to create flight")
			return
		}

		data := dto.FlightToResponse(created)
		httpresp.Created(w, data)
	}
}

func (h *FlightHandler) HandleDeleteFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Failed to parse query param `id`. Should be UUID")
			return
		}

		err = h.flight_srv.DeleteFlight(ctx, id)
		if err != nil {
			if errors.Is(err, domain_err.ErrFlightNotFound) {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			httperr.Write(w, http.StatusInternalServerError, "Failed to delete flight")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *FlightHandler) HandleCancelFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Failed to parse param `id`. Should be UUID")
			return
		}

		err = h.flight_srv.CancelFlight(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, domain_err.ErrFlightNotFound):
				httperr.Write(w, http.StatusNotFound, "Flight not found")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
		}

		httpresp.OK(w, map[string]string{"message": "flight cancelled"})
	}
}

func (h *FlightHandler) HandleRescheduleFlight() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Failed to parse param `id`. Should be UUID")
			return
		}

		req := &dto.RescheduleFlightRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if err := h.validator.Struct(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.flight_srv.RescheduleFlight(ctx, id, req.NewDate)
		if err != nil {
			switch {
			case errors.Is(err, domain_err.ErrFlightNotFound):
				httperr.Write(w, http.StatusNotFound, "flights not found")
			case errors.Is(err, domain_err.ErrIncorrectFlightTime):
				httperr.Write(w, http.StatusBadRequest, "new time in future")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		httpresp.OK(w, map[string]string{"message": fmt.Sprintf("flight rescheduled to %s", req.NewDate)})
	}
}

func (h *FlightHandler) HandleChangeStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "Failed to parse param `id`. Should be UUID")
			return
		}

		req := &dto.ChangeStatusRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if err := h.validator.Struct(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.flight_srv.ChangeStatus(ctx, id, req.Status)

		if err != nil {
			switch {
			case errors.Is(err, domain_err.ErrFlightNotFound):
				httperr.Write(w, http.StatusNotFound, "flights not found")
			case errors.Is(err, domain_err.ErrInvalidFlightStatus):
				httperr.Write(w, http.StatusBadRequest, "invalid flight")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		httpresp.OK(w, map[string]string{"message": fmt.Sprintf("flight status changed to %s", req.Status)})
	}
}
