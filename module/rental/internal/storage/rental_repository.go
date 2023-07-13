package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dragonator/rental-service/module/rental/internal/models"
)

// RentalRepository hold DB operations over rental entities.
type RentalRepository struct {
	db *sql.DB
}

// NewRentalRepository is a constructor function for RentalRepository.
func NewRentalRepository(db *sql.DB) *RentalRepository {
	return &RentalRepository{
		db: db,
	}
}

// GetByID returns a single rental object corresponding to the requested id.
// If no such rental exists it returns an error.
func (rr *RentalRepository) GetByID(ctx context.Context, rentalID int) (*models.Rental, error) {
	result := new(models.Rental)
	result.User = new(models.User)

	sqlStatement := `
		SELECT * 
		FROM rentals
		JOIN users ON users.id = rentals.user_id
		WHERE rentals.id = $1;
	`

	err := rr.db.QueryRowContext(ctx, sqlStatement, rentalID).
		Scan(
			&result.ID,
			&result.UserID,
			&result.Name,
			&result.Type,
			&result.Description,
			&result.Sleeps,
			&result.PricePerDay,
			&result.HomeCity,
			&result.HomeState,
			&result.HomeZip,
			&result.HomeCountry,
			&result.VehicleMake,
			&result.VehicleModel,
			&result.VehicleYear,
			&result.VehicleLength,
			&result.Created,
			&result.Updated,
			&result.Latitude,
			&result.Longitude,
			&result.PrimaryImageURL,
			&result.User.ID,
			&result.User.FirstName,
			&result.User.LastName,
		)
	if err != nil {
		return nil, fmt.Errorf("getting rental by id: %w", err)
	}

	return result, nil
}

// List returns a list of rentals based on the given filters. If no results are found it returns an empty list.
func (rr *RentalRepository) List(ctx context.Context, filters *RentalFilters) (models.Rentals, error) {
	return nil, nil
}
