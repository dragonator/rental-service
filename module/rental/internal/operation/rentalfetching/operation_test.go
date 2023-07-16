package rentalfetching

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/dragonator/rental-service/module/rental/internal/http/service/svc"
	"github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
	"github.com/google/go-cmp/cmp"
)

var (
	_rentals = model.Rentals{
		{
			ID:              1,
			UserID:          2,
			Name:            "Rental 1",
			Type:            "Type 1",
			Description:     "Description 1",
			Sleeps:          4,
			PricePerDay:     1000,
			HomeCity:        "City 1",
			HomeState:       "State 1",
			HomeZip:         "Zip 1",
			HomeCountry:     "Country 1",
			VehicleMake:     "Make 1",
			VehicleModel:    "Model 1",
			VehicleYear:     2022,
			VehicleLength:   10.5,
			Latitude:        40.1234,
			Longitude:       -75.5678,
			PrimaryImageURL: "ImageURL 1",
			User: &model.User{
				ID:        2,
				FirstName: "FirstName 1",
				LastName:  "LastName 1",
			},
		},
		{
			ID:              2,
			UserID:          3,
			Name:            "Rental 2",
			Type:            "Type 2",
			Description:     "Description 2",
			Sleeps:          6,
			PricePerDay:     1500,
			HomeCity:        "City 2",
			HomeState:       "State 2",
			HomeZip:         "Zip 2",
			HomeCountry:     "Country 2",
			VehicleMake:     "Make 2",
			VehicleModel:    "Model 2",
			VehicleYear:     2023,
			VehicleLength:   12.5,
			Latitude:        35.6789,
			Longitude:       -80.9012,
			PrimaryImageURL: "ImageURL 2",
			User: &model.User{
				ID:        3,
				FirstName: "FirstName 2",
				LastName:  "LastName 2",
			},
		},
	}
)

func TestOperation_GetRentalByID(t *testing.T) {
	testCases := []struct {
		name            string
		mockRentalStore *RentalStoreMock
		rentalID        int
		expectedResult  *model.Rental
		expectedErr     error
	}{
		{
			name: "Existing rental",
			mockRentalStore: &RentalStoreMock{
				GetByIDFunc: func(ctx context.Context, rentalID int) (*model.Rental, error) {
					return _rentals[0], nil
				},
			},
			rentalID:       int(_rentals[0].ID),
			expectedResult: _rentals[0],
			expectedErr:    nil,
		},
		{
			name: "Not found error",
			mockRentalStore: &RentalStoreMock{
				GetByIDFunc: func(ctx context.Context, rentalID int) (*model.Rental, error) {
					return nil, sql.ErrNoRows
				},
			},
			rentalID:       int(_rentals[0].ID),
			expectedResult: nil,
			expectedErr:    svc.ErrNotFound,
		},
		{
			name: "Other error",
			mockRentalStore: &RentalStoreMock{
				GetByIDFunc: func(ctx context.Context, rentalID int) (*model.Rental, error) {
					return nil, sql.ErrConnDone
				},
			},
			expectedResult: nil,
			expectedErr:    sql.ErrConnDone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			operation := NewOperation(tc.mockRentalStore)

			rental, err := operation.GetRentalByID(ctx, tc.rentalID)

			calls := tc.mockRentalStore.GetByIDCalls()
			if len(calls) != 1 {
				t.Fatalf("Unexpected number of calls to GetRentalByID:\nexpected: 1\ngot      %d", len(calls))
			}

			if calls[0].Ctx != ctx {
				t.Fatalf("Unexpected context:\nexpected: %v\ngot:      %v", ctx, calls[0].Ctx)
			}

			if calls[0].RentalID != tc.rentalID {
				t.Fatalf("Unexpected rental id:\nexpected: %v\ngot:      %v", calls[0].RentalID, tc.rentalID)
			}

			if err != tc.expectedErr && !errors.Is(err, tc.expectedErr) {
				t.Fatalf("Unexpected error:\nexpected: %v\ngot:      %v", tc.expectedErr, err)
			}

			if !cmp.Equal(rental, tc.expectedResult) {
				t.Fatalf("Unxpected rental:\nexpected: %v\ngot:      %v", tc.expectedResult, rental)
			}
		})
	}
}

func TestOperation_ListRentals(t *testing.T) {
	testCases := []struct {
		name            string
		mockRentalStore *RentalStoreMock
		filters         *storage.RentalFilters
		expectedResult  model.Rentals
		expectedErr     error
	}{
		{
			name: "With all filters",
			mockRentalStore: &RentalStoreMock{
				ListFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			filters: &storage.RentalFilters{
				IDs:      []int32{1, 2},
				PriceMin: toPtr(int64(12)),
				PriceMax: toPtr(int64(-38)),
				Near: &storage.Location{
					Latitude:  34.3,
					Longitude: -83.23,
				},
				OrderBy: toPtr("price_per_day"),
				Pagination: storage.Pagination{
					Limit:  toPtr(10),
					Offset: toPtr(8),
				},
			},
			expectedResult: _rentals,
			expectedErr:    nil,
		},
		{
			name: "Without filters",
			mockRentalStore: &RentalStoreMock{
				ListFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			filters:        nil,
			expectedResult: _rentals,
			expectedErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			operation := NewOperation(tc.mockRentalStore)

			rentals, err := operation.ListRentals(ctx, tc.filters)

			calls := tc.mockRentalStore.ListCalls()
			if len(calls) != 1 {
				t.Fatalf("Unexpected number of calls to List:\nexpected: 1\ngot      %d", len(calls))
			}

			if calls[0].Ctx != ctx {
				t.Fatalf("Unexpected context:\nexpected: %v\ngot:      %v", ctx, calls[0].Ctx)
			}

			if !cmp.Equal(calls[0].Filters, tc.filters) {
				t.Fatalf("Unexpected rental id:\nexpected: %v\ngot:      %v", calls[0].Filters, tc.filters)
			}

			if err != tc.expectedErr && !errors.Is(err, tc.expectedErr) {
				t.Fatalf("Unexpected error:\nexpected: %v\ngot:      %v", tc.expectedErr, err)
			}

			if !cmp.Equal(rentals, tc.expectedResult) {
				t.Fatalf("Unxpected rental:\nexpected: %v\ngot:      %v", tc.expectedResult, rentals)
			}
		})
	}
}

func toPtr[T any](v T) *T {
	return &v
}
