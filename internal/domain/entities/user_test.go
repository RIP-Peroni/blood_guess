package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUser_NewUser(t *testing.T) {
	t.Run("creating new user", func(t *testing.T) {
		user := NewUser(12345, "test_user")

		assert.Equal(t, int64(12345), user.TelegramID())
		assert.Equal(t, "test_user", user.Username())
		assert.Equal(t, 0, user.Balance())
		assert.NotEmpty(t, user.ID())
		assert.WithinDuration(t, time.Now(), user.CreatedAt(), time.Second)
	})

	t.Run("user with negative balance cannot be created", func(t *testing.T) {
		user := NewUser(12345, "test_user")
		err := user.AddPoints(-100)
		assert.Error(t, err)
		assert.Equal(t, 0, user.Balance())
	})
}

func TestUser_AddPoints(t *testing.T) {
	user := NewUser(12345, "test_user")

	t.Run("successfully add points", func(t *testing.T) {
		err := user.AddPoints(100)
		assert.NoError(t, err)
		assert.Equal(t, 100, user.Balance())

		err = user.AddPoints(50)
		assert.NoError(t, err)
		assert.Equal(t, 150, user.Balance())
	})

	t.Run("adding negative points causes an error", func(t *testing.T) {
		err := user.AddPoints(-10)
		assert.Error(t, err)
		assert.Equal(t, 150, user.Balance())
	})

	t.Run("adding zero points is acceptable", func(t *testing.T) {
		err := user.AddPoints(0)
		assert.NoError(t, err)
		assert.Equal(t, 150, user.Balance())
	})
}

func TestUser_DeductPoints(t *testing.T) {
	user := NewUser(12345, "test_user")
	err := user.AddPoints(200)
	assert.NoError(t, err)

	t.Run("successful write-off of points", func(t *testing.T) {
		err := user.DeductPoints(100)
		assert.NoError(t, err)
		assert.Equal(t, 100, user.Balance())
	})
	t.Run("writing off more points causes an error", func(t *testing.T) {
		err := user.DeductPoints(150)
		assert.Error(t, err)
		assert.Equal(t, 100, user.Balance())
	})
	t.Run("writing off negative points causes an error", func(t *testing.T) {
		err := user.DeductPoints(-10)
		assert.Error(t, err)
		assert.Equal(t, 100, user.Balance())
	})
}

func TestUser_UserCanMakePrediction(t *testing.T) {
	t.Run("user can make prediction", func(t *testing.T) {
		user := NewUser(12345, "test_user")
		canPredict := user.CanMakePrediction("game-123") //just a stub for now. Logic will be in Game entity
		assert.True(t, canPredict)
	})
}
