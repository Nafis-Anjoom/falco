package chat

import (
	"chat/database"
	"log"
)

type MessageRequest struct {
	SenderId  uint64 `json:"sender_id"`
	ChatId    uint64 `json:"chat_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type userConnection struct {
	userId uint64
	conn   *Client
}

type MessageService struct {
	MessageBuff       chan *MessageRequest
	nextId            uint64
	models            *database.Models
	activeConnections map[uint64]*Client
	register          chan *userConnection
	deregister        chan *userConnection
}

func NewMessageService(models *database.Models) *MessageService {
	return &MessageService{
		MessageBuff:       make(chan *MessageRequest, 512),
		activeConnections: make(map[uint64]*Client),
		register:          make(chan *userConnection),
		deregister:        make(chan *userConnection),
        models: models,
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
			msg := &database.Message{
				ChatId:   msgIn.ChatId,
				SenderId: msgIn.SenderId,
				Content:  msgIn.Content,
			}
			m.models.Messages.InsertMessage(msg)
			log.Println("successfully inserted message")
		}
	}
}
