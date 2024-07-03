package handlers

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/lilpipidron/time-tracker/internal/httpserver/requests"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func EndTaskHandler(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("handling end task request")

		var endTaskRequest requests.EndTaskRequest
		var req interface{} = &endTaskRequest
		requests.Decode(w, r, &req)

		log.Debug("Decoded body", "body", endTaskRequest)

		if err := validator.New().Struct(endTaskRequest); err != nil {
			var validationError validator.ValidationErrors
			errors.As(err, &validationError)

			render.Status(r, http.StatusBadRequest)
			log.Error("Invalid request", "error", validationError)
			render.JSON(w, r, validationError.Error())
			return
		}

		userTask := models.UserTask{}
		err := storage.DB.Where("user_id = ? AND task_id = ?", endTaskRequest.UserID, endTaskRequest.TaskID).First(&userTask).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, nil)
				log.Info("Record not found", "error", err)
				return
			}
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			log.Error("Failed to query user task", "error", err)
			return
		}

		if !userTask.EndTime.IsZero() || userTask.StartTime.IsZero() {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "Invalid request")
			log.Error("Invalid request")
			return
		}

		userTask.EndTime = time.Now()
		storage.DB.Save(&userTask)
		log.Debug("Saved userTask", "userTask", userTask)
	}
}
