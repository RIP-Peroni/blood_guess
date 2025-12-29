package usecases

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClosePredictionsUseCase(t *testing.T) {
	// Create all repositories
	gameRepo := persistence.NewInMemoryGameRepository()

	// Create a use case
	useCase := NewClosePredictionsUseCase(gameRepo)

	t.Run("successful closing of predictions", func(t *testing.T) {
		// Creating a game with open predictions
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1", "townsfolk")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.ClosePredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)
		assert.Equal(t, string(game.ID()), response.ID)
		assert.Equal(t, entities.GameStatusPredictionsClosed, response.Status)
	})

	t.Run("Error closing predictions for a non-existent game", func(t *testing.T) {
		command := dto.ClosePredictionsCommand{
			GameID:  "non-existent-game",
			AdminID: 12345,
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("Error closing predictions for a non-game creator", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1", "townsfolk")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.ClosePredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 99999, // another user
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("error closing predictions for a game in an invalid state", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1", "townsfolk")
		require.NoError(t, err)

		// The game has been created, but predictions have not been opened
		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.ClosePredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err) // Predictions not open
	})
}
