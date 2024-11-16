package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	Id           uint32 `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	IsDeleted    bool   `json:"isDeleted"`
}

type UserModel struct {
	dbPool *pgxpool.Pool
}

func (um *UserModel) InsertUser(user *User) (uint32, error) {
	sqlStmt := `INSERT INTO public.users(firstName, lastName, email, passwordHash) VALUES ($1, $2, $3, $4) RETURNING id`
	row := um.dbPool.QueryRow(context.Background(), sqlStmt, user.FirstName, user.LastName, user.Email, user.PasswordHash)
	var id uint32
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (um *UserModel) GetUserById(userId uint64) (*User, error) {
	sqlStmt := `
        SELECT id, firstName, lastName, email, isDeleted
        FROM public.users
        WHERE id = $1`

	row := um.dbPool.QueryRow(context.Background(), sqlStmt, userId)

	var user User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.IsDeleted)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, RecordNotFoundError
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (um *UserModel) DeleteUserById(userId uint64) error {
	sqlStmt := `UPDATE public.users SET isDeleted = true WHERE id = $1`

	_, err := um.dbPool.Exec(context.Background(), sqlStmt, userId)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return RecordNotFoundError
		default:
			return err
		}
	}

	return nil
}
