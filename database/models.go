package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	RecordNotFoundError = errors.New("Record not found")
	InsertionError      = errors.New("Insertion error")
)

type Models struct {
	Users    UserModel
	Messages MessageModel
}

func NewModels(dbPool *pgxpool.Pool) *Models {
	return &Models{
		Users:    UserModel{dbPool: dbPool},
		Messages: MessageModel{dbPool: dbPool},
	}
}
