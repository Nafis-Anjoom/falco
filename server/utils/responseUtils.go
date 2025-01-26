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
	HTTPCode  int       `json:"http_status_code"`
	Message   string    `json:"message"`
	TimeStamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Details   string    `json:"details"`
}

func WriteJSONResponse(writer http.ResponseWriter, status int, data any) error {
	writer.WriteHeader(status)
	writer.Header().Set("Content-Type", "application/json")
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	json = append(json, '\n')
	_, err = writer.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func WriteErrorResponse(writer http.ResponseWriter, request *http.Request, statusCode int, err error) {
	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")

	errorResponse := ErrorResponse{
		HTTPCode:  statusCode,
		Message:   errorMessageMap[statusCode],
		TimeStamp: time.Now().UTC(),
		Path:      request.URL.RequestURI(),
		Method:    request.Method,
		Details:   err.Error(),
	}

	json, err := json.MarshalIndent(errorResponse, "", "\t")
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	json = append(json, '\n')
	_, err = writer.Write(json)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
