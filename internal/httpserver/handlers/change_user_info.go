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
	"reflect"
)

func ChangeUserInfoHandler(storage *postgresql.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Handling change user info request")

		var changeUserInfo requests.ChangeUserInfo
		var req interface{} = &changeUserInfo
		requests.Decode(w, r, &req)

		log.Debug("Decoded body", "body", changeUserInfo)

		if err := validator.New().Struct(changeUserInfo); err != nil {
			var validationError validator.ValidationErrors
			errors.As(err, &validationError)

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, validationError)
			log.Error("Invalid request", "error", changeUserInfo)
			return
		}

		user := models.User{}
		if err := storage.DB.Where("passport_number = ?", changeUserInfo.PassportNumber).First(&user).Error; err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, nil)
			log.Info("User not found", "error", err)
			return
		}

		userVal := reflect.ValueOf(&user).Elem()
		changeUserInfoVal := reflect.ValueOf(&changeUserInfo).Elem()
		for i := 0; i < changeUserInfoVal.NumField(); i++ {
			field := changeUserInfoVal.Field(i)
			fieldName := changeUserInfoVal.Type().Field(i).Name

			if field.Kind() == reflect.String && field.String() != "" {
				userField := userVal.FieldByName(fieldName)

				if userField.IsValid() && userField.CanSet() {
					userField.SetString(field.String())
				}
			}
		}

		if err := storage.DB.Save(&user).Error; err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, err.Error())
			log.Error("Failed to update user", "error", err)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, user)
		log.Debug("User updated successfully", "user", user)

	}
}
