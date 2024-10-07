package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Message struct {
    Id uint64
    ChatId uint64
    SenderId uint64
    Content string
}

type MessageModel struct {
    dbPool *pgxpool.Pool
}

func (mm *MessageModel) InsertMessage(msg *Message) error {

    return nil
}

func (mm *MessageModel) GetMessage(chatId uint64, msgId uint64) (*Message, error) {
    return nil, nil
}

// func (mm *MessageModel) GetMessagesByChat(chatId) ([]*Message, error) {
//
// }
