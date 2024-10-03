package database

type User struct {
    Id uint64
    PhoneNumber string
}

type UserModel struct {
    userStorage map[uint64]*User
    nextId uint64
}

func (um *UserModel) InsertUser(user *User) error {
    um.userStorage[um.nextId] = user
    user.Id = um.nextId
    um.nextId += 1

    return nil
}

func (um *UserModel) GetUser(userId uint64) (*User, error) {
    user, ok := um.userStorage[userId]
    if !ok {
        return nil, RecordNotFoundError
    }

    return user, nil
}
