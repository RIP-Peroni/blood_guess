package usecases

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPlayersUseCase(t *testing.T) {
	gameRepo := persistence.NewInMemoryGameRepository()
	useCase := NewAddPlayersUseCase(gameRepo)

	t.Run("successfully add multiple players", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayersCommand{
			GameID:      string(game.ID()),
			PlayerNames: []string{"Вася", "Петя", "Коля"},
			AdminID:     12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)

		assert.Equal(t, 3, len(response.Players))
		assert.Equal(t, "Вася", response.Players[0].Name)
		assert.Equal(t, "Петя", response.Players[1].Name)
		assert.Equal(t, "Коля", response.Players[2].Name)
	})

	t.Run("error adding duplicate players", func(t *testing.T) {
		game := entities.NewGame("Test Game 2", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayersCommand{
			GameID:      string(game.ID()),
			PlayerNames: []string{"Вася", "Петя", "Вася"}, // Дубликат
			AdminID:     12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists in this game")
	})

	t.Run("error when not creator", func(t *testing.T) {
		game := entities.NewGame("Test Game 3", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayersCommand{
			GameID:      string(game.ID()),
			PlayerNames: []string{"Вася", "Петя"},
			AdminID:     99999, // Не создатель
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only game creator")
	})

	t.Run("error when game not in created state", func(t *testing.T) {
		game := entities.NewGame("Test Game 4", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)
		err = game.OpenPredictions()
		require.NoError(t, err)
		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayersCommand{
			GameID:      string(game.ID()),
			PlayerNames: []string{"Вася", "Петя"},
			AdminID:     12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid game state")
	})
}
