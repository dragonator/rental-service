package rentalfetching

import (
	"context"

	"github.com/dragonator/rental-service/module/rental/internal/model"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
)

// RentalStore is a contract to a rental storage.
//
//go:generate moq -rm -pkg rentalfetching_test -out rental_store_mock_test.go . RentalStore
type RentalStore interface {
	GetByID(ctx context.Context, rentalID int) (*model.Rental, error)
	List(ctx context.Context, filters *storage.RentalFilters) (model.Rentals, error)
}
