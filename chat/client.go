package chat

import (
	"log"
	"net/http"
	"strconv"

	"chat/chat/protocol"
	"chat/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	userId        int64
	conn          *websocket.Conn
	receiveBuffer chan *protocol.MessageReceieve
}

func (client *Client) readClient(ms *MessageService, conn *websocket.Conn) {
	for {
        _, buff, err := conn.ReadMessage()
        if err != nil {
            // TODO: handle error "websocket: bad close code 2545"
            switch {
            case websocket.IsCloseError(err, websocket.CloseMessage):
                log.Println("Closing connection with client: ", conn.RemoteAddr().String())
            case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
                log.Println("Abnormal closure. Closing connection with client: ", conn.RemoteAddr().String())
            default: 
                log.Println("unknown error: ", err)
            }
            break
        }

        packet := protocol.PacketFromBytes(buff)
        switch packet.PayloadType {
        case protocol.MSG_SEND:
            err = client.handleMessageSend(ms, packet)
        default:
            log.Println("unhandled packet type: ", packet.PayloadType)
        }

        // TODO: notify client if error
        if err != nil {
            log.Println(err)
        }

		// ms.MessageBuff <- &msg
	}
}

func (client *Client) handleMessageSend(ms *MessageService, packet *protocol.Packet) error {
    var err error
    var message protocol.MessageSend
    err = message.UnmarshalBinary(packet.Payload)
    if err != nil {
        return err
    }
    ms.MessageBuff <- &message

    return err
}

func (client *Client) writeClient(conn *websocket.Conn) {
	for {
		message := <-client.receiveBuffer

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
	userId, err := strconv.ParseInt(qp.Get("userId"), 10, 32)
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
		userId:        userId,
		conn:          conn,
		receiveBuffer: make(chan *protocol.MessageReceieve, 256),
	}

	ms.register <- client

	go client.readClient(ms, conn)
	go client.writeClient(conn)
}
