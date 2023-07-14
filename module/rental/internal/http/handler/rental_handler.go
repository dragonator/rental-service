package handler

import (
	"context"
	"net/http"

	"github.com/dragonator/rental-service/module/rental/internal/http/contract"
	"github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
	"github.com/go-chi/chi"
)

// RentalFetchingOp is a contract to a rental fetching operation.
type RentalFetchingOp interface {
	GetRentalByID(ctx context.Context, rentalID string) (*model.Rental, error)
	ListRentals(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error)
}

// RentalHandler holds implementation of handlers for rentals.
type RentalHandler struct {
	rentalFetchingOp RentalFetchingOp
}

// NewRentalHandler is a construction function for RentalHandler.
func NewRentalHandler(rentalFetchingOp RentalFetchingOp) *RentalHandler {
	return &RentalHandler{
		rentalFetchingOp: rentalFetchingOp,
	}
}

// GetRentalByID returns a handle that is fetching a rental by id.
func (rh *RentalHandler) GetRentalByID(method, path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rental, err := rh.rentalFetchingOp.GetRentalByID(r.Context(), chi.URLParam(r, "id"))
		if err != nil {
			errorResponse(w, err)
			return
		}

		successResponse(w, toGetRentalByIDResponse(rental))

		return
	}
}

func toGetRentalByIDResponse(rental *model.Rental) *contract.GetRentalByIDResponse {
	return &contract.GetRentalByIDResponse{
		ID:              rental.ID,
		Name:            rental.Name,
		Description:     rental.Description,
		Type:            rental.Type,
		Make:            rental.VehicleMake,
		Model:           rental.VehicleModel,
		Year:            rental.VehicleYear,
		Length:          rental.VehicleLength,
		Sleeps:          rental.Sleeps,
		PrimaryImageURL: rental.PrimaryImageURL,
		Price: contract.Price{
			Day: rental.PricePerDay,
		},
		Location: contract.Location{
			City:      rental.HomeCity,
			State:     rental.HomeState,
			Zip:       rental.HomeZip,
			Country:   rental.HomeCountry,
			Latitude:  rental.Latitude,
			Longitude: rental.Longitude,
		},
	}
}
