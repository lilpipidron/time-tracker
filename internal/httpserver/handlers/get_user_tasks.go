package handlers

import (
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"net/http"
	"time"
)

type TaskTime struct {
	TaskID       int     `json:"task_id"`
	TotalHours   float64 `json:"total_hours"`
	TotalMinutes int64   `json:"total_minutes"`
}

func GetUserTasksHandler(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("handling get user tasks request")

		userID := chi.URLParam(r, "userID")
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			log.Error("Invalid start date", "error", err)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, err.Error())
			log.Error("Invalid end date", "error", err)
			return
		}

		var taskTimes []TaskTime

		if err = storage.DB.Table("user_tasks").
			Select("task_id, FLOOR(EXTRACT(EPOCH FROM SUM(end_time - start_time)) / 3600) AS total_hours, FLOOR((EXTRACT(EPOCH FROM SUM(end_time - start_time)) / 60) % 60) AS total_minutes").
			Where("user_id = ? AND start_time >= ? AND end_time <= ?", userID, startDate, endDate).
			Group("task_id").
			Order("total_hours DESC, total_minutes DESC").
			Scan(&taskTimes).Error; err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			log.Error("Failed to fetch task times", "error", err)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, taskTimes)
		log.Debug("finished get user tasks request")
	}
}
