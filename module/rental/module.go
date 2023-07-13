package rental

import (
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

	return &RentalModule{
		Storage: storage.NewRentalRepository(db),
	}, nil
}
