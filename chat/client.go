package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
    conn *websocket.Conn
    buff chan []byte
}

var upgrader = websocket.Upgrader {
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func (c *Client) readClient(server *Server, conn *websocket.Conn) {
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            log.Println("error occured during read from client:", conn.RemoteAddr())
        }

        server.Broadcast <- message
    }
}

func (c *Client) writeClient(conn *websocket.Conn) {
    for {
        message:= <- c.buff 

        err := conn.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Println("error occured during write from client:", conn.RemoteAddr())
        }
    }
}

func ServeWs(server *Server, w http.ResponseWriter, r *http.Request) {
    log.Println("attempting to set up socket. Source: ", r.RemoteAddr)

    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("error during upgrade:", err)
        return
    }

    client := &Client {
        conn: conn,
        buff: make(chan []byte, 256),
    }

    server.Register <- client

    go client.readClient(server, conn)
    go client.writeClient(conn)
}
