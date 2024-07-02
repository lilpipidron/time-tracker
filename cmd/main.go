package main

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/lilpipidron/time-tracker/internal/config"
	"github.com/lilpipidron/time-tracker/internal/httpserver/handlers"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
	"strconv"
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
	storage, err := postgresql.NewPostgresDB(dsn, cfg.PostgresDB)

	if err != nil {
		log.Debug("Failed to connect to database: %v", err)
	}
	log.Debug("Connected to PostgreSQL")

	log.Info("Connected to database successfully")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/user/{userID}", nil)

	router.Get("/user/{userID}/worklogs", nil)

	router.Post("/task/start", nil)

	router.Post("/tasks/stop", nil)

	router.Delete("/user/{userID}", handlers.DeleteUser(storage))

	router.Put("/user", nil)

	router.Post("/user", handlers.AddUser(storage, *cfg))

	addr := cfg.ServiceHost + ":" + strconv.Itoa(cfg.ServicePort)

	if err = http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}
