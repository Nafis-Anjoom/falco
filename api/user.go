package main

import (
	"chat/database"
	"chat/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type createUserRequest struct {
	Id        uint64 `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (app *application) createUserHandler(writer http.ResponseWriter, request *http.Request) {
	var input createUserRequest

	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	user := &database.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	id, err := app.models.Users.InsertUser(user)
	if err != nil {
		log.Println(err)
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}

	input.Id = id
	log.Println(input)
	utils.WriteJSONResponse(writer, http.StatusCreated, input)
}

func (app *application) getUserByIdHandler(writer http.ResponseWriter, request *http.Request) {
	param := request.PathValue("id")
	if param == "" {
		err := errors.New("missing id param")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	userId, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		err := errors.New("id parameter is not an integer")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	user, err := app.models.Users.GetUserById(userId)
	if err != nil {
		switch {
		case errors.Is(err, database.RecordNotFoundError):
			utils.WriteErrorResponse(writer, request, http.StatusNotFound, err)
		default:
			utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		}
		return
	}

	err = utils.WriteJSONResponse(writer, http.StatusOK, user)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}
}
