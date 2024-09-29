package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Welcome to the homepage!")
    })

    mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("This is the about page.")
        echo(w, r)
    })

    mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("This is the echo page.")
        echo(w, r)
    })

    s := http.Server {
        Addr: ":3000",
        Handler: mux,
    }

    s.ListenAndServe()
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
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("error", err)
        }

        log.Println("message receieved:", message)
        
        err = conn.WriteMessage(messageType, message)
        if err != nil {
            log.Println("error", err)
        }
    }
}
