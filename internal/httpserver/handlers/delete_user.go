package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
)

// DeleteUserHandler handles deleting a user by ID
//
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Param			userID	path	int	true	"User ID"
//	@Success		204		"No Content"
//	@Failure		404		{object}	map[string]string	"Not Found"
//	@Failure		500		{object}	map[string]string	"Internal Server Error"
//	@Router			/user/{userID} [delete]
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

		render.Status(r, http.StatusNoContent)

		log.Debug("User deleted")
	}
}
