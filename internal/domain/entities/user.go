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

// AddPoints adds points (positive amount)
func (u *User) AddPoints(amount int) error {
	if amount < 0 {
		return ErrNegativePoints
	}
	u.balance += amount
	return nil
}

// DeductPoints subtracts points but does not allow the balance to become negative
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

// ChangeBalance changes the balance, but does not allow it to become negative
func (u *User) ChangeBalance(amount int) (int, error) {
	newBalance := u.balance + amount

	// If the balance becomes negative, set it to 0
	if newBalance < 0 {
		u.balance = 0
		return 0, nil
	}

	u.balance = newBalance
	return newBalance, nil
}

// CanMakePrediction - stub
func (u *User) CanMakePrediction(gameID string) bool {
	return true
}

func (u *User) ID() UserID           { return u.id }
func (u *User) TelegramID() int64    { return u.telegramID }
func (u *User) Username() string     { return u.username }
func (u *User) Balance() int         { return u.balance }
func (u *User) CreatedAt() time.Time { return u.createdAt }
