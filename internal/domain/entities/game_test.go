// internal/domain/entities/game_test.go
package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGame_AddPlayer_OptionalRole(t *testing.T) {
	t.Run("adding player without role", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		err := game.AddPlayer("Player 1")
		assert.NoError(t, err)
		assert.Len(t, game.Players(), 1)

		player := game.Players()[0]
		assert.Equal(t, "Player 1", player.Name)
		assert.Equal(t, "", player.AssignedRole) // Роль не указана
		assert.Empty(t, player.RealRole)
		assert.False(t, player.IsRealRoleSet)
	})

	t.Run("adding player with role", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		err := game.AddPlayer("Player 1", "townsfolk")
		assert.NoError(t, err)
		assert.Len(t, game.Players(), 1)

		player := game.Players()[0]
		assert.Equal(t, "Player 1", player.Name)
		assert.Equal(t, "townsfolk", player.AssignedRole)
	})
}

func TestGame_AddPlayers(t *testing.T) {
	t.Run("adding multiple players at once", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		names := []string{"Вася", "Петя", "Коля", "Миша", "Алекс"}
		err := game.AddPlayers(names)
		assert.NoError(t, err)
		assert.Len(t, game.Players(), 5)

		for i, name := range names {
			assert.Equal(t, name, game.Players()[i].Name)
			assert.Equal(t, "", game.Players()[i].AssignedRole)
		}
	})

	t.Run("error when adding duplicate player", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		names := []string{"Вася", "Петя", "Вася"} // Дубликат
		err := game.AddPlayers(names)
		assert.Error(t, err)
		assert.Len(t, game.Players(), 2) // Только первые два добавлены
	})
}

func TestGame_CopyPlayersFrom(t *testing.T) {
	t.Run("copy players from another game", func(t *testing.T) {
		// Создаем исходную игру с игроками
		sourceGame := NewGame("Source Game", 12345)
		err := sourceGame.AddPlayer("Вася", "townsfolk")
		require.NoError(t, err)
		err = sourceGame.AddPlayer("Петя", "demon")
		require.NoError(t, err)

		// Создаем целевую игру
		targetGame := NewGame("Target Game", 12345)

		// Копируем игроков
		err = targetGame.CopyPlayersFrom(sourceGame)
		assert.NoError(t, err)
		assert.Len(t, targetGame.Players(), 2)

		// Проверяем, что имена скопированы, а роли - нет
		assert.Equal(t, "Вася", targetGame.Players()[0].Name)
		assert.Equal(t, "", targetGame.Players()[0].AssignedRole) // Роль не копируется
		assert.Equal(t, "Петя", targetGame.Players()[1].Name)
		assert.Equal(t, "", targetGame.Players()[1].AssignedRole) // Роль не копируется
	})

	t.Run("copy players to game with existing players", func(t *testing.T) {
		sourceGame := NewGame("Source Game", 12345)
		err := sourceGame.AddPlayer("Вася")
		require.NoError(t, err)

		targetGame := NewGame("Target Game", 12345)
		err = targetGame.AddPlayer("Петя")
		require.NoError(t, err)

		err = targetGame.CopyPlayersFrom(sourceGame)
		assert.NoError(t, err)
		assert.Len(t, targetGame.Players(), 2) // Оба игрока
	})
}

func TestGame_SetAssignedRole(t *testing.T) {
	t.Run("set assigned role to player", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		playerID := game.Players()[0].ID
		err = game.SetAssignedRole(playerID, "demon")

		assert.NoError(t, err)
		assert.Equal(t, "demon", game.Players()[0].AssignedRole)
	})

	t.Run("error when setting empty role", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		playerID := game.Players()[0].ID
		err = game.SetAssignedRole(playerID, "")

		assert.Error(t, err)
		assert.Equal(t, "", game.Players()[0].AssignedRole)
	})
}

func TestGame_GetPlayerNames(t *testing.T) {
	game := NewGame("Test Game", 12345)
	err := game.AddPlayer("Вася")
	require.NoError(t, err)
	err = game.AddPlayer("Петя")
	require.NoError(t, err)
	err = game.AddPlayer("Коля")
	require.NoError(t, err)

	names := game.GetPlayerNames()
	assert.Equal(t, []string{"Вася", "Петя", "Коля"}, names)
}

func TestGame_HasPlayersWithRoles(t *testing.T) {
	t.Run("game without roles", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Вася")
		require.NoError(t, err)
		err = game.AddPlayer("Петя")
		require.NoError(t, err)

		assert.False(t, game.HasPlayersWithRoles())
	})

	t.Run("game with some roles", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Вася")
		require.NoError(t, err)
		err = game.AddPlayer("Петя", "townsfolk")
		require.NoError(t, err)

		assert.True(t, game.HasPlayersWithRoles())
	})

	t.Run("game with all roles", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer("Вася", "townsfolk")
		require.NoError(t, err)
		err = game.AddPlayer("Петя", "demon")
		require.NoError(t, err)

		assert.True(t, game.HasPlayersWithRoles())
	})
}
