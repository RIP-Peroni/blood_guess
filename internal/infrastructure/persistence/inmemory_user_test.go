package persistence

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryUserRepository(t *testing.T) {
	repo := NewInMemoryUserRepository()

	t.Run("saving and getting user by ID", func(t *testing.T) {
		user := entities.NewUser(12345, "test_user")

		err := repo.Save(user)
		require.NoError(t, err)

		found, err := repo.FindById(user.ID())
		require.NoError(t, err)
		assert.Equal(t, user.ID(), found.ID())
		assert.Equal(t, user.TelegramID(), found.TelegramID())
		assert.Equal(t, user.Username(), found.Username())
	})

	t.Run("search for user by telegram ID", func(t *testing.T) {
		user := entities.NewUser(67890, "another_user")
		err := repo.Save(user)
		require.NoError(t, err)

		found, err := repo.FindByTelegramID(67890)
		require.NoError(t, err)
		assert.Equal(t, user.ID(), found.ID())
		assert.Equal(t, int64(67890), found.TelegramID())
	})

	t.Run("updating user", func(t *testing.T) {
		user := entities.NewUser(11111, "old_username")
		err := repo.Save(user)
		require.NoError(t, err)

		// Changing the balance
		err = user.AddPoints(100)
		require.NoError(t, err)

		err = repo.Update(user)
		require.NoError(t, err)

		updated, err := repo.FindById(user.ID())
		require.NoError(t, err)
		assert.Equal(t, 100, updated.Balance())
	})

	t.Run("searching for a non-existent user returns an error", func(t *testing.T) {
		_, err := repo.FindById(entities.UserID("non-existent"))
		assert.Error(t, err)

		_, err = repo.FindByTelegramID(99999)
		assert.Error(t, err)
	})

	t.Run("cannot save a nil user", func(t *testing.T) {
		err := repo.Save(nil)
		assert.Error(t, err)
	})

	t.Run("cannot update a non-existent user", func(t *testing.T) {
		user := entities.NewUser(22222, "test")
		err := repo.Update(user)
		assert.Error(t, err)
	})
}
