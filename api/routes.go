package main

import (
	"chat/chat"
	"log"
	"net/http"
)

func (app *application) routes() http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Welcome to the homepage!")
    })

    mux.HandleFunc("/echo", app.echoHandler)

    mux.HandleFunc("GET /chat", func(w http.ResponseWriter, r *http.Request) {
        chat.ServeWs(app.messageService, w, r)
    })

    mux.HandleFunc("GET /user/{id}", app.userService.getUserByIdHandler) 
    mux.HandleFunc("DELETE /user/{id}", app.userService.deleteUserByIdHandler) 
    mux.HandleFunc("POST /user", app.userService.createUserHandler) 

    return mux
}
