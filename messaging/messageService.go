package messaging

import (
	"chat/auth"
	"chat/database"
	"chat/messaging/idGenerator"
	protocol "chat/messaging/protocol_v2"
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
	authService       *auth.AuthService
}

type ChatPreview struct {
	UserId   int64     `json:"userId"`
	UserName string    `json:"userName"`
	Message  string    `json:"message"`
	SentAt   time.Time `json:"sentAt"`
}

func NewMessageService(models *database.Models, idGenerator *idGenerator.IdGenerator, authService *auth.AuthService) *MessageService {
	return &MessageService{
		IdGenerator:       idGenerator,
		authService:       authService,
		MessageBuff:       make(chan *protocol.MessageSend, 512),
		activeConnections: make(map[int64]*Client),
		register:          make(chan *Client),
		deregister:        make(chan int64),
		models:            models,
	}
}

// TODO: implement
func (ms *MessageService) GetChatPreviewsHandler(writer http.ResponseWriter, request *http.Request) {
	chats := []ChatPreview{
		{UserId: 101, UserName: "Alice", Message: "Hey, how are you?", SentAt: time.Now().Add(-time.Minute * 5)},
		{UserId: 102, UserName: "Bob", Message: "I'm good, thanks!", SentAt: time.Now().Add(-time.Minute * 10)},
		{UserId: 103, UserName: "Charlie", Message: "Let's meet up soon.", SentAt: time.Now().Add(-time.Minute * 15)},
		{UserId: 104, UserName: "David", Message: "Sounds good to me!", SentAt: time.Now().Add(-time.Minute * 20)},
		{UserId: 105, UserName: "Eva", Message: "What's the plan for tomorrow?", SentAt: time.Now().Add(-time.Minute * 25)},
		{UserId: 106, UserName: "Frank", Message: "I think we should head to the park.", SentAt: time.Now().Add(-time.Minute * 30)},
		{UserId: 107, UserName: "Grace", Message: "I'm in, let's go!", SentAt: time.Now().Add(-time.Minute * 35)},
		{UserId: 108, UserName: "Hannah", Message: "Is it going to rain?", SentAt: time.Now().Add(-time.Minute * 40)},
		{UserId: 109, UserName: "Ivy", Message: "I heard it's supposed to be sunny.", SentAt: time.Now().Add(-time.Minute * 45)},
		{UserId: 110, UserName: "Jack", Message: "Great, I'll bring a picnic!", SentAt: time.Now().Add(-time.Minute * 50)},
	}

	utils.WriteJSONResponse(writer, http.StatusOK, chats)
}

func (ms *MessageService) InitializeClientHandler(writer http.ResponseWriter, request *http.Request) {
	userId := utils.ContextGetUser(request)
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

func (ms *MessageService) GetMessageThreadHandler(writer http.ResponseWriter, request *http.Request) {
    userId1 := utils.ContextGetUser(request);

	param := request.PathValue("id")
	if param == "" {
		err := errors.New("missing id param")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	userId2, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		err := errors.New("id parameter is not an integer")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

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

func (ms *MessageService) GetTotalPagesCountHandler(writer http.ResponseWriter, request *http.Request) {
    userId1 := utils.ContextGetUser(request);

	param := request.PathValue("id")
	if param == "" {
		err := errors.New("missing id param")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	userId2, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		err := errors.New("id parameter is not an integer")
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	pageCount, err := ms.models.Messages.GetTotalMessagesPages(userId1, userId2)
	if err != nil {
		utils.WriteErrorResponse(writer, request, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONResponse(writer, http.StatusOK, pageCount)
}

func (ms *MessageService) Run() {
	log.Println("started message service")
	for {
		select {
		case user := <-ms.register:
			ms.activeConnections[user.userId] = user
			log.Println("user registered", user.userId)
		case userId := <-ms.deregister:
			if client, found := ms.activeConnections[userId]; found {
                client.conn.Close()
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

	// ms.ackMessage(messageId, message.SenderId, message.RecipientId, message.SentAt, timestamp)
	ms.ackMessage(messageId, message, timestamp)

	messageReceive := &protocol.MessageReceieve{
		Id:          messageId,
		SenderId:    message.SenderId,
		RecipientId: message.RecipientId,
		Content:     message.Content,
		Timestamp:   timestamp,
	}

	ms.sendMessageToRecipient(messageReceive)
}

func (ms *MessageService) sendMessageToRecipient(message *protocol.MessageReceieve) {
	client, found := ms.activeConnections[message.RecipientId]
	// if client not active, then enqueue in message queue
	// TODO: implement messageQueue using redis
	if !found {
		log.Printf("user %d not online\n", message.SenderId)
		return
	}

	packet := protocol.NewPacket(protocol.MSG_RECEIVE, message)

	client.writePacket(&packet)
}

// func (ms *MessageService) ackMessage(messageId int64, senderId int64, recipientId int64,
func (ms *MessageService) ackMessage(messageId int64, message *protocol.MessageSend, timestamp time.Time) {
	messageSentAck := &protocol.MessageSentSuccess{
		MessageId:   messageId,
		RecipientId: message.RecipientId,
		Timestamp:   timestamp,
		SentAt:      message.SentAt,
        LocalUUID: message.LocalUUID,
	}

	packet := protocol.NewPacket(protocol.MSG_SENT_SUCCESS, messageSentAck)
	client, found := ms.activeConnections[message.SenderId]
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
