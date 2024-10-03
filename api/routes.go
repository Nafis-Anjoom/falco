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
    mux.HandleFunc("/temp", func(w http.ResponseWriter, r *http.Request) {
        qp := r.URL.Query()

        log.Println(qp)
    })
    mux.HandleFunc("/echo", app.echoHandler)

    mux.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
        chat.ServeWs(app.messageService, w, r)
    })

    return mux
}

