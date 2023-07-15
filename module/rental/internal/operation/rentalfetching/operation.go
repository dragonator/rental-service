package rentalfetching

import (
	"context"
	"fmt"

	"github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
)

// Operation provides an API for fetching single or multiple rentals.
type Operation struct {
	rentalStore RentalStore
}

// NewOperation is a contruction function for Operation.
func NewOperation(rentalStore RentalStore) *Operation {
	return &Operation{
		rentalStore: rentalStore,
	}
}

// GetRentalByID returns a rental for the given id.
func (o *Operation) GetRentalByID(ctx context.Context, rentalID int) (*model.Rental, error) {
	rental, err := o.rentalStore.GetByID(ctx, rentalID)
	if err != nil {
		return nil, fmt.Errorf("operation GetRentalByID: %w", err)
	}

	return rental, nil
}

// ListRentals returns a list of rentals based on the specified filters.
// If no rentals are found it returns an empty list.
func (o *Operation) ListRentals(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error) {
	rentals, err := o.rentalStore.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("operation ListRentals: %w", err)
	}

	return rentals, nil
}
