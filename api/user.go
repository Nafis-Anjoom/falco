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
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type getUserResponse struct {
	Id        uint32 `josn:"id"`
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

	output := getUserResponse{
		Id:        id,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	utils.WriteJSONResponse(writer, http.StatusCreated, output)
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

	output := getUserResponse{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}

	err = utils.WriteJSONResponse(writer, http.StatusOK, output)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}
}

func (app *application) deleteUserById(writer http.ResponseWriter, request *http.Request) {
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

	err = app.models.Users.DeleteUserById(userId)
	if err != nil {
		switch {
		case errors.Is(err, database.RecordNotFoundError):
			utils.WriteErrorResponse(writer, request, http.StatusNotFound, err)
		default:
			utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		}
		return
	}

	err = utils.WriteJSONResponse(writer, http.StatusOK, nil)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}
}
