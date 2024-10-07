package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
    Id uint64 `json:"id"`
    FirstName string `json:"first_name"`
    LastName string `json:"last_name"`
    Email string `json:"email"`
}

type UserModel struct {
    dbPool *pgxpool.Pool
}

func (um *UserModel) InsertUser(user *User) (uint64, error) {
    sqlStmt := `INSERT INTO public.user(firstName, lastName, email) VALUES ($1, $2, $3) RETURNING id`
    row := um.dbPool.QueryRow(context.Background(), sqlStmt, user.FirstName, user.LastName, user.Email)
    var id uint64
    err := row.Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (um *UserModel) GetUserById(userId uint64) (*User, error) {
    sqlStmt := `SELECT * FROM public.user where id = $1`

    row := um.dbPool.QueryRow(context.Background(), sqlStmt, userId)

    var user User
    err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email)
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
