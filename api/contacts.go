package main

import (
	"chat/database"
	"chat/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type createContactRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Contact struct {
	ContactId int64  `json:"contact_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type ContactsService struct {
	models *database.Models
}

func NewContactsService(models *database.Models) *ContactsService {
	return &ContactsService{models: models}
}

func (cs *ContactsService) createContactHandler(writer http.ResponseWriter, request *http.Request) {
	param := request.PathValue("userId")
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

    authenticatedUser := utils.ContextGetUser(request)
    if authenticatedUser != userId {
        err = errors.New("user not authorized")
        utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, err)
        return
    }

	var input createContactRequest
	err = json.NewDecoder(request.Body).Decode(&input)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	contactUser, err := cs.models.Users.GetUserByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, database.UserNotFoundError):
			utils.WriteErrorResponse(writer, request, http.StatusNotFound, err)
		default:
			utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		}
		return
	}

    contactId := contactUser.Id
    contactName := input.Name

	err = cs.models.Contacts.InsertContact(userId, contactId, contactName)
	if err != nil {
        if errors .Is(err, database.DuplicateContactError) {
            utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
        } else {
            utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
        }
		return
	}

	utils.WriteJSONResponse(writer, http.StatusAccepted, "contact created")
}

func (cs *ContactsService) GetContactsByUserHandler(writer http.ResponseWriter, request *http.Request) {
	param := request.PathValue("userId")
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

    authenticatedUser := utils.ContextGetUser(request)
    if authenticatedUser != userId {
        err = errors.New("user not authorized")
        utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, err)
        return
    }

    // maybe create a DTO type for contacts???
	output, err := cs.models.Contacts.GetContactsByUserId(userId)
	if err != nil {
		switch {
		case errors.Is(err, database.UserNotFoundError):
			utils.WriteErrorResponse(writer, request, http.StatusNotFound, err)
		default:
			utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		}
		return
	}

	err = utils.WriteJSONResponse(writer, http.StatusOK, output)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusInternalServerError, err)
		return
	}
}
