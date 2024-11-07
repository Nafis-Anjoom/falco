package database

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OneToOneMessage struct {
	Id         int64
	SenderId   int64
	ReceiverId int64
	Content    string
	TimeStamp  time.Time
}

type MessageModel struct {
	dbPool *pgxpool.Pool
}

func (mm *MessageModel) InsertOneToOneMessage(msg *OneToOneMessage) error {
	sqlStmt := `insert into public.oneToOneMessages(senderId, receiverId, content, timestamp) values($1, $2, $3, $4)`
	_, err := mm.dbPool.Exec(context.Background(), sqlStmt, msg.SenderId, msg.ReceiverId, msg.Content, msg.TimeStamp)
	if err != nil {
		log.Println("error inserting message:", err)
		return err
	}
	return nil
}

func (mm *MessageModel) GetOneToOneMessage(msgId int64) (*OneToOneMessage, error) {
	return nil, nil
}

func (mm *MessageModel) GetMessagesByChat(chatId int64) ([]OneToOneMessage, error) {
	sqlStmt := `select id, senderId, content from public.messages where chatId = $1`
	rows, _ := mm.dbPool.Query(context.Background(), sqlStmt, chatId)

	messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[OneToOneMessage])
	if err != nil {
		log.Println("error retrieving chat messages:", err)
		return nil, err
	}

	return messages, nil
}
