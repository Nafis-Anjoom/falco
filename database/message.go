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

func (mm *MessageModel) GetMessage(chatId uint64, msgId uint64) (*Message, error) {
    msg, ok := mm.messageStorage[msgId]
    if !ok {
        log.Printf("message not found. ChatId: %d. msgId: %d\n", chatId, msgId)
        return nil, RecordNotFoundError
    }

    return msg, nil
}

// func (mm *MessageModel) GetMessagesByChat(chatId) ([]*Message, error) {
//
// }
