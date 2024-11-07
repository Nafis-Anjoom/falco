package chat

import (
	"chat/chat/protocol"
	"chat/database"
	"fmt"
	"log"
	"time"
)

type userConnection struct {
	userId int64
	client   *Client
}

type MessageService struct {
	MessageBuff       chan *protocol.MessageReceieve
	models            *database.Models
	activeConnections map[int64]*Client
	register          chan *Client
	deregister        chan int64
}

func NewMessageService(models *database.Models) *MessageService {
	return &MessageService{
		MessageBuff:       make(chan *protocol.MessageReceieve, 512),
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
        // TODO: handle messageReceive
        case messageReceive := <- m.MessageBuff:
            fmt.Printf("Message received: %+v\n", messageReceive)
		}
	}
}

func (m *MessageService) handleOneToOneMessage(message *protocol.MessageReceieve) error {
    oneToOneMessage := &database.OneToOneMessage {
        SenderId: message.SenderId,
        ReceiverId: message.RecipientId,
        Content: message.Content,
        TimeStamp: time.Now().UTC(),
    }

    err := m.models.Messages.InsertOneToOneMessage(oneToOneMessage)
    if err != nil {
        log.Println(err)
    }
    return nil
}

// TODO: implement
func (m *MessageService) handleGroupMessage(message *protocol.MessageReceieve) error {
    return nil
}
