package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/google/go-cmp/cmp"

	"github.com/dragonator/rental-service/module/rental/internal/http/contract"
	"github.com/dragonator/rental-service/module/rental/internal/http/handler"
	"github.com/dragonator/rental-service/module/rental/internal/http/service/svc"
	"github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
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

func TestRentalHandler_GetRentalByID(t *testing.T) {
	testCases := []struct {
		name                 string
		rentalID             string
		mockRentalFetchingOp *RentalFetchingOpMock
		expectedRentalID     int
		expectedCalls        int
		expectedCode         int
		expectedRental       contract.GetRentalByIDResponse
		expectedError        contract.ErrorResponse
	}{
		{
			name:     "Valid rental ID",
			rentalID: "1",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				GetRentalByIDFunc: func(ctx context.Context, rentalID int) (*model.Rental, error) {
					return _rentals[0], nil
				},
			},
			expectedRentalID: 1,
			expectedCalls:    1,
			expectedCode:     http.StatusOK,
			expectedRental:   contract.GetRentalByIDResponse{Rental: *toRentalContract(_rentals[0])},
		},
		{
			name:     "Missing rental ID",
			rentalID: "3",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				GetRentalByIDFunc: func(ctx context.Context, rentalID int) (*model.Rental, error) {
					return nil, svc.ErrNotFound
				},
			},
			expectedRentalID: 3,
			expectedCalls:    1,
			expectedCode:     http.StatusNotFound,
			expectedError: contract.ErrorResponse{
				Error: "not found",
			},
		},
		{
			name:                 "Invalid rental ID",
			rentalID:             "invalid",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: id",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rentalHandler := handler.NewRentalHandler(tc.mockRentalFetchingOp)

			router := chi.NewRouter()
			router.Get("/rentals/{id}", rentalHandler.GetRentalByID("GET", "/rentals/{id}"))

			request := httptest.NewRequest("GET", "/rentals/"+tc.rentalID, nil)
			responseRecorder := httptest.NewRecorder()

			router.ServeHTTP(responseRecorder, request)

			calls := tc.mockRentalFetchingOp.GetRentalByIDCalls()
			if tc.expectedCalls != len(calls) {
				t.Fatalf("Unexpected number of calls to GetRentalByID:\nexpected: %d\ngot      %d", tc.expectedCalls, len(calls))
			}

			if tc.expectedCalls > 0 {
				if calls[0].RentalID != tc.expectedRentalID {
					t.Fatalf("Unexpected rental id:\nexpected: %v\ngot:      %v", calls[0].RentalID, tc.expectedRentalID)
				}
			}

			if responseRecorder.Code != tc.expectedCode {
				t.Fatalf("Unexpected status code:\nexpected %d\ngot:      %d", tc.expectedCode, responseRecorder.Code)
			}

			if tc.expectedCode == http.StatusOK {
				var responseBody contract.GetRentalByIDResponse

				err := json.NewDecoder(responseRecorder.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if !cmp.Equal(responseBody, tc.expectedRental) {
					t.Fatalf("Unexpected rental:\nexpected: %v\ngot:      %v", tc.expectedRental, responseBody)
				}
			} else {
				var errorResponse contract.ErrorResponse

				err := json.NewDecoder(responseRecorder.Body).Decode(&errorResponse)
				if err != nil {
					t.Fatalf("Failed to decode error response body: %v", err)
				}

				if !cmp.Equal(errorResponse, tc.expectedError) {
					t.Fatalf("Unexpected error message:\nexpected: %s\ngot      %s", tc.expectedError, errorResponse.Error)
				}
			}
		})
	}
}

func TestRentalHandler_ListRentals(t *testing.T) {
	testCases := []struct {
		name                 string
		query                string
		mockRentalFetchingOp *RentalFetchingOpMock
		expectedFilters      *storage.RentalFilters
		expectedCalls        int
		expectedCode         int
		expectedRental       contract.ListRentalsResponse
		expectedError        contract.ErrorResponse
	}{
		{
			name:  "Valid IDs",
			query: "?ids=1,2",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				IDs: []int32{1, 2},
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:                 "Invalid IDs",
			query:                "?ids=a,b",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: unmashalling query: schema: error converting value for index 0 of \"ids\"",
			},
		},
		{
			name:  "PriceMin",
			query: "?price_min=100",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				PriceMin: toPtr[int64](100),
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:                 "Invalid PriceMin",
			query:                "?price_min=invalid",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: unmashalling query: schema: error converting value for \"price_min\"",
			},
		},
		{
			name:  "PriceMax",
			query: "?price_max=10000",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				PriceMax: toPtr[int64](10000),
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:                 "Invalid PriceMax",
			query:                "?price_max=invalid",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: unmashalling query: schema: error converting value for \"price_max\"",
			},
		},
		{
			name:  "Near",
			query: "?near=13.28,-43.76",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				Near: &storage.Location{
					Latitude:  13.28,
					Longitude: -43.76,
				},
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:                 "Invalid number of values for Near",
			query:                "?near=1",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: invalid number of values for near (expected 2)",
			},
		},
		{
			name:                 "Invalid Near",
			query:                "?near=a,b",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: unmashalling query: schema: error converting value for index 0 of \"near\"",
			},
		},
		{
			name:  "Limit",
			query: "?limit=3",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				Pagination: storage.Pagination{
					Limit: toPtr(3),
				},
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:                 "Invalid Limit",
			query:                "?limit=invalid",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: unmashalling query: schema: error converting value for \"limit\"",
			},
		},
		{
			name:  "Offset",
			query: "?offset=8",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				Pagination: storage.Pagination{
					Offset: toPtr(8),
				},
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:                 "Invalid Offset",
			query:                "?offset=invalid",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: unmashalling query: schema: error converting value for \"offset\"",
			},
		},
		{
			name:  "Sort",
			query: "?sort=price_per_day",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				OrderBy: toPtr("price_per_day"),
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:                 "Invalid Sort",
			query:                "?sort=invalid",
			mockRentalFetchingOp: &RentalFetchingOpMock{},
			expectedCalls:        0,
			expectedCode:         http.StatusBadRequest,
			expectedError: contract.ErrorResponse{
				Error: "invalid query parameters: unexpected sort field: expected one of [id name type make model year length sleeps price_per_day]",
			},
		},
		{
			name:  "All filters",
			query: "?ids=1,2&price_min=100&price_max=10000&near=13.28,-43.76&limit=3&offset=8&sort=price_per_day",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return _rentals, nil
				},
			},
			expectedFilters: &storage.RentalFilters{
				IDs:      []int32{1, 2},
				PriceMin: toPtr[int64](100),
				PriceMax: toPtr[int64](10000),
				Near: &storage.Location{
					Latitude:  13.28,
					Longitude: -43.76,
				},
				Pagination: storage.Pagination{
					Limit:  toPtr(3),
					Offset: toPtr(8),
				},
				OrderBy: toPtr("price_per_day"),
			},
			expectedCalls: 1,
			expectedCode:  http.StatusOK,
			expectedRental: contract.ListRentalsResponse{
				toRentalContract(_rentals[0]),
				toRentalContract(_rentals[1]),
			},
		},
		{
			name:  "No results",
			query: "",
			mockRentalFetchingOp: &RentalFetchingOpMock{
				ListRentalsFunc: func(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
					return nil, nil
				},
			},
			expectedFilters: &storage.RentalFilters{},
			expectedCalls:   1,
			expectedCode:    http.StatusOK,
			expectedRental:  contract.ListRentalsResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rentalHandler := handler.NewRentalHandler(tc.mockRentalFetchingOp)

			router := chi.NewRouter()
			router.Get("/rentals", rentalHandler.ListRentals("GET", "/rentals"))

			request := httptest.NewRequest("GET", "/rentals"+tc.query, nil)
			responseRecorder := httptest.NewRecorder()

			router.ServeHTTP(responseRecorder, request)

			calls := tc.mockRentalFetchingOp.ListRentalsCalls()
			if tc.expectedCalls != len(calls) {
				t.Fatalf("Unexpected number of calls to ListRentals:\nexpected: %d\ngot      %d", tc.expectedCalls, len(calls))
			}

			if tc.expectedCalls > 0 {
				if !cmp.Equal(calls[0].Filters, tc.expectedFilters) {
					t.Fatalf("Unexpected filters:\nexpected: %v\ngot:      %v", calls[0].Filters, tc.expectedFilters)
				}
			}

			if responseRecorder.Code != tc.expectedCode {
				t.Fatalf("Unexpected status code:\nexpected %d\ngot:      %d", tc.expectedCode, responseRecorder.Code)
			}

			if tc.expectedCode == http.StatusOK {
				var responseBody contract.ListRentalsResponse

				err := json.NewDecoder(responseRecorder.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if !cmp.Equal(responseBody, tc.expectedRental) {
					t.Fatalf("Unexpected rental:\nexpected: %v\ngot:      %v", tc.expectedRental, responseBody)
				}
			} else {
				var errorResponse contract.ErrorResponse

				err := json.NewDecoder(responseRecorder.Body).Decode(&errorResponse)
				if err != nil {
					t.Fatalf("Failed to decode error response body: %v", err)
				}

				if !cmp.Equal(errorResponse, tc.expectedError) {
					t.Fatalf("expected error message %s, but got %s", tc.expectedError, errorResponse.Error)
				}
			}
		})
	}
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
func toPtr[T any](v T) *T {
	return &v
}
