package chat

import (
	"log"
	"net/http"
	"strconv"

    "chat/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	userId uint64
	conn   *websocket.Conn
	buff   chan *MessageRequest
}

func (c *Client) readClient(ms *MessageService, conn *websocket.Conn) {
	for {
        var msg MessageRequest
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("error occured during read from client:", err)
			continue
		}
		log.Printf("%+v\n", msg)
		ms.MessageBuff <- &msg
	}
}

func (c *Client) writeClient(conn *websocket.Conn) {
	for {
		message := <-c.buff

		err := conn.WriteJSON(*message)
		if err != nil {
			log.Println("error occured during write from client:", conn.RemoteAddr())
		}
	}
}

func ServeWs(ms *MessageService, w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()

	if !qp.Has("userId") || !qp.Has("chatId") {
		log.Println("missing userId or chatId")
		return
	}

	log.Println("attempting to set up socket. Source: ", r.RemoteAddr)
	userId, err := strconv.ParseUint(qp.Get("userId"), 10, 64)
	if err != nil {
		log.Println("error parsing userId")
	}

	log.Println("userId", userId)

	conn, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error during upgrade:", err)
		return
	}

	client := &Client{
		userId: userId,
		conn:   conn,
		buff:   make(chan *MessageRequest, 256),
	}

	userConn := &userConnection{
		userId: userId,
		conn:   client,
	}

	ms.register <- userConn

	go client.readClient(ms, conn)
	go client.writeClient(conn)
}
