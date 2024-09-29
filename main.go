package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Welcome to the homepage!")
    })

    http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("This is the about page.")
    })

    http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("This is the echo page.")
        echo(w, r)
    })

    // http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
    //     conn, err := upgrader.Upgrade(w, r, nil)
    //     if err != nil {
    //         log.Println(err)
    //     }
    //
    //     client := client.Client {
    //         SendBuf: make(chan []byte, 1024),
    //         Conn: conn,
    //     }
    // })

    http.ListenAndServe(":3000", nil)
}

var upgrader = websocket.Upgrader {
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func echo(w http.ResponseWriter, r *http.Request) {
    log.Println("echo")
    conn, err := upgrader.Upgrade(w, r, nil)
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
