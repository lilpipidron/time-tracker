package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
	"strconv"
)

func GetUsers(storage *postgresql.Storage) http.HandlerFunc {
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
			pageInt, _ := strconv.Atoi(page)
			limitInt, _ := strconv.Atoi(limit)
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
