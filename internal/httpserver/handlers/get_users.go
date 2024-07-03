package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
	"strconv"
)

func GetUsersHandler(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Get users")

		var users []models.User
		db := storage.DB

		firstName := r.URL.Query().Get("first_name")
		if firstName != "" {
			db = db.Where("first_name = ?", firstName)
		}

		lastName := r.URL.Query().Get("last_name")
		if lastName != "" {
			db = db.Where("last_name = ?", lastName)
		}

		email := r.URL.Query().Get("email")
		if email != "" {
			db = db.Where("email = ?", email)
		}

		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		if page != "" && limit != "" {
			pageInt, err := strconv.Atoi(page)
			if err != nil {
				log.Debug("Error converting page to int")
				log.Info("Error converting page to int")
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, render.M{"error": "Unable to parse page to int"})
				return
			}
			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				log.Debug("Error converting limit to int")
				log.Info("Error converting limit to int")
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, render.M{"error": "Unable to parse limit to int"})
				return
			}
			offset := (pageInt - 1) * limitInt
			db = db.Offset(offset).Limit(limitInt)
		}

		db.Find(&users)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, users)
		log.Debug("Successfully got users")
		log.Info("Successfully got users")
	}
}
