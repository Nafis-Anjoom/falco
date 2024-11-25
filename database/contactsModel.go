package database

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Contact struct {
	ContactId int64  `json:"contact_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

type ContactsModel struct {
	dbPool *pgxpool.Pool
}

var (
    DuplicateContactError = errors.New("Contact with email already exists")
)

// TODO: implement
func (cm *ContactsModel) InsertContact(userId int64, contactId int64, contactName string) error {
	sqlStmt := "insert into public.contacts(userId, contactId, name) values($1, $2, $3)"
	_, err := cm.dbPool.Exec(context.Background(), sqlStmt, userId, contactId, contactName)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "contacts_pkey" (SQLSTATE 23505)` {
			return DuplicateContactError 
		} else {
            return fmt.Errorf("%w: %w", InsertionError, err)
		}
	}
	return nil
}

// TODO: implement
func (cm *ContactsModel) DeleteContact(userId, contactId int64) error {
	return nil
}

// TODO: implement
func (cm *ContactsModel) GetContactsByUserId(userId int64) ([]Contact, error) {
	sqlStmt := `
    SELECT contacts.contactId, contacts.name, users.email
    FROM public.contacts
    JOIN public.users on contacts.contactId = users.id
    WHERE contacts.userId = $1`

	rows, _ := cm.dbPool.Query(context.Background(), sqlStmt, userId)
	contacts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Contact])
	if err != nil {
		log.Println("error retrieving contacts:", err)
		return nil, err
	}

	return contacts, nil
}
