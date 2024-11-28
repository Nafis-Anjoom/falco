package main

import (
	"chat/database"
	"chat/utils"
	"encoding/json"
	"errors"
	"net/http"
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
    userId := utils.ContextGetUser(request)

	var input createContactRequest
    err := json.NewDecoder(request.Body).Decode(&input)
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

func (cs *ContactsService) GetContactsHandler(writer http.ResponseWriter, request *http.Request) {
    userId := utils.ContextGetUser(request)
    query := request.URL.Query()

    var output []database.Contact
    var err error
    if query.Has("q") {
        output, err = cs.models.Contacts.GetFilteredContacts(userId, query.Get("q"))
    } else {
        output, err = cs.models.Contacts.GetContacts(userId)
    }

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
