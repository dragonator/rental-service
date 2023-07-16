package storage_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	model "github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
	"github.com/dragonator/rental-service/pkg/config"
)

var (
	_nearThresholdRadius = 100
	_rentals             = model.Rentals{
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
	_columns = []string{
		"rentals.id",
		"rentals.user_id",
		"rentals.name",
		"rentals.type",
		"rentals.description",
		"rentals.sleeps",
		"rentals.price_per_day",
		"rentals.home_city",
		"rentals.home_state",
		"rentals.home_zip",
		"rentals.home_country",
		"rentals.vehicle_make",
		"rentals.vehicle_model",
		"rentals.vehicle_year",
		"rentals.vehicle_length",
		"rentals.lat",
		"rentals.lng",
		"rentals.primary_image_url",
		"users.id as users_id",
		"users.first_name",
		"users.last_name",
	}
)

func TestRentalRepository_GetByID(t *testing.T) {
	selectQuery := fmt.Sprintf(
		"SELECT %s FROM rentals JOIN users ON users.id = rentals.user_id WHERE rentals.id = \\$1",
		strings.Join(_columns, ", "),
	)

	testCases := []struct {
		name           string
		idParam        int
		expectedRental *model.Rental
		expectedError  error
		mockFunc       func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "Valid rental ID",
			idParam:        1,
			expectedRental: _rentals[0],
			expectedError:  nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(selectQuery).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(sqlmock.NewRows(_columns).AddRow(rentalValues(_rentals[0])...))
			},
		},
		{
			name:           "Missing rental ID",
			idParam:        77,
			expectedRental: nil,
			expectedError:  sql.ErrNoRows,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(selectQuery).
					WithArgs([]driver.Value{77}...).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %s", err)
			}
			defer db.Close()

			repo := storage.NewRentalRepository(&config.Config{}, db)

			tc.mockFunc(mock)

			rental, err := repo.GetByID(context.Background(), tc.idParam)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}

			if !cmp.Equal(rental, tc.expectedRental) {
				t.Fatalf("result expectation mismatch: %s", cmp.Diff(rental, tc.expectedRental))
			}

			if err != tc.expectedError && !errors.Is(err, tc.expectedError) {
				t.Fatalf("error expectation mismatch: %s", cmp.Diff(err, tc.expectedError))
			}
		})
	}
}

func TestRentalRepository_List(t *testing.T) {
	sq := "SELECT %s FROM rentals JOIN users ON users.id = rentals.user_id"

	testCases := []struct {
		name           string
		filters        *storage.RentalFilters
		expectedResult model.Rentals
		expectedError  error
		mockFunc       func(mock sqlmock.Sqlmock)
	}{
		{
			name:           "List without filters",
			expectedResult: _rentals,
			expectedError:  nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				selectQuery := fmt.Sprintf(sq, strings.Join(_columns, ", "))
				mock.ExpectQuery(selectQuery).
					WillReturnRows(sqlmock.NewRows(_columns).
						AddRow(rentalValues(_rentals[0])...).
						AddRow(rentalValues(_rentals[1])...))
			},
		},
		{
			name: "List with ID filter",
			filters: &storage.RentalFilters{
				IDs: []int32{1, 2},
			},
			expectedResult: model.Rentals{
				_rentals[1],
			},
			expectedError: nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				selectQuery := fmt.Sprintf(sq, strings.Join(_columns, ", "))
				mock.ExpectQuery(selectQuery + " WHERE rentals.id IN \\(1, 2\\)").
					// WithArgs(2).
					WillReturnRows(sqlmock.NewRows(_columns).
						AddRow(rentalValues(_rentals[1])...))
			},
		},
		{
			name: "List with price range filter",
			filters: &storage.RentalFilters{
				PriceMin: toPtr[int64](1500),
				PriceMax: toPtr[int64](2000),
			},
			expectedResult: model.Rentals{
				_rentals[1],
			},
			expectedError: nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				selectQuery := fmt.Sprintf(sq, strings.Join(_columns, ", "))
				mock.ExpectQuery(selectQuery + " WHERE price_per_day >= 1500 AND price_per_day <= 2000").
					WillReturnRows(sqlmock.NewRows(_columns).
						AddRow(rentalValues(_rentals[1])...))
			},
		},
		{
			name: "List with location filter",
			filters: &storage.RentalFilters{
				Near: &storage.Location{
					Latitude:  53.28,
					Longitude: -129.12,
				},
			},
			expectedResult: model.Rentals{
				_rentals[1],
			},
			expectedError: nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				subqueryColumns := strings.Join(_columns, ", ")
				subqueryColumns += ", ABS\\(lat - 53.28\\) as a"
				subqueryColumns += ", ABS\\(lng - -129.12\\) as b"

				parentQueryColumns := strings.Join(_columns, ", ")
				parentQueryColumns = strings.ReplaceAll(parentQueryColumns, "users.id as users_id", "subquery.users_id")
				parentQueryColumns = strings.ReplaceAll(parentQueryColumns, "rentals.", "subquery.")
				parentQueryColumns = strings.ReplaceAll(parentQueryColumns, "users.", "subquery.")

				selectQuery := fmt.Sprintf(sq, subqueryColumns) + " WHERE ABS\\(lat - 53.28\\) <= 100 AND ABS\\(lng - -129.12\\) <= 100"
				mock.ExpectQuery(
					fmt.Sprintf("SELECT %s FROM \\(%s\\) subquery WHERE SQRT\\(POW\\(a, 2\\) \\+ POW\\(b, 2\\)\\) <= 100", parentQueryColumns, selectQuery)).
					WillReturnRows(sqlmock.NewRows(_columns).
						AddRow(rentalValues(_rentals[1])...))
			},
		},
		{
			name: "List with order by filter",
			filters: &storage.RentalFilters{
				OrderBy: toPtr("price_per_day"),
			},
			expectedResult: _rentals,
			expectedError:  nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				selectQuery := fmt.Sprintf(sq, strings.Join(_columns, ", "))
				mock.ExpectQuery(selectQuery + " ORDER BY price_per_day").
					WillReturnRows(sqlmock.NewRows(_columns).
						AddRow(rentalValues(_rentals[0])...).
						AddRow(rentalValues(_rentals[1])...))
			},
		},
		{
			name: "List with limit filter",
			filters: &storage.RentalFilters{
				Pagination: storage.Pagination{
					Limit: toPtr(3),
				},
			},
			expectedResult: _rentals,
			expectedError:  nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				selectQuery := fmt.Sprintf(sq, strings.Join(_columns, ", "))
				mock.ExpectQuery(selectQuery + " LIMIT 3").
					WillReturnRows(sqlmock.NewRows(_columns).
						AddRow(rentalValues(_rentals[0])...).
						AddRow(rentalValues(_rentals[1])...))
			},
		},
		{
			name: "List with offset filter",
			filters: &storage.RentalFilters{
				Pagination: storage.Pagination{
					Offset: toPtr(8),
				},
			},
			expectedResult: _rentals,
			expectedError:  nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				selectQuery := fmt.Sprintf(sq, strings.Join(_columns, ", "))
				mock.ExpectQuery(selectQuery + " OFFSET 8").
					WillReturnRows(sqlmock.NewRows(_columns).
						AddRow(rentalValues(_rentals[0])...).
						AddRow(rentalValues(_rentals[1])...))
			},
		},
		{
			name:           "List with error",
			expectedResult: nil,
			expectedError:  sql.ErrConnDone,
			mockFunc: func(mock sqlmock.Sqlmock) {
				selectQuery := fmt.Sprintf(sq, strings.Join(_columns, ", "))
				mock.ExpectQuery(selectQuery).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create sqlmock: %s", err)
			}
			defer db.Close()

			repo := storage.NewRentalRepository(&config.Config{NearThresholdRadius: _nearThresholdRadius}, db)

			tc.mockFunc(mock)

			rentals, err := repo.List(context.Background(), tc.filters)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("there were unfulfilled expectations: %s", err)
			}

			if !cmp.Equal(rentals, tc.expectedResult) {
				t.Fatalf("result expectation mismatch: %s", cmp.Diff(rentals, tc.expectedResult))
			}

			if err != tc.expectedError && !errors.Is(err, tc.expectedError) {
				t.Fatalf("error expectation mismatch: %s", cmp.Diff(err.Error(), tc.expectedError.Error()))
			}
		})
	}
}

// Helper function to create a pointers to values.
func toPtr[T any](v T) *T {
	return &v
}

func rentalValues(rental *model.Rental) []driver.Value {
	return []driver.Value{
		rental.ID,
		rental.UserID,
		rental.Name,
		rental.Type,
		rental.Description,
		rental.Sleeps,
		rental.PricePerDay,
		rental.HomeCity,
		rental.HomeState,
		rental.HomeZip,
		rental.HomeCountry,
		rental.VehicleMake,
		rental.VehicleModel,
		rental.VehicleYear,
		rental.VehicleLength,
		rental.Latitude,
		rental.Longitude,
		rental.PrimaryImageURL,
		rental.User.ID,
		rental.User.FirstName,
		rental.User.LastName,
	}
}
