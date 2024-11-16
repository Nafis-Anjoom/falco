package main

import (
	"chat/database"
	"chat/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	models *database.Models
}

func NewUserService(models *database.Models) *UserService {
	return &UserService{models: models}
}

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type getUserResponse struct {
	Id        uint32 `josn:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func (us *UserService) createUserHandler(writer http.ResponseWriter, request *http.Request) {
	var input createUserRequest

	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
        return
	}

	user := &database.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

    passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
        return
	}

    user.PasswordHash = passwordHash

	id, err := us.models.Users.InsertUser(user)
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

func (us *UserService) getUserByIdHandler(writer http.ResponseWriter, request *http.Request) {
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

	user, err := us.models.Users.GetUserById(userId)
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

func (us *UserService) deleteUserByIdHandler(writer http.ResponseWriter, request *http.Request) {
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

	err = us.models.Users.DeleteUserById(userId)
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
