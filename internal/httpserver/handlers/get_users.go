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
		log.Debug("handling get users")

		var users []models.User
		db := storage.DB

		firstName := r.URL.Query().Get("name")
		if firstName != "" {
			db = db.Where("name = ?", firstName)
		}

		lastName := r.URL.Query().Get("surname")
		if lastName != "" {
			db = db.Where("surname = ?", lastName)
		}

		email := r.URL.Query().Get("patronymic")
		if email != "" {
			db = db.Where("patronymic = ?", email)
		}

		address := r.URL.Query().Get("address")
		if address != "" {
			db = db.Where("address = ?", address)
		}

		passportNumber := r.URL.Query().Get("passportNumber")
		if passportNumber != "" {
			db = db.Where("passport_number = ?", passportNumber)
		}

		page := r.URL.Query().Get("page")
		limit := r.URL.Query().Get("limit")

		if page != "" && limit != "" {
			pageInt, err := strconv.Atoi(page)
			if err != nil {
				log.Error("Error converting page to int")
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, err.Error())
				return
			}
			limitInt, err := strconv.Atoi(limit)
			if err != nil {
				log.Error("Error converting limit to int")
				render.Status(r, http.StatusInternalServerError)
				render.JSON(w, r, err.Error())
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
