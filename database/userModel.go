package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	DuplicateEmailError = errors.New("User with provided email already exist")
	UserNotFoundError   = errors.New("User does not exist")
)

type User struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	IsDeleted    bool   `json:"isDeleted"`
}

type UserModel struct {
	dbPool *pgxpool.Pool
}

func (um *UserModel) InsertUser(user *User) (int64, error) {
	sqlStmt := `
        INSERT INTO public.users(firstName, lastName, email, passwordHash)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	row := um.dbPool.QueryRow(context.Background(), sqlStmt, user.FirstName, user.LastName, user.Email, user.PasswordHash)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			return -1, DuplicateEmailError
		} else {
			return -1, err
		}
	}
	return id, nil
}

func (um *UserModel) GetUserById(userId int64) (*User, error) {
	sqlStmt := `
        SELECT id, firstName, lastName, email, isDeleted, passwordHash
        FROM public.users
        WHERE id = $1`

	row := um.dbPool.QueryRow(context.Background(), sqlStmt, userId)

	var user User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.IsDeleted, &user.PasswordHash)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, UserNotFoundError
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (um *UserModel) GetUserByEmail(email string) (*User, error) {
	sqlStmt := `
        SELECT id, firstName, lastName, email, isDeleted, passwordHash
        FROM public.users
        WHERE email = $1`

	row := um.dbPool.QueryRow(context.Background(), sqlStmt, email)

	var user User
	err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.IsDeleted, &user.PasswordHash)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, UserNotFoundError
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (um *UserModel) DeleteUserById(userId int64) error {
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
