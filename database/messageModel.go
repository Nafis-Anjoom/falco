package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Message struct {
	Id       int64
	ChatId   int64
	SenderId int32
	Content  string
}

type MessageModel struct {
	dbPool *pgxpool.Pool
}

func (mm *MessageModel) InsertMessage(msg *Message) error {
    sqlStmt := `insert into public.messages(chatId, senderId, content) values($1, $2, $3)`
    _, err := mm.dbPool.Exec(context.Background(), sqlStmt, msg.ChatId, msg.SenderId, msg.Content)
    if err != nil {
        log.Println("error inserting message:", err)
        return err
    }
    return nil
}

func (mm *MessageModel) GetMessage(chatId uint64, msgId uint64) (*Message, error) {
	return nil, nil
}

func (mm *MessageModel) GetMessagesByChat(chatId uint64) ([]Message, error) {
    sqlStmt := `select id, senderId, content from public.messages where chatId = $1`
    rows, _ := mm.dbPool.Query(context.Background(), sqlStmt, chatId)

    messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[Message])
    if err != nil {
        log.Println("error retrieving chat messages:", err)
        return nil, err
    }
    
    return messages, nil
}
