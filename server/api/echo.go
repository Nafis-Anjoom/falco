package main

import (
	"log"
	"net/http"
    
	"chat/utils"

	"github.com/gorilla/websocket"
)


func (app *application) echoHandler(w http.ResponseWriter, r *http.Request) {
    log.Println("echo")
    conn, err := utils.Upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }

    defer conn.Close()

    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("error", err)
        }

        log.Println("message receieved:", string(message))
        
        err = conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Println("error", err)
        }
    }
}
