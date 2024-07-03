package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
)

func DeleteUserHandler(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("handling delete request")

		userID := chi.URLParam(r, "userID")
		var user models.User
		if err := storage.DB.First(&user, userID).Error; err != nil {
			log.Error("Failed to find user:", err)
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, err.Error())
			return
		}

		if err := storage.DB.Delete(&user).Error; err != nil {
			log.Error("Failed to delete user:", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			return
		}

		log.Debug("User deleted")
	}
}
