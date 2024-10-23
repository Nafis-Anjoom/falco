package chat

import (
	"chat/database"
	"log"
)

type MessageRequest struct {
	SenderId  uint32 `json:"sender_id"`
	ChatId    uint64 `json:"chat_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type userConnection struct {
	userId uint32
	conn   *Client
}

type MessageService struct {
	MessageBuff       chan *MessageRequest
	models            *database.Models
	activeConnections map[uint32]*Client
	register          chan *userConnection
	deregister        chan *userConnection
}

func NewMessageService(models *database.Models) *MessageService {
	return &MessageService{
		MessageBuff:       make(chan *MessageRequest, 512),
		activeConnections: make(map[uint32]*Client),
		register:          make(chan *userConnection),
		deregister:        make(chan *userConnection),
		models:            models,
	}
}

func (m *MessageService) Run() {
	log.Println("started message service")
	for {
		select {
		case user := <-m.register:
			m.activeConnections[user.userId] = user.conn
			log.Println("user registered", user.userId)
		case user := <-m.deregister:
			if _, found := m.activeConnections[user.userId]; found {
				delete(m.activeConnections, user.userId)
			}
		case msgIn := <-m.MessageBuff:
            // converting to signed int because the db is using signed integers
            // have to convert to unsigned during reception
			msg := &database.Message{
				ChatId:   int64(msgIn.ChatId),
				SenderId: int32(msgIn.SenderId),
				Content:  msgIn.Content,
			}
			log.Printf("message received: %+v\n", *msg)
			err := m.models.Messages.InsertMessage(msg)
			if err != nil {
				m.reportNotSent(msgIn, err)
			}
		}
	}
}

func (m *MessageService) reportNotSent(msg *MessageRequest, err error) {
	client, found := m.activeConnections[msg.SenderId]
	if found {
		// TODO: implement a struct for sending message status : sent, not sent, read by
		errorMsg := &MessageRequest{
			SenderId: msg.SenderId,
			ChatId:   msg.ChatId,
			Content:  err.Error(),
		}

		client.buff <- errorMsg
	}
}
