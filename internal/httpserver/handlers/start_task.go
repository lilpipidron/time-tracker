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

// StartTaskHandler handles starting a task for a user
//
//	@Summary		Start a task
//	@Description	Start a new task for a user
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Param			startTaskRequest	body		requests.StartTaskRequest	true	"Start Task Request"
//	@Success		201					{string}	string						"Successfully started task"
//	@Failure		400					{string}	string						"Invalid request"
//	@Router			/task/start [post]
func StartTaskHandler(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("handling start task request")

		var startTaskRequest requests.StartTaskRequest
		var req interface{} = &startTaskRequest
		requests.Decode(w, r, &req)

		log.Debug("received start task request")

		if err := validator.New().Struct(req); err != nil {
			var validationError validator.ValidationErrors
			errors.As(err, &validationError)

			render.Status(r, http.StatusBadRequest)

			log.Info("Invalid request", err)

			render.JSON(w, r, validationError.Error())

			return
		}

		userTask := models.UserTask{
			UserID:    startTaskRequest.UserID,
			TaskID:    startTaskRequest.TaskID,
			StartTime: time.Now(),
		}

		storage.DB.Create(&userTask)

		log.Info("Successfully started task")
		render.Status(r, http.StatusCreated)
		log.Info("Successfully started task")
	}
}
