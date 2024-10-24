package chat

import (
	"log"
	"net/http"
	"strconv"

	"chat/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	userId uint32
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

	if !qp.Has("userId") {
		log.Println("missing userId")
		return
	}
    userId, err := strconv.ParseUint(qp.Get("userId"), 10, 32)
	if err != nil {
        log.Println("error parsing userId: ", err.Error())
	}


	log.Println("attempting to set up socket. Source: ", r.RemoteAddr)

	conn, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error during upgrade:", err)
		return
	}

	client := &Client{
		userId: uint32(userId),
		conn:   conn,
		buff:   make(chan *MessageRequest, 256),
	}

	userConn := &userConnection{
		userId: uint32(userId),
		conn:   client,
	}

	ms.register <-userConn

	go client.readClient(ms, conn)
	go client.writeClient(conn)
}
