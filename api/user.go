package main

import (
	"chat/auth"
	"chat/database"
	"chat/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type UserService struct {
	models      *database.Models
	authService *auth.AuthService
}

func NewUserService(models *database.Models, as *auth.AuthService) *UserService {
	return &UserService{models: models, authService: as}
}

type createUserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type createUserResponse struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type getUserResponse struct {
	Id        int64  `josn:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Id    int64  `josn:"id"`
	Token string `json:"token"`
}

var (
    emailPasswordMismatchError = errors.New("email or password incorrect")
)

func (us *UserService) LoginHandler(writer http.ResponseWriter, request *http.Request) {
	var input loginRequest
	err := json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	user, err := us.models.Users.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, database.UserNotFoundError) {
			utils.WriteErrorResponse(writer, request, http.StatusNotFound, err)
		} else {
			utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		}
		return
	}

	if !us.authService.PasswordMatches(input.Password, user.PasswordHash) {
		err = errors.New("Password does not match")
		utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, emailPasswordMismatchError)
		return
	}

	tokenString, err := us.authService.NewToken(user.Id)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}

	output := loginResponse{
        Id: user.Id,
		Token: tokenString,
	}
    
    authCookie := &http.Cookie{
        Name: "authToken",
        Value: tokenString,
        Path: "/",
        Domain: "localhost",
        SameSite: http.SameSiteLaxMode,
        MaxAge: 24 * 60 * 60,
        Secure: false,
        HttpOnly: true,
    }

    userIdCookie := &http.Cookie{
        Name: "userId",
        Value: strconv.FormatInt(user.Id, 10),
        Path: "/",
        Domain: "localhost",
        SameSite: http.SameSiteLaxMode,
        MaxAge: 24 * 60 * 60,
        Secure: false,
        HttpOnly: false,
    }
    http.SetCookie(writer, authCookie)
    http.SetCookie(writer, userIdCookie)
	utils.WriteJSONResponse(writer, http.StatusOK, output)
}

func (us *UserService) LogoutHandler(writer http.ResponseWriter, request *http.Request) {
    authCookie := &http.Cookie{
        Name: "authToken",
        Value: "",
        Path: "/",
        Domain: "localhost",
        SameSite: http.SameSiteLaxMode,
        MaxAge: -1,
        Secure: false,
        HttpOnly: true,
    }

    userIdCookie := &http.Cookie{
        Name: "userId",
        Value: "",
        Path: "/",
        Domain: "localhost",
        SameSite: http.SameSiteLaxMode,
        MaxAge: -1,
        Secure: false,
        HttpOnly: false,
    }
    http.SetCookie(writer, authCookie)
    http.SetCookie(writer, userIdCookie)
	utils.WriteJSONResponse(writer, http.StatusOK, "successfully logged out")
}

func (us *UserService) ValidateHandler(writer http.ResponseWriter, request *http.Request) {
    utils.WriteJSONResponse(writer, http.StatusOK, "validated")
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

	passwordHash, err := us.authService.HashPassword(input.Password)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}

	user.PasswordHash = passwordHash
	id, err := us.models.Users.InsertUser(user)
	if err != nil {
		if errors.Is(err, database.DuplicateEmailError) {
			utils.WriteErrorResponse(writer, request, http.StatusUnprocessableEntity, err)
		} else {
			utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		}
		return
	}

	tokenString, err := us.authService.NewToken(id)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}

	output := createUserResponse{
		Id:    id,
		Token: tokenString,
	}

    authCookie := &http.Cookie{
        Name: "authToken",
        Value: tokenString,
        Path: "/",
        Domain: "http://localhost:3001",
        MaxAge: 24 * 60 * 60,
        Secure: false,
        HttpOnly: true,
    }

    userIdCookie := &http.Cookie{
        Name: "userId",
        Value: strconv.FormatInt(id, 10),
        Path: "/",
        Domain: "http://localhost:3001",
        MaxAge: 24 * 60 * 60,
        Secure: false,
        HttpOnly: true,
    }

    http.SetCookie(writer, authCookie)
    http.SetCookie(writer, userIdCookie)
	utils.WriteJSONResponse(writer, http.StatusCreated, output)
}

func (us *UserService) getCurrentUserHandler(writer http.ResponseWriter, request *http.Request) {
    userId := utils.ContextGetUser(request)
	user, err := us.models.Users.GetUserById(userId)
	if err != nil {
		switch {
		case errors.Is(err, database.UserNotFoundError):
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

func (us *UserService) getUserByIdHandler(writer http.ResponseWriter, request *http.Request) {
	param := request.PathValue("id")
	if param == "" {
		err := errors.New("missing id param")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	userId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		err := errors.New("id parameter is not an integer")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	user, err := us.models.Users.GetUserById(userId)
	if err != nil {
		switch {
		case errors.Is(err, database.UserNotFoundError):
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

	userId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		err := errors.New("id parameter is not an integer")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	err = us.models.Users.DeleteUserById(userId)
	if err != nil {
		switch {
		case errors.Is(err, database.UserNotFoundError):
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
