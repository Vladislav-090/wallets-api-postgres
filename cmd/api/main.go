package main

import (
	"log"
	"net/http"
	"wallets-api-postgres/internal/config"
	"wallets-api-postgres/internal/database"
	"wallets-api-postgres/internal/handlers"
	"wallets-api-postgres/internal/repository"
	"wallets-api-postgres/internal/router"
	"wallets-api-postgres/internal/service"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, cfg.JWT.Secret)
	userHandler := handlers.NewUserHandler(userService)

	walletRepository := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepository)
	walletHandler := handlers.NewWalletHandler(walletService)

	appRouter := router.New(
		userHandler,
		walletHandler,
		cfg.JWT.Secret,
	)

	address := ":" + cfg.Server.Port

	log.Println("server started on port", cfg.Server.Port)

	if err := http.ListenAndServe(address, appRouter); err != nil {
		log.Fatal(err)
	}
}
