package handlers

import (
	"encoding/json"
	"errors"
	usecase "flights/internal/application/usecase"
	domain_err "flights/internal/domain/errors"
	"flights/internal/transport/dto"
	"net/http"
	"pkg/httperr"
	"pkg/httpresp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CargoHandler struct {
	flightUseCase usecase.CargoUseCase
	validator     *validator.Validate
}

func NewCargoHandler(useCase usecase.CargoUseCase) (*CargoHandler, error) {
	return &CargoHandler{
		flightUseCase: useCase,
		validator:     validator.New(),
	}, nil
}

func (h *CargoHandler) HandleGetCargoManifest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid id")
			return
		}

		manifest, err := h.flightUseCase.GetCargoManifest(ctx, usecase.GetCargoManifestRequest{ID: id})
		if err != nil {
			switch {
			case errors.Is(err, domain_err.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "Cargo manifest not found")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// mapping to DTO
		response := dto.CargoManifestToResponse(manifest)

		httpresp.OK(w, response)
	}
}

func (h *CargoHandler) HandleGetCargoItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *CargoHandler) HandleAddCargoItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid id")
			return
		}

		req := &dto.AddCargoItemRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if err := h.validator.Struct(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		item, err := h.flightUseCase.AddCargoItem(ctx, usecase.AddCargoItemRequest{
			ManifestID:  id,
			CargoType:   req.CargoType,
			WeightKg:    req.WeightKg,
			PassengerID: req.PassengerID,
			Description: req.Description,
		})

		if err != nil {
			switch {
			case errors.Is(err, domain_err.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "cargo manifest not found")
			case errors.Is(err, domain_err.ErrIncorrectCargoWeigh):
				httperr.Write(w, http.StatusUnprocessableEntity, "incorrect cargo weight")
			case errors.Is(err, domain_err.ErrLuggageWithoutOwner):
				httperr.Write(w, http.StatusUnprocessableEntity, "luggage without owner")
			case errors.Is(err, domain_err.ErrMaxWeightExceeded):
				httperr.Write(w, http.StatusUnprocessableEntity, "max weight exceeded")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		httpresp.Created(w, dto.CargoItemToResponse(item))
	}
}

func (h *CargoHandler) HandleDeleteCargoItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *CargoHandler) HandleMoveCargoItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
