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
	_rental              = &model.Rental{
		ID:              1,
		UserID:          2,
		Name:            "Rental",
		Type:            "Type",
		Description:     "Description",
		Sleeps:          42,
		PricePerDay:     64,
		HomeCity:        "HomeCity",
		HomeState:       "HomeState",
		HomeZip:         "HomeZip",
		HomeCountry:     "HomeCountry",
		VehicleMake:     "VehicleMake",
		VehicleModel:    "VehicleModel",
		VehicleYear:     2023,
		VehicleLength:   13.5,
		Latitude:        53.28,
		Longitude:       -129.12,
		PrimaryImageURL: "PrimaryImageURL",
		User: &model.User{
			ID:        2,
			FirstName: "FirstName",
			LastName:  "LastName",
		},
	}
	_rentalValues = []driver.Value{
		_rental.ID,
		_rental.UserID,
		_rental.Name,
		_rental.Type,
		_rental.Description,
		_rental.Sleeps,
		_rental.PricePerDay,
		_rental.HomeCity,
		_rental.HomeState,
		_rental.HomeZip,
		_rental.HomeCountry,
		_rental.VehicleMake,
		_rental.VehicleModel,
		_rental.VehicleYear,
		_rental.VehicleLength,
		_rental.Latitude,
		_rental.Longitude,
		_rental.PrimaryImageURL,
		_rental.User.ID,
		_rental.User.FirstName,
		_rental.User.LastName,
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
			expectedRental: _rental,
			expectedError:  nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(selectQuery).
					WithArgs([]driver.Value{1}...).
					WillReturnRows(sqlmock.NewRows(_columns).AddRow(_rentalValues...))
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
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

			if !cmp.Equal(rental, tc.expectedRental) {
				t.Errorf("result expectation mismatch: %s", cmp.Diff(rental, tc.expectedRental))
			}

			if err != tc.expectedError && !errors.Is(err, tc.expectedError) {
				t.Errorf("error expectation mismatch: %s", cmp.Diff(err, tc.expectedError))
			}
		})
	}
}
