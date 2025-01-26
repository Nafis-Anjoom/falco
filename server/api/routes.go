package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {})

    mux.HandleFunc("GET /echo", app.echoHandler)

    mux.HandleFunc("GET /chat/{id}", app.messageService.GetMessageThreadHandler) 
    mux.HandleFunc("GET /chat/{id}/totalPages", app.messageService.GetTotalPagesCountHandler) 
    mux.HandleFunc("GET /inbox", app.messageService.GetChatPreviewsHandler)
    mux.HandleFunc("GET /ws", app.messageService.InitializeClientHandler)

    mux.HandleFunc("GET /user/me", app.userService.getCurrentUserHandler) 
    mux.HandleFunc("GET /user/{id}", app.userService.getUserByIdHandler) 
    mux.HandleFunc("DELETE /user/{id}", app.userService.deleteUserByIdHandler) 
    mux.HandleFunc("POST /user", app.userService.createUserHandler) 
    mux.HandleFunc("POST /login", app.userService.LoginHandler)
    mux.HandleFunc("POST /logout", app.userService.LogoutHandler)
    mux.HandleFunc("POST /validate", app.userService.ValidateHandler)

    mux.HandleFunc("GET /contacts", app.contactsService.GetContactsHandler) 
    mux.HandleFunc("POST /contacts", app.contactsService.createContactHandler) 

    return mux
}
