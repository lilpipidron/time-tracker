package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
)

func DeleteUser(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("DeleteUser")
		defer r.Body.Close()

		userID := chi.URLParam(r, "userID")
		var user models.User
		if err := storage.DB.First(&user, userID).Error; err != nil {
			log.Info("Failed to find user:", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if err := storage.DB.Delete(&user).Error; err != nil {
			log.Info("Failed to delete user:", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		log.Debug("User deleted")
		log.Info("User deleted")
	}
}
