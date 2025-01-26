package messaging

import (
	"log"
	"net/http"

	protocol "chat/messaging/protocol_v2"
	"chat/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	userId int64
	conn   *websocket.Conn
}

func (client *Client) readClient(ms *MessageService) {
	for {
		_, buff, err := client.conn.ReadMessage()
		if err != nil {
			// TODO: handle error "websocket: bad close code 2545"
			switch {
			case websocket.IsCloseError(err, websocket.CloseMessage):
				log.Println("Closing connection with client: ", client.conn.RemoteAddr().String())
			case websocket.IsCloseError(err, websocket.CloseAbnormalClosure):
				log.Println("Abnormal closure. Closing connection with client: ", client.conn.RemoteAddr().String())
			case websocket.IsCloseError(err, websocket.CloseGoingAway):
				log.Println("Client going away. Closing connection with client: ", client.conn.RemoteAddr().String())
			default:
				log.Println("unknown error: ", err)
			}
            client.conn.Close()
			break
		}

		packet := protocol.PacketFromBytes(buff)
		switch packet.PayloadType {
		case protocol.MSG_SEND:
			var message protocol.MessageSend
			err := message.UnmarshalBinary(packet.Payload)
			if err == nil {
				ms.MessageBuff <- &message
			}

		default:
			log.Println("unhandled packet type: ", packet.PayloadType)
		}

		// TODO: notify client if error
		if err != nil {
			log.Println(err)
		}
	}
}

func (client *Client) writePacket(packet *protocol.Packet) {
	err := client.conn.WriteMessage(websocket.BinaryMessage, packet.ToBytes())
	if err != nil {
		log.Println(err)
	}
}

func ServeWs(ms *MessageService, writer http.ResponseWriter, request *http.Request) {
    userId := utils.ContextGetUser(request)

	log.Println("attempting to set up socket. Source: ", request.RemoteAddr)
	conn, err := utils.Upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("error during upgrade:", err)
		return
	}

	client := &Client{
		userId: userId,
		conn:   conn,
	}

	ms.register <- client

	go client.readClient(ms)
}
