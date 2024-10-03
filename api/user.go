package main

import (
	// "chat/database"
	"encoding/json"
	"net/http"
)

type createUserRequest struct {
    PhoneNumber string `json:"phoneNumber"`
}

func createUserHandler(writer http.ResponseWriter, request *http.Request) {
    var input createUserRequest

    err := json.NewDecoder(request.Body).Decode(&input)
    if err != nil {
        http.Error(writer, err.Error(), http.StatusBadRequest)
    }

    // user := &database.User {
    //     PhoneNumber: input.PhoneNumber,
    // }
    //
    //
    //

}
