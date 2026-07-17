package main

import (
	"log"
	"net/http"
	"wallets-api-postgres/internal/config"
	"wallets-api-postgres/internal/router"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	appRouter := router.New()
	address := ":" + cfg.Server.Port

	log.Println("server started on port", cfg.Server.Port)

	if err := http.ListenAndServe(address, appRouter); err != nil {
		log.Fatal(err)
	}
}
