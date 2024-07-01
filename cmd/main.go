package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/lilpipidron/time-tracker/internal/config"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	log.Debug("Loading configuration...")
	cfg := config.MustLoad()
	log.Debug("Configuration loaded")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	log.Debug("Connecting to PostgreSQL...")
	_, err := postgresql.NewPostgresDB(dsn, cfg.PostgresDB)
	if err != nil {
		log.Debug("Failed to connect to database: %v", err)
	}
	log.Debug("Connected to PostgreSQL")

	log.Info("Connected to database successfully")
}
