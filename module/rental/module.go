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

type RentalModule struct {
	RentalService *service.Service
}

func NewRentalModule(config *config.Config, logger *logger.Logger) (*RentalModule, error) {
	db, err := db.OpenPGX(config, logger.Desugar())
	if err != nil {
		return nil, fmt.Errorf("creating rental model: %w", err)
	}

	rentalStore := storage.NewRentalRepository(config, db)
	rentalFetchingOp := rentalfetching.NewOperation(rentalStore)
	rentalHandler := handler.NewRentalHandler(rentalFetchingOp)
	router := service.NewRouter(rentalHandler)

	rentalService, err := service.New(config, logger, router)
	if err != nil {
		return nil, fmt.Errorf("creating rental model: %w", err)
	}

	return &RentalModule{
		RentalService: rentalService,
	}, nil
}

func toPtr[T any](v T) *T {
	return &v
}
