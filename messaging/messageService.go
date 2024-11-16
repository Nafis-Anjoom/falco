package messaging

import (
	"chat/database"
	"chat/messaging/idGenerator"
	"chat/messaging/protocol"
	"chat/utils"
	"errors"
	"log"
	"net/http"
	"strconv"
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

func (ms *MessageService) GetMessageThreadHandler(writer http.ResponseWriter, request *http.Request) {
	qp := request.URL.Query()

    var err error
	if !qp.Has("userId1") {
        err = errors.New("missing userId1")
        utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
        return
	}
	if !qp.Has("userId2") {
        err = errors.New("missing userId2")
        utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
        return
	}
    
    userId1, err := strconv.ParseInt(qp.Get("userId1"), 10, 64)
    userId2, err := strconv.ParseInt(qp.Get("userId2"), 10, 64)
    
    if err != nil {
        utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
        return
    }

    messages, err := ms.models.Messages.GetOneToOneMessageThread(userId1, userId2)
    if err != nil {
        utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
        return
    }

    utils.WriteJSONResponse(writer, http.StatusOK, messages)
}

func (ms *MessageService) Run() {
	log.Println("started message service")
	for {
		select {
		case user := <-ms.register:
			ms.activeConnections[user.userId] = user
			log.Println("user registered", user.userId)
		case userId := <-ms.deregister:
			if _, found := ms.activeConnections[userId]; found {
				delete(ms.activeConnections, userId)
			}
		case messageSend := <-ms.MessageBuff:
			ms.handleOneToOneMessage(messageSend)
		}
	}
}

func (ms *MessageService) handleOneToOneMessage(message *protocol.MessageSend) {
	messageId := ms.IdGenerator.Generate()
	timestamp := time.Now().UTC()
	oneToOneMessage := &database.OneToOneMessage{
		MessageId:   messageId,
		SenderId:    message.SenderId,
		RecipientId: message.RecipientId,
		Content:     message.Content,
		TimeStamp:   timestamp,
	}

	err := ms.models.Messages.InsertOneToOneMessage(oneToOneMessage)
	// Message not stored. Needs to notify the user that message not sent
	if err != nil {
		log.Println(err)
		ms.failToSend(message)
		return
	}

	ms.ackMessage(messageId, message.SenderId, message.RecipientId, message.SentAt, timestamp)

	messageReceive := &protocol.MessageReceieve{
		Id:          messageId,
		SenderId:    message.SenderId,
		RecipientId: message.RecipientId,
		Content:     message.Content,
		Timestamp:   timestamp,
	}

	ms.sendToRecipient(messageReceive)
}

func (ms *MessageService) sendToRecipient(message *protocol.MessageReceieve) {
	client, found := ms.activeConnections[message.RecipientId]
	// if client not active, then enqueue in message queue
	if !found {
		log.Printf("user %d not online\n")
		return
	}

    packet := protocol.NewPacket(protocol.MSG_RECEIVE, message)

	client.writePacket(&packet)
}

func (ms *MessageService) ackMessage(messageId int64, senderId int64, recipientId int64,
	sentAt time.Time, timestamp time.Time) {
	messageSentAck := &protocol.MessageSentSuccess{
		MessageId:   messageId,
		RecipientId: recipientId,
		Timestamp:   timestamp,
		SentAt:      sentAt,
	}

	packet := protocol.NewPacket(protocol.MSG_SENT_SUCCESS, messageSentAck)
	client, found := ms.activeConnections[senderId]
	// if client not active, then enqueue in message queue
	if !found {
		log.Printf("user %d not online\n")
		return
	}

	client.writePacket(&packet)
}

// TODO: implement
func (ms *MessageService) handleGroupMessage(message *protocol.MessageSend) error {
	return nil
}

func (ms *MessageService) failToSend(message *protocol.MessageSend) {
	log.Printf("Unable to store message: %+v\n", message)
}
