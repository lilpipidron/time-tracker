package handlers

import (
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

func AddUserHandler(storage *postgresql.Storage, cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Handling request to add user")

		var addUserRequest requests.AddUserRequest
		var req interface{} = &addUserRequest
		requests.Decode(w, r, &req)

		log.Debug("Decoded body", "add user request", addUserRequest)

		if err := validator.New().Struct(addUserRequest); err != nil {
			var validateError validator.ValidationErrors
			errors.As(err, &validateError)

			render.Status(r, http.StatusBadRequest)
			log.Info("Invalid request", "error", err)
			render.JSON(w, r, validateError.Error())
			return
		}

		passportData := strings.Split(addUserRequest.PassportNumber, " ")
		if len(passportData) < 2 {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "passport is required")
			log.Info("Passport is required")
			return
		}

		series, err := strconv.Atoi(passportData[0])
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, fmt.Sprintf("passport series is required: %s", err.Error()))
			log.Info("Passport series is required", "error", err)
			return
		}

		number, err := strconv.Atoi(passportData[1])
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, fmt.Sprintf("passport number is required: %s", err.Error()))
			log.Info("Passport number is required", "error", err)
			return
		}

		requestToApi := RequestToApi{
			PassportSerie:  series,
			PassportNumber: number,
		}

		response, err := sendRequestToApi(requestToApi, cfg)
		if err != nil {
			if errors.Is(err, fmt.Errorf("bad request")) {
				render.Status(r, http.StatusBadRequest)
			} else if errors.Is(err, fmt.Errorf("internal server error")) {
				render.Status(r, http.StatusInternalServerError)
			}
			log.Error("Failed to send request to API", "error", err)
			return
		}

		log.Debug("Response from API", "response", response)
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
			render.JSON(w, r, nil)
			log.Error("Error saving user", "error", err)
			return
		}
		log.Debug("Successfully saved user", "user", user)
	}
}

func sendRequestToApi(request RequestToApi, cfg config.Config) (*ResponseFromApi, error) {
	log.Debug("Sending request to API")

	client := http.Client{}

	req, err := http.NewRequest("GET", cfg.ApiUrl+"/info"+"?passportSerie="+strconv.Itoa(request.PassportSerie)+
		"&passportNumber="+strconv.Itoa(request.PassportNumber), nil)
	if err != nil {
		log.Error("Failed creating request to API", "error", err)
		return nil, fmt.Errorf("error creating HTTP request: %s", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed sending request to API", "error", err)
		return nil, fmt.Errorf("error sending HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		log.Debug("Received status bad request from API")
		return nil, fmt.Errorf("bad request")
	} else if resp.StatusCode == http.StatusInternalServerError {
		log.Debug("Received status internal server error from API")
		return nil, fmt.Errorf("internal server error")
	}

	var responseFromApi ResponseFromApi
	err = render.DecodeJSON(resp.Body, &responseFromApi)
	if errors.Is(err, io.EOF) {
		log.Error("Request body is empty", "error", err)
		return nil, fmt.Errorf("request body is empty: %s", err.Error())
	}
	if err != nil {
		log.Error("Failed to decode response body", "error", err)
		return nil, fmt.Errorf("failed to decode response body: %s", err.Error())
	}

	if err := validator.New().Struct(responseFromApi); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)

		log.Error("Invalid response from API", "response", responseFromApi, "error", err)
		return nil, fmt.Errorf("invalid response: %s", err.Error())
	}

	log.Debug("Successfully received response from API", "response", responseFromApi)
	return &responseFromApi, nil
}
