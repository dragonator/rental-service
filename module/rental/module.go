package rental

import (
	"fmt"

	"github.com/dragonator/rental-service/module/rental/internal/db"
	"github.com/dragonator/rental-service/module/rental/internal/http/handler"
	"github.com/dragonator/rental-service/module/rental/internal/http/service"
	"github.com/dragonator/rental-service/module/rental/internal/operation/rentalfetching"
	"github.com/dragonator/rental-service/module/rental/internal/storage"
	"github.com/dragonator/rental-service/pkg/config"
	"github.com/dragonator/rental-service/pkg/logger"
)

// RentalService provides methods for starting and stopping a rental service.
type RentalService interface {
	Start()
	Stop()
}

// RentalModule provides access to the functionality of rental module.
type RentalModule struct {
	RentalService RentalService
}

// NewRentalModule is a construction function for RentalModule.
func NewRentalModule(config *config.Config, logger *logger.Logger) (*RentalModule, error) {
	db, err := db.OpenPGX(config, logger.Desugar())
	if err != nil {
		return nil, fmt.Errorf("creating rental module: %w", err)
	}

	rentalStore := storage.NewRentalRepository(config, db)
	rentalFetchingOp := rentalfetching.NewOperation(rentalStore)
	rentalHandler := handler.NewRentalHandler(rentalFetchingOp)
	router := service.NewRouter(rentalHandler)

	rentalService, err := service.New(config, logger, router)
	if err != nil {
		return nil, fmt.Errorf("creating rental module: %w", err)
	}

	return &RentalModule{
		RentalService: rentalService,
	}, nil
}
