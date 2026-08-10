package main

import (
	"log/slog"
	"net/http"
	"os"
	"wallets-api-postgres/internal/config"
	"wallets-api-postgres/internal/database"
	"wallets-api-postgres/internal/handlers"
	"wallets-api-postgres/internal/repository"
	"wallets-api-postgres/internal/router"
	"wallets-api-postgres/internal/service"
	_ "wallets-api-postgres/docs"
)

// @title Wallets API
// @version 1.0
// @description REST API for users, wallets, and transfers.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error(
			"failed to load config",
			"error", err,
		)
		os.Exit(1)
	}

	db, err := database.Connect(cfg.Database)
	if err != nil {
		slog.Error(
			"failed to connect to database",
			"error", err,
		)
		os.Exit(1)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository, cfg.JWT.Secret)
	userHandler := handlers.NewUserHandler(userService)

	walletRepository := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepository)
	walletHandler := handlers.NewWalletHandler(walletService)

	transferRepository := repository.NewTransferRepository(db)
	transferService := service.NewTransferService(transferRepository)
	transferHandler := handlers.NewTransferHandler(transferService)

	appRouter := router.New(
		userHandler,
		walletHandler,
		transferHandler,
		cfg.JWT.Secret,
	)

	address := ":" + cfg.Server.Port

	slog.Info(
		"server started",
		"port", cfg.Server.Port,
	)

	if err := http.ListenAndServe(address, appRouter); err != nil {
		slog.Error(
			"server stopped",
			"error", err,
		)
		os.Exit(1)
	}
}
