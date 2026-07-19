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
	userService := service.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	appRouter := router.New(userHandler)

	address := ":" + cfg.Server.Port

	log.Println("server started on port", cfg.Server.Port)

	if err := http.ListenAndServe(address, appRouter); err != nil {
		log.Fatal(err)
	}
}
