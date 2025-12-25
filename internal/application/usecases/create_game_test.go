package usecases_test

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/usecases"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGameUseCase(t *testing.T) {
	gameRepo := persistence.NewInMemoryGameRepository()
	useCase := usecases.NewCreateGameUseCase(gameRepo)

	t.Run("game creation successful", func(t *testing.T) {
		command := dto.CreateGameCommand{
			Name:      "Test Game",
			CreatorID: 12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)
		assert.Equal(t, "Test Game", response.Name)
		assert.Equal(t, int64(12345), response.CreatorID)
		assert.Equal(t, "created", string(response.Status))
	})

	t.Run("error with empty game name", func(t *testing.T) {
		command := dto.CreateGameCommand{
			Name:      "",
			CreatorID: 12345,
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("Error with invalid creator ID", func(t *testing.T) {
		command := dto.CreateGameCommand{
			Name:      "Test Game",
			CreatorID: 0,
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("Creating a game with spaces in the name", func(t *testing.T) {
		command := dto.CreateGameCommand{
			Name:      " Test Game ",
			CreatorID: 12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)
		assert.Equal(t, " Test Game ", response.Name) // Spaces are preserved
	})
}
