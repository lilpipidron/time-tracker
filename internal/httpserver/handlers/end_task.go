package handlers

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/lilpipidron/time-tracker/internal/httpserver/requests"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
	"time"
)

func EndTaskHandler(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("handling end task request")

		var endTaskRequest requests.EndTaskRequest
		var req interface{} = &endTaskRequest
		requests.Decode(w, r, &req)

		log.Debug("Decoded body", "end task handler")

		if err := validator.New().Struct(endTaskRequest); err != nil {
			var validationError validator.ValidationErrors
			errors.As(err, &validationError)

			render.Status(r, http.StatusBadRequest)
			log.Info("Invalid request", "error", validationError)
			render.JSON(w, r, validationError)
			return
		}

		userTask := models.UserTask{}
		storage.DB.First(&userTask, userTask.UserID, userTask.TaskID)
		if userTask.EndTime.IsZero() || !userTask.StartTime.IsZero() {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, nil)
			log.Info("Invalid request", "error", "Invalid request")
			return
		}

		userTask.EndTime = time.Now()
		storage.DB.Save(&userTask)
		log.Debug("Saved userTask", "userTask", userTask)
	}
}
