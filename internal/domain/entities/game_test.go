// internal/domain/entities/game_test.go
package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGame_AddPlayer_NoRole(t *testing.T) {
	t.Run("adding player without any role", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		err := game.AddPlayer("Player 1")
		assert.NoError(t, err)
		assert.Len(t, game.Players(), 1)

		player := game.Players()[0]
		assert.Equal(t, "Player 1", player.Name)
		assert.Equal(t, "", player.RealRole) // Роль не указана
		assert.False(t, player.IsRealRoleSet)
	})
}

func TestGame_SetPlayerRealRole(t *testing.T) {
	t.Run("set real role to player after game", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		// Открываем и закрываем прогнозы, начинаем игру
		err = game.OpenPredictions()
		require.NoError(t, err)
		err = game.ClosePredictions()
		require.NoError(t, err)
		err = game.Start()
		require.NoError(t, err)

		playerID := game.Players()[0].ID

		// Устанавливаем реальную роль
		err = game.SetPlayerRealRole(playerID, "demon")
		assert.NoError(t, err)

		player := game.Players()[0]
		assert.Equal(t, "demon", player.RealRole)
		assert.True(t, player.IsRealRoleSet)
	})

	t.Run("error when setting real role before game starts", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		playerID := game.Players()[0].ID

		// Пытаемся установить роль до начала игры
		err = game.SetPlayerRealRole(playerID, "demon")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "game must be in progress or finished")
	})
}

func TestGame_AllRealRolesSet(t *testing.T) {
	t.Run("check when all real roles are set", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)
		err = game.AddPlayer("Player 2")
		require.NoError(t, err)

		// Начинаем игру
		err = game.OpenPredictions()
		require.NoError(t, err)
		err = game.ClosePredictions()
		require.NoError(t, err)
		err = game.Start()
		require.NoError(t, err)

		// Устанавливаем реальные роли
		err = game.SetPlayerRealRole(game.Players()[0].ID, "demon")
		require.NoError(t, err)
		err = game.SetPlayerRealRole(game.Players()[1].ID, "townsfolk")
		require.NoError(t, err)

		assert.True(t, game.AllRealRolesSet())
	})

	t.Run("check when not all real roles are set", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)
		err = game.AddPlayer("Player 2")
		require.NoError(t, err)

		// Начинаем игру
		err = game.OpenPredictions()
		require.NoError(t, err)
		err = game.ClosePredictions()
		require.NoError(t, err)
		err = game.Start()
		require.NoError(t, err)

		// Устанавливаем реальную роль только для одного игрока
		err = game.SetPlayerRealRole(game.Players()[0].ID, "demon")
		require.NoError(t, err)

		assert.False(t, game.AllRealRolesSet())
	})
}
