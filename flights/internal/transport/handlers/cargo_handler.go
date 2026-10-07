package handlers

import (
	"errors"
	usecase "flights/internal/application/usecase"
	domain_err "flights/internal/domain/errors"
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

		id, err := uuid.Parse(r.PathValue("flight_id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid flight id")
			return
		}

		manifest, err := h.flightUseCase.GetCargoManifest(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, domain_err.ErrCargoManifestNotFound):
				httperr.Write(w, http.StatusNotFound, "Cargo manifest not found")
			default:
				httperr.Write(w, http.StatusInternalServerError, "internal server error")
			}
		}

		// mapping to DTO

		httpresp.OK(w, manifest)
	}
}
