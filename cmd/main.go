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

//	@title			Time Tracker API
//	@version		1.0
//	@description	This is a sample server for tracking tasks and user activities.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/
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

	router.Get("/users", handlers.GetUsersHandler(storage))

	router.Get("/user/{userID}/tasks", handlers.GetUserTasksHandler(storage))

	router.Post("/task/start", handlers.StartTaskHandler(storage))

	router.Put("/task/end", handlers.EndTaskHandler(storage))

	router.Delete("/user/{userID}", handlers.DeleteUserHandler(storage))

	router.Put("/user", handlers.ChangeUserInfoHandler(storage))

	router.Post("/user", handlers.AddUserHandler(storage, *cfg))

	addr := cfg.ServiceHost + ":" + strconv.Itoa(cfg.ServicePort)

	if err = http.ListenAndServe(addr, router); err != nil {
		panic(err)
	}
}
