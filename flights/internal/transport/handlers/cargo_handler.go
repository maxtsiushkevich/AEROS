package handlers

import (
	"encoding/json"
	"errors"
	usecase "flights/internal/application/usecase"
	domainerr "flights/internal/domain/errors"
	"flights/internal/transport/dto"
	"net/http"
	"pkg/httperr"
	"pkg/httpresp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CargoHandler struct {
	cargoUseCase usecase.CargoUseCase
	validator    *validator.Validate
}

func NewCargoHandler(useCase usecase.CargoUseCase) (*CargoHandler, error) {
	return &CargoHandler{
		cargoUseCase: useCase,
		validator:    validator.New(),
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

		manifest, err := h.cargoUseCase.GetCargoManifest(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, domainerr.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "Cargo manifest not found")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		response := dto.CargoManifestToResponse(manifest)

		httpresp.OK(w, response)
	}
}

func (h *CargoHandler) HandleGetCargoItems() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cargoId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid manifest id")
			return
		}

		items, err := h.cargoUseCase.GetCargoItems(ctx, cargoId)
		if err != nil {
			switch {
			case errors.Is(err, domainerr.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "cargo manifest not found")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		response := make([]dto.CargoItemResponse, 0, len(items))
		for i := range items {
			response = append(response, dto.CargoItemToResponse(&items[i]))
		}

		httpresp.OK(w, response)
	}
}

func (h *CargoHandler) HandleGetCargoItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cargoId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid manifest id")
			return
		}

		itemId, err := uuid.Parse(r.PathValue("item_id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid item id")
			return
		}

		item, err := h.cargoUseCase.GetCargoItem(ctx, usecase.GetCargoItemRequest{
			ManifestID: cargoId,
			ItemID:     itemId,
		})
		if err != nil {
			switch {
			case errors.Is(err, domainerr.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "cargo manifest not found")
			case errors.Is(err, domainerr.ErrCargoItemNotFound):
				httperr.Write(w, http.StatusNotFound, "cargo item not found")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		httpresp.OK(w, dto.CargoItemToResponse(item))
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

		item, err := h.cargoUseCase.AddCargoItem(ctx, usecase.AddCargoItemRequest{
			ManifestID:  id,
			CargoType:   req.CargoType,
			WeightKg:    req.WeightKg,
			PassengerID: req.PassengerID,
			Description: req.Description,
		})

		if err != nil {
			switch {
			case errors.Is(err, domainerr.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "cargo manifest not found")
			case errors.Is(err, domainerr.ErrIncorrectCargoWeigh):
				httperr.Write(w, http.StatusUnprocessableEntity, "incorrect cargo weight")
			case errors.Is(err, domainerr.ErrLuggageWithoutOwner):
				httperr.Write(w, http.StatusUnprocessableEntity, "luggage without owner")
			case errors.Is(err, domainerr.ErrMaxWeightExceeded):
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
		ctx := r.Context()

		cargoId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid manifest id")
			return
		}

		itemId, err := uuid.Parse(r.PathValue("item_id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid item id")
			return
		}

		err = h.cargoUseCase.DeleteCargoItem(ctx, cargoId, itemId)
		if err != nil {
			httperr.Write(w, http.StatusInternalServerError, "Failed to delete cargo item")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *CargoHandler) HandleMoveCargoItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		cargoId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid manifest id")
			return
		}

		req := &dto.MoveCargoItemRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		if err := h.validator.Struct(req); err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.cargoUseCase.MoveCargoItem(ctx, usecase.MoveCargoItemRequest{
			ItemID:     req.ItemID,
			CargoID:    cargoId,
			NewCargoID: req.NewCargoID,
		})

		if err != nil {
			switch {
			case errors.Is(err, domainerr.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "cargo manifest not found")
			case errors.Is(err, domainerr.ErrCargoItemNotFound):
				httperr.Write(w, http.StatusNotFound, "cargo item not found")
			case errors.Is(err, domainerr.ErrMaxWeightExceeded):
				httperr.Write(w, http.StatusUnprocessableEntity, "max weight exceeded")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
