package chat

import (
	"chat/chat/idGenerator"
	"chat/chat/protocol"
	"chat/database"
	"log"
	"time"
)

type MessageService struct {
	IdGenerator       *idGenerator.IdGenerator
	MessageBuff       chan *protocol.MessageSend
	models            *database.Models
	activeConnections map[int64]*Client
	register          chan *Client
	deregister        chan int64
}

func NewMessageService(models *database.Models, idGenerator *idGenerator.IdGenerator) *MessageService {
	return &MessageService{
		IdGenerator:       idGenerator,
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
		case messageSend := <-m.MessageBuff:
			m.handleOneToOneMessage(messageSend)
		}
	}
}

func (m *MessageService) handleOneToOneMessage(message *protocol.MessageSend) {
    messageId := m.IdGenerator.Generate()
    timestamp := time.Now().UTC()
	oneToOneMessage := &database.OneToOneMessage{
		MessageId:   messageId,
		SenderId:    message.SenderId,
		RecipientId: message.RecipientId,
		Content:     message.Content,
		TimeStamp:   timestamp,
	}

	err := m.models.Messages.InsertOneToOneMessage(oneToOneMessage)
	// Message not stored. Needs to notify the user that message not sent
	if err != nil {
		log.Println(err)
		m.failToSend(message)
        return
	}

    m.ackMessage(messageId, message.SenderId, message.RecipientId, message.SentAt, timestamp)
}

func (m *MessageService) ackMessage(messageId int64, senderId int64, recipientId int64,
    sentAt time.Time, timestamp time.Time) {
    messageSentAck := &protocol.MessageSentSuccess{
        MessageId: messageId,
        RecipientId: recipientId,
        Timestamp: timestamp,
        SentAt: sentAt,
    }

    packet := protocol.NewPacket(protocol.MSG_SENT_SUCCESS, messageSentAck)
    client, found := m.activeConnections[senderId]
    // if client not active, then enqueue in message queue
    if !found {
        log.Printf("user %d not online\n")
        return
    }

    client.writePacket(&packet)
}

// TODO: implement
func (m *MessageService) handleGroupMessage(message *protocol.MessageSend) error {
	return nil
}

func (m *MessageService) failToSend(message *protocol.MessageSend) {
	log.Printf("Unable to store message: %+v\n", message)
}
