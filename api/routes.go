package main

import (
	"chat/messaging"
	"net/http"
)

func (app *application) routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {})
    mux.HandleFunc("OPTIONS /*", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        writer.WriteHeader(http.StatusNoContent)
    })

    mux.HandleFunc("GET /echo", app.echoHandler)

    mux.HandleFunc("GET /thread", app.messageService.GetMessageThreadHandler) 
    mux.HandleFunc("GET /ws2", app.messageService.InitializeClientHandler)
    mux.HandleFunc("GET /chat/preview", app.messageService.GetChatPreviewsHandler)
    mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
        messaging.ServeWs(app.messageService, w, r)
    })

    mux.HandleFunc("GET /user/me", app.userService.getCurrentUserHandler) 
    mux.HandleFunc("GET /user/{id}", app.userService.getUserByIdHandler) 
    mux.HandleFunc("DELETE /user/{id}", app.userService.deleteUserByIdHandler) 
    mux.HandleFunc("POST /user", app.userService.createUserHandler) 
    mux.HandleFunc("POST /login", app.userService.LoginHandler)

    mux.HandleFunc("GET /contacts", app.contactsService.GetContactsHandler) 
    mux.HandleFunc("POST /contacts", app.contactsService.createContactHandler) 

    return mux
}
