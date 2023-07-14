package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	rentalModule, err := rental.NewRentalModule(cfg, logger)
	if err != nil {
		panic(err)
	}

	rentalModule.RentalService.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sig := <-stop

	log.Printf("Signal caught (%s), stopping...", sig.String())
	rentalModule.RentalService.Stop()
	log.Print("Service stopped.")
}
