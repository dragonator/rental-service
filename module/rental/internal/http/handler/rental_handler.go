package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dragonator/rental-service/module/rental/internal/http/contract"
	"github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
	"github.com/go-chi/chi"
)

var _invalidParameter = "invalid parameter: %s"

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
	ids := parseIDs(r)
	priceMin, errMin := parsePriceMin(r)
	priceMax, errMax := parsePriceMax(r)
	near, errNear := parseNear(r)
	sort, errSort := parseSort(r)
	limit, errLimit := parseLimit(r)
	offset, errOffset := parseOffset(r)

	err := errors.Join(errMin, errMax, errNear, errSort, errLimit, errOffset)
	if err != nil {
		return nil, err
	}

	return &storage.RentalFilters{
		IDs:      ids,
		PriceMin: priceMin,
		PriceMax: priceMax,
		Near:     near,
		OrderBy:  sort,
		Pagination: storage.Pagination{
			Limit:  limit,
			Offset: offset,
		},
	}, nil
}

func parseIDs(r *http.Request) []string {
	if ids := r.URL.Query().Get("ids"); len(ids) != 0 {
		return strings.Split(ids, ",")
	}

	return nil
}

func parsePriceMin(r *http.Request) (*int64, error) {
	if priceMin := r.URL.Query().Get("price_min"); len(priceMin) != 0 {
		pm, err := strconv.ParseInt(priceMin, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(_invalidParameter, "price_min")
		}

		return &pm, nil
	}

	return nil, nil
}

func parsePriceMax(r *http.Request) (*int64, error) {
	if priceMax := r.URL.Query().Get("price_max"); len(priceMax) != 0 {
		pm, err := strconv.ParseInt(priceMax, 10, 64)
		if err != nil {
			return nil, fmt.Errorf(_invalidParameter, "price_max")
		}

		return &pm, nil
	}

	return nil, nil
}

func parseNear(r *http.Request) (*storage.Location, error) {
	if near := r.URL.Query().Get("near"); len(near) != 0 {
		near := strings.Split(near, ",")

		if len(near) != 2 {
			return nil, fmt.Errorf(_invalidParameter, "near")
		}

		lat, err1 := strconv.ParseFloat(near[0], 64)
		lng, err2 := strconv.ParseFloat(near[1], 64)
		if err1 != nil || err2 != nil {
			return nil, errors.New("invalid argument: near")
		}

		return &storage.Location{
			Latitude:  float32(lat),
			Longitude: float32(lng),
		}, nil
	}

	return nil, nil
}

func parseSort(r *http.Request) (*string, error) {
	if sort := r.URL.Query().Get("sort"); len(sort) != 0 {
		if _, ok := storage.RentalSortFields[sort]; !ok {
			return nil, errors.New("invalid argument: near")
		}

		return &sort, nil
	}

	return nil, nil
}

func parseLimit(r *http.Request) (*int, error) {
	if limit := r.URL.Query().Get("limit"); len(limit) != 0 {
		lim, err := strconv.Atoi(limit)
		if err != nil {
			return nil, fmt.Errorf(_invalidParameter, "limit")
		}

		return &lim, nil
	}

	return nil, nil
}

func parseOffset(r *http.Request) (*int, error) {
	if offset := r.URL.Query().Get("offset"); len(offset) != 0 {
		os, err := strconv.Atoi(offset)
		if err != nil {
			return nil, fmt.Errorf(_invalidParameter, "offset")
		}

		return &os, nil
	}

	return nil, nil
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
	}
}

func toListRentalsResponse(rentals model.Rentals) *contract.ListRentalsResponse {
	resp := contract.ListRentalsResponse{}

	for _, r := range rentals {
		resp = append(resp, toRentalContract(r))
	}

	return &resp
}
