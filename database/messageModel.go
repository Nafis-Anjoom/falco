package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OneToOneMessage struct {
	MessageId   int64
	SenderId    int64
	RecipientId int64
	Content     string
	TimeStamp   time.Time
}

type MessageModel struct {
	dbPool *pgxpool.Pool
}

func (mm *MessageModel) InsertOneToOneMessage(msg *OneToOneMessage) error {
	sqlStmt := "insert into public.oneToOneMessages(messageId, senderId, recipientId, content, timestamp) values($1, $2, $3, $4, $5)"
	_, err := mm.dbPool.Exec(context.Background(), sqlStmt, msg.MessageId, msg.SenderId,
		msg.RecipientId, msg.Content, msg.TimeStamp)
	if err != nil {
		return fmt.Errorf("%w: %w", InsertionError, err)
	}
	return nil
}

func (mm *MessageModel) GetOneToOneMessage(msgId int64) (*OneToOneMessage, error) {
	return nil, nil
}

func (mm *MessageModel) GetOneToOneMessageThread(userId1, userId2 int64) ([]OneToOneMessage, error) {
    sqlStmt := "select * from public.onetoonemessages where (senderid = $1 and recipientid = $2) or (senderid = $2 and recipientid = $1);"
	rows, _ := mm.dbPool.Query(context.Background(), sqlStmt, userId1, userId2)

	messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[OneToOneMessage])
	if err != nil {
		log.Println("error retrieving chat messages:", err)
		return nil, err
	}

	return messages, nil
}
