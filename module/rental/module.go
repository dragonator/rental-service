package rental

import (
	"context"
	"fmt"

	"github.com/dragonator/rental-service/module/rental/internal/db"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
	"github.com/dragonator/rental-service/pkg/config"
	"github.com/dragonator/rental-service/pkg/logger"
)

type RentalModule struct {
	Storage *storage.RentalRepository
}

func NewRentalModule(config *config.Config, logger *logger.Logger) (*RentalModule, error) {
	db, err := db.OpenPGX(config, logger.Desugar())
	if err != nil {
		return nil, fmt.Errorf("creating rental model: %w", err)
	}

	rm := &RentalModule{
		Storage: storage.NewRentalRepository(config, db),
	}

	rentals, err := rm.Storage.List(context.Background(), &storage.RentalFilters{
		Pagination: storage.Pagination{
			// Limit: toPtr(3),
			// Offset: toPtr(3),
		},
		// PriceMin: toPtr(int64(9000)),
		// PriceMax: toPtr(int64(13000)),
		// OrderBy:  toPtr("price_per_day"),
		// IDs:      []string{"24", "26"},
		Near: &storage.Location{
			Latitude:  33.64,
			Longitude: -117.93,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(len(rentals))
	fmt.Println(rentals)
	for _, r := range rentals {
		fmt.Println(r.ID)
	}

	return &RentalModule{
		Storage: storage.NewRentalRepository(config, db),
	}, nil
}

func toPtr[T any](v T) *T {
	return &v
}
