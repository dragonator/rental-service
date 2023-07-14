package main

import (
	"github.com/dragonator/rental-service/module/rental"
	"github.com/dragonator/rental-service/pkg/config"
	"github.com/dragonator/rental-service/pkg/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logger := logger.NewLogger(cfg.LoggerLevel)

	_, err = rental.NewRentalModule(cfg, logger)
	if err != nil {
		panic(err)
	}

	// a, err := rentalModule.Storage.GetByID(context.Background(), 1)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(a)
	// fmt.Println(a.User)
}
