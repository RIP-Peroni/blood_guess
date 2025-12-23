package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNegativePoints     = errors.New("negative points")
	ErrInsufficientPoints = errors.New("insufficient points")
)

type UserID string

type User struct {
	id         UserID
	telegramID int64
	username   string
	balance    int
	createdAt  time.Time
}

func NewUser(telegramID int64, username string) *User {
	return &User{
		id:         UserID(uuid.New().String()),
		telegramID: telegramID,
		username:   username,
		balance:    0,
		createdAt:  time.Now(),
	}
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) TelegramID() int64 {
	return u.telegramID
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Balance() int {
	return u.balance
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) AddPoints(amount int) error {
	if amount < 0 {
		return ErrNegativePoints
	}
	u.balance += amount
	return nil
}

func (u *User) DeductPoints(amount int) error {
	if amount < 0 {
		return ErrNegativePoints
	}
	if amount > u.balance {
		return ErrInsufficientPoints
	}
	u.balance -= amount
	return nil
}

func (u *User) CanMakePrediction(gameID string) bool {
	//stub
	return true
}
