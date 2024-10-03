package database

import "errors"

var (
    RecordNotFoundError = errors.New("record not found")
)

type Models struct {
    Users UserModel
    Messages MessageModel
}

func NewModels() Models {
    return Models{
        Users: UserModel{
            userStorage: make(map[uint64]*User),
            nextId: 0,
        },
        Messages: MessageModel{
            messageStorage: make(map[uint64]*Message),
            nextId: 0,
        },
    }
}
