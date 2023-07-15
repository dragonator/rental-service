package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dragonator/rental-service/module/rental/internal/http/contract"
	"github.com/dragonator/rental-service/module/rental/internal/http/service/svc"
	"github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
	"github.com/go-chi/chi"
	"github.com/gorilla/schema"
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

		successResponse(w, toRentalContract(rental))

		return
	}
}

// ListRentals returns a handle that is listing rentals based on filters.
func (rh *RentalHandler) ListRentals(method, path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		filters, err := rentalFiltersFromRequest(r)
		if err != nil {
			errorResponse(w, err)
			return
		}

		rentals, err := rh.rentalFetchingOp.ListRentals(r.Context(), filters)
		if err != nil {
			errorResponse(w, err)
			return
		}

		successResponse(w, toListRentalsResponse(rentals))

		return
	}
}

func rentalFiltersFromRequest(r *http.Request) (*storage.RentalFilters, error) {
	var query contract.ListRentalsQuery

	if err := r.ParseForm(); err != nil {
		return nil, fmt.Errorf("%w: parsing form: %w", svc.ErrInvalidQueryParameters, err)
	}

	if err := schema.NewDecoder().Decode(&query, r.Form); err != nil {
		return nil, fmt.Errorf("%w: unmashalling query: %w", svc.ErrInvalidQueryParameters, err)
	}

	filters := &storage.RentalFilters{
		IDs:      query.Ids,
		PriceMin: query.PriceMin,
		PriceMax: query.PriceMax,
		OrderBy:  query.Sort,
		Pagination: storage.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}

	if len(query.Near) > 0 {
		if len(query.Near) != 2 {
			return nil, fmt.Errorf("%w: invalid number of values for near (expected 2)", svc.ErrInvalidQueryParameters)
		}

		filters.Near = &storage.Location{
			Latitude:  query.Near[0],
			Longitude: query.Near[1],
		}
	}

	return filters, nil
}

func toRentalContract(rental *model.Rental) *contract.Rental {
	return &contract.Rental{
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
		User: contract.User{
			ID:        rental.User.ID,
			FirstName: rental.User.FirstName,
			LastName:  rental.User.LastName,
		},
	}
}

func toListRentalsResponse(rentals model.Rentals) *contract.ListRentalsResponse {
	resp := contract.ListRentalsResponse{}

	for _, r := range rentals {
		resp = append(resp, toRentalContract(r))
	}

	return &resp
}
