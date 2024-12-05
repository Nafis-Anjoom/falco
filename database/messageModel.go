package database

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OneToOneMessage struct {
	MessageId   int64            `json:"messageId"`
	SenderId    int64            `json:"senderId"`
	RecipientId int64            `json:"recipientId"`
	Content     string           `json:"content"`
	TimeStamp   time.Time        `json:"timestamp"`
	SeenAt      pgtype.Timestamp `json:"seenAt"`
	DeliveredAt pgtype.Timestamp `json:"deliveredAt"`
}

type MessageModel struct {
	dbPool *pgxpool.Pool
}

const MESSAGES_PER_PAGE = 30

func (mm *MessageModel) GetTotalMessagesPages(userId1, userId2 int64) (int, error) {
	sqlStmt := `
    SELECT COUNT(messageId)
    FROM public.oneToOneMessages
    WHERE (senderId = $1 AND recipientId = $2)
    OR (senderId = $2 AND recipientId = $2);
    `

	row := mm.dbPool.QueryRow(context.Background(), sqlStmt, userId1, userId2)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}

	return int(math.Ceil(float64(count) / MESSAGES_PER_PAGE)), nil
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
	sqlStmt := "select  from public.onetoonemessages where (senderid = $1 and recipientid = $2) or (senderid = $2 and recipientid = $1);"
	rows, _ := mm.dbPool.Query(context.Background(), sqlStmt, userId1, userId2)

	messages, err := pgx.CollectRows(rows, pgx.RowToStructByName[OneToOneMessage])
	if err != nil {
		log.Println("error retrieving chat messages:", err)
		return nil, err
	}

	return messages, nil
}
