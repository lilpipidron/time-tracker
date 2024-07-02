package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/lilpipidron/time-tracker/internal/config"
	"github.com/lilpipidron/time-tracker/internal/httpserver/requests"
	"github.com/lilpipidron/time-tracker/internal/models"
	"github.com/lilpipidron/time-tracker/internal/storage/postgresql"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type RequestToApi struct {
	PassportSerie  int `json:"passportSerie"`
	PassportNumber int `json:"passportNumber"`
}

type ResponseFromApi struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}

func AddUser(storage *postgresql.Storage, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Add user")

		var addUserRequest requests.AddUserRequest
		var req interface{} = &addUserRequest
		requests.Decode(w, r, &req)

		log.Debug("Decoded body", "add user request", addUserRequest)

		if err := validator.New().Struct(addUserRequest); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			render.Status(r, http.StatusBadRequest)

			log.Info("invalid request", err)

			render.JSON(w, r, validateErr.Error())

			return
		}

		passportData := strings.Split(addUserRequest.PassportNumber, " ")
		if passportData == nil || len(passportData) < 2 {
			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, "passport is required")

			log.Info("passport is required")

			return
		}

		series, err := strconv.Atoi(passportData[0])
		if err != nil {
			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, fmt.Sprintf("passport series is required: %s", err.Error()))

			log.Info(fmt.Sprintf("passport series is required: %s", err.Error()))

			return
		}

		number, err := strconv.Atoi(passportData[1])
		if err != nil {
			render.Status(r, http.StatusBadRequest)

			render.JSON(w, r, fmt.Sprintf("passport number is required: %s", err.Error()))

			log.Info(fmt.Sprintf("passport number is required: %s", err.Error()))

			return
		}

		requestToApi := RequestToApi{
			PassportSerie:  series,
			PassportNumber: number,
		}

		response, err := sendRequestToApi(requestToApi, cfg)
		if errors.Is(err, fmt.Errorf("bad request")) {
			render.Status(r, http.StatusBadRequest)
			return
		} else if errors.Is(err, fmt.Errorf("internal server error")) {
			render.Status(r, http.StatusInternalServerError)
			return
		}

		log.Debug("Response from API:", response)
		render.JSON(w, r, response)

		user := models.User{
			Name:           response.Name,
			Surname:        response.Surname,
			Patronymic:     response.Patronymic,
			Address:        response.Address,
			PassportNumber: addUserRequest.PassportNumber,
		}

		if err = storage.DB.Create(&user).Error; err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, fmt.Sprintf("error saving user: %s", err.Error()))
			log.Debug("error saving user:", err)
		}
	}
}

func sendRequestToApi(request RequestToApi, cfg config.Config) (*ResponseFromApi, error) {
	log.Debug("Send request to API")

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Debug("Failed marshalling request to API:", "error", err)
		return nil, fmt.Errorf("error marshalling request body: %s", err.Error())
	}

	client := http.Client{}

	req, err := http.NewRequest("GET", cfg.ApiUrl+"/info", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Debug("Failed creating request to API:", "error", err)
		return nil, fmt.Errorf("error creating HTTP request: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed sending request to API:", "error", err)
		return nil, fmt.Errorf("error sending HTTP request: %s", err.Error())
	}

	if resp.StatusCode == http.StatusBadRequest {
		log.Debug("Get status bad request")
		return nil, fmt.Errorf("bad request")
	} else if resp.StatusCode == http.StatusInternalServerError {
		log.Debug("Get status internal server error")
		return nil, fmt.Errorf("internal server error")
	}

	var responseFromApi ResponseFromApi
	err = render.DecodeJSON(resp.Body, &responseFromApi)
	if errors.Is(err, io.EOF) {
		log.Debug("request body is empty")

		return nil, fmt.Errorf("request body is empty: %s", err.Error())
	}

	if err != nil {
		log.Error("failed to decode request body", err)
		return nil, fmt.Errorf("failed to decode request body: %s", err.Error())
	}

	if err := validator.New().Struct(responseFromApi); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		log.Debug("invalid response", "response", responseFromApi)

		return nil, fmt.Errorf("invalid response: %s", err.Error())
	}

	return &responseFromApi, nil
}
