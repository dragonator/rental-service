package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	model "github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/pkg/config"
)

var (
	rentalColums = []string{
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
		"rentals.created",
		"rentals.updated",
	}
	userColumns = []string{
		"users.id as users_id",
		"users.first_name",
		"users.last_name",
	}
)

type rowScanner interface {
	Scan(dest ...any) error
}

// RentalRepository hold DB operations over rental entities.
type RentalRepository struct {
	db                  *sql.DB
	nearThresholdRadius int
}

// NewRentalRepository is a constructor function for RentalRepository.
func NewRentalRepository(config *config.Config, db *sql.DB) *RentalRepository {
	return &RentalRepository{
		db:                  db,
		nearThresholdRadius: config.NearThresholdRadius,
	}
}

// GetByID returns a single rental object corresponding to the requested id.
// If no such rental exists it returns an error.
func (rr *RentalRepository) GetByID(ctx context.Context, rentalID string) (*model.Rental, error) {
	qb := NewQueryBuilder().
		Select().
		Columns(rentalColums...).
		Columns(userColumns...).
		From("rentals").
		Join("users ON users.id = rentals.user_id").
		Where("rentals.id = $1")

	rental, err := scanRental(rr.db.QueryRowContext(ctx, qb.String(), rentalID))
	if err != nil {
		return nil, fmt.Errorf("getting rental by id: %w", err)
	}

	return rental, nil
}

// List returns a list of rentals based on the given filters. If no results are found it returns an empty list.
func (rr *RentalRepository) List(ctx context.Context, filters *RentalFilters) (model.Rentals, error) {
	rentals := make(model.Rentals, 0, 10)

	sqlQuery := rr.buildListQuery(filters)

	rows, err := rr.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		rental, err := scanRental(rows)
		if err != nil {
			return nil, fmt.Errorf("scanning rental: %w", err)
		}

		rentals = append(rentals, rental)
	}

	return rentals, nil
}

func (rr *RentalRepository) buildListQuery(f *RentalFilters) string {
	qb := NewQueryBuilder().
		Select().
		Columns(rentalColums...).
		Columns(userColumns...).
		From("rentals").
		Join("users ON users.id = rentals.user_id")

	if len(f.IDs) > 0 {
		qb.Where(fmt.Sprintf("rentals.id IN (%s)", strings.Join(f.IDs, ",")))
	}

	if f.PriceMin != nil {
		qb.Where(fmt.Sprintf("price_per_day >= %d", *f.PriceMin))
	}

	if f.PriceMax != nil {
		qb.Where(fmt.Sprintf("price_per_day <= %d", *f.PriceMax))
	}

	if f.Near != nil {
		qb.Columns(
			fmt.Sprintf("ABS(lat - %f) as a", f.Near.Latitude),
			fmt.Sprintf("ABS(lng - %f) as b", f.Near.Longitude),
		)
		qb.Where(fmt.Sprintf("ABS(lat - %f) <= %d", f.Near.Latitude, rr.nearThresholdRadius))
		qb.Where(fmt.Sprintf("ABS(lng - %f) <= %d", f.Near.Longitude, rr.nearThresholdRadius))
	}

	if f.Near != nil {
		newRentalColumns := changeColumnTable("rentals.", "subquery.", rentalColums...)

		qbTmp := NewQueryBuilder().
			Select().
			Columns(newRentalColumns...).
			Columns("subquery.users_id", "subquery.first_name", "subquery.last_name").
			From(fmt.Sprintf("(%s) subquery", qb.String())).
			Where(fmt.Sprintf("SQRT(a*a + b*b) <= %d", rr.nearThresholdRadius))

		qb = qbTmp
	}

	if f.OrderBy != nil {
		qb.OrderBy(*f.OrderBy)
	}

	if f.Limit != nil {
		qb.Limit(*f.Limit)
	}

	if f.Offset != nil {
		qb.Offset(*f.Offset)
	}

	return qb.String()
}

func changeColumnTable(oldPrefix string, newPrefix string, columns ...string) []string {
	newColumns := make([]string, 0, len(columns))

	for _, c := range columns {
		newColumns = append(newColumns, strings.Replace(c, oldPrefix, newPrefix, 1))
	}

	return newColumns
}

func scanRental(row rowScanner) (*model.Rental, error) {
	rental := new(model.Rental)
	rental.User = new(model.User)

	if err := row.Scan(
		&rental.ID,
		&rental.UserID,
		&rental.Name,
		&rental.Type,
		&rental.Description,
		&rental.Sleeps,
		&rental.PricePerDay,
		&rental.HomeCity,
		&rental.HomeState,
		&rental.HomeZip,
		&rental.HomeCountry,
		&rental.VehicleMake,
		&rental.VehicleModel,
		&rental.VehicleYear,
		&rental.VehicleLength,
		&rental.Latitude,
		&rental.Longitude,
		&rental.PrimaryImageURL,
		&rental.Created,
		&rental.Updated,
		&rental.User.ID,
		&rental.User.FirstName,
		&rental.User.LastName,
	); err != nil {
		return nil, err
	}

	return rental, nil
}
