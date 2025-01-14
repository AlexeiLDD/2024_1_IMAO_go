package storage

import (
	"2024_1_IMAO/internal/usecase"
	"errors"
	"math/rand"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	errUserNotExists    = errors.New("user does not exist")
	errUserExists       = errors.New("user already exists")
	errSessionNotExists = errors.New("session does not exist")
)

type UnauthorizedUser struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

type User struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type UsersList struct {
	//Ключ - id сессии, значение - id пользователя
	Sessions   map[string]uint
	Users      map[string]*User
	UsersCount uint
}

type UsersInfo interface {
	UserExists(email string) bool
	CreateUser(email, password string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserBySession(sessionID string) (*User, error)
	GetLastID() uint

	SessionExists(sessionID string) bool
	AddSession(email string) string
	RemoveSession(sessionID string)
}

func randString(length int) string {
	var result string
	for i := 0; i < length; i++ {
		result += string(letters[rand.Intn(len(letters))])
	}

	return result
}

func (active *UsersList) UserExists(email string) bool {
	_, exists := active.Users[email]

	return exists
}

func (active *UsersList) GetLastID() uint {
	active.UsersCount++

	return active.UsersCount
}

func (active *UsersList) CreateUser(email, passwordHash string) (*User, error) {
	if active.UserExists(email) {
		return nil, errUserExists
	}

	active.Users[email] = &User{
		ID:           active.GetLastID(),
		PasswordHash: passwordHash,
		Email:        email,
	}

	return active.Users[email], nil
}

func (active *UsersList) GetUserByEmail(email string) (*User, error) {
	if !active.UserExists(email) {
		return nil, errUserNotExists
	}

	return active.Users[email], nil
}

func (active *UsersList) GetUserBySession(sessionID string) (*User, error) {
	id := active.Sessions[sessionID]

	for _, val := range active.Users {
		if val.ID == id {
			return val, nil
		}
	}

	return nil, errUserNotExists
}

func (active *UsersList) SessionExists(sessionID string) bool {
	_, exists := active.Sessions[sessionID]

	return exists
}

func (active *UsersList) AddSession(email string) string {
	sessionID := randString(32)
	user := active.Users[email]

	active.Sessions[sessionID] = user.ID

	return sessionID
}

func (active *UsersList) RemoveSession(sessionID string) error {
	if !active.SessionExists(sessionID) {
		return errSessionNotExists
	}

	delete(active.Sessions, sessionID)

	return nil
}

func NewActiveUser() *UsersList {
	return &UsersList{
		Sessions: make(map[string]uint, 1),
		Users: map[string]*User{
			"example@mail.ru": {
				ID:           1,
				Email:        "example@mail.ru",
				PasswordHash: usecase.HashPassword("123456"),
			},
		},
		UsersCount: 1,
	}
}
