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
	cfg := config.MustLoad()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", cfg.PostgresHost, cfg.PostgresPort,
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	_, err := postgresql.NewPostgresDB(psqlInfo, "postgres")
	if err != nil {
		log.Fatal(err)
	}
}
