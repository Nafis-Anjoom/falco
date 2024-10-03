package database

import "log"

type Message struct {
    Id uint64
    ChatId uint64
    SenderId uint64
    Content string
}

type MessageModel struct {
    messageStorage map[uint64]*Message
    nextId uint64
}

func (mm *MessageModel) InsertMessage(msg *Message) error {
    mm.messageStorage[mm.nextId] = msg
    log.Println("message stored", msg)
    msg.Id = mm.nextId
    mm.nextId += 1

    return nil
}

func (mm *MessageModel) GetUser(msgId uint64) (*Message, error) {
    msg, ok := mm.messageStorage[msgId]
    if !ok {
        return nil, RecordNotFoundError
    }

    return msg, nil
}
