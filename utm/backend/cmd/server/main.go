package main

import (
	"log"
	"os"

	"github.com/rs/zerolog"
	"github.com/utm/backend/internal/config"
	"github.com/utm/backend/internal/database"
	"github.com/utm/backend/internal/hapi"
	scrpclient "github.com/utm/backend/internal/scrp"
)

func main() {
	// Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if cfg.App.Env == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Database connections
	pgPool, err := database.NewPostgresPool(cfg.Postgres)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer pgPool.Close()

	mongoDB, err := database.NewMongoClient(cfg.Mongo)
	if err != nil {
		log.Fatalf("connect mongo: %v", err)
	}

	// SCRP gRPC client
	scrpClient, err := scrpclient.NewClient(cfg.SCRP.Address)
	if err != nil {
		log.Fatalf("connect scrp: %v", err)
	}
	defer scrpClient.Close()

	// HTTP server
	r := hapi.NewRouter(cfg, pgPool, mongoDB, scrpClient)

	port := cfg.App.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("UTM HTTP server starting on :%s (env=%s)", port, cfg.App.Env)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
		os.Exit(1)
	}
}
