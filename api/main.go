package main

import (
	"chat/chat"
	"chat/database"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type application struct {
    messageService *chat.MessageService
}

func (app *application) serve() {
    srv := &http.Server{
        Addr: ":3000",
        Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: time.Minute,
        WriteTimeout: time.Minute,
    }

    go app.messageService.Run()

    log.Printf("starting server on localhost:%d", 3000)
    log.Fatal(srv.ListenAndServe())
}

func main() {
    models := database.NewModels()
    app := &application{
        messageService: chat.NewMessageService(&models),
    }
    app.serve()
}

var upgrader = websocket.Upgrader {
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
