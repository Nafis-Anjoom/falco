package chat

import (
	"chat/chat/protocol"
	"chat/database"
	"fmt"
	"log"
	"time"
)

type MessageService struct {
	MessageBuff       chan *protocol.MessageSend
	models            *database.Models
	activeConnections map[int64]*Client
	register          chan *Client
	deregister        chan int64
}

func NewMessageService(models *database.Models) *MessageService {
	return &MessageService{
		MessageBuff:       make(chan *protocol.MessageSend, 512),
		activeConnections: make(map[int64]*Client),
		register:          make(chan *Client),
		deregister:        make(chan int64),
		models:            models,
	}
}

func (m *MessageService) Run() {
	log.Println("started message service")
	for {
		select {
		case user := <-m.register:
			m.activeConnections[user.userId] = user
			log.Println("user registered", user.userId)
		case userId := <-m.deregister:
			if _, found := m.activeConnections[userId]; found {
				delete(m.activeConnections, userId)
			}
        // TODO: handle messageSend
        case messageSend := <- m.MessageBuff:
            m.handleOneToOneMessage(messageSend)
            fmt.Printf("Message in MessageService: %+v\n", messageSend)
		}
	}
}

func (m *MessageService) handleOneToOneMessage(message *protocol.MessageSend) error {
    oneToOneMessage := &database.OneToOneMessage {
        SenderId: message.SenderId,
        ReceiverId: message.RecipientId,
        Content: message.Content,
        TimeStamp: time.Now().UTC(),
    }

    err := m.models.Messages.InsertOneToOneMessage(oneToOneMessage)
    // Message not stored. Needs to notify the user that message not sent
    if err != nil {
        m.failToSend(message)
    }
    return nil
}

// TODO: implement
func (m *MessageService) handleGroupMessage(message *protocol.MessageSend) error {
    return nil
}

func (m *MessageService) failToSend(message *protocol.MessageSend) {
    log.Println("Unable to store message: %+v\n", message)
}
