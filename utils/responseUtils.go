package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var errorMessageMap = map[int]string{
	http.StatusNotFound:            "the requested resource is not found",
	http.StatusInternalServerError: "the server encountered an error and could not process your request",
	http.StatusBadRequest:          "invalid request",
	http.StatusUnprocessableEntity: "invalid data",
	http.StatusUnauthorized:        "user not authorized",
}

type ErrorResponse struct {
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Details   string    `json:"details"`
}

func WriteJSONResponse(writer http.ResponseWriter, status int, data any) error {
	writer.WriteHeader(status)
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	json = append(json, '\n')
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func WriteErrorResponse(writer http.ResponseWriter, request *http.Request, statusCode int, err error) {
	writer.WriteHeader(statusCode)

	errorResponse := ErrorResponse{
		Message:   errorMessageMap[statusCode],
		TimeStamp: time.Now().UTC(),
		Path:      request.URL.RequestURI(),
		Details:   err.Error(),
	}

	json, err := json.MarshalIndent(errorResponse, "", "\t")
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	json = append(json, '\n')
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(json)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
