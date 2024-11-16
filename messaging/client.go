package messaging

import (
	"log"
	"net/http"
	"strconv"

	"chat/messaging/protocol"
	"chat/utils"

	"github.com/gorilla/websocket"
)

type Client struct {
	userId        int64
	conn          *websocket.Conn
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

func (client *Client) writePacket(packet *protocol.Packet) {
    err := client.conn.WriteMessage(websocket.BinaryMessage, packet.ToBytes())
    if err != nil {
        log.Println(err)
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
	}

	ms.register <- client

	go client.readClient(ms)
}
