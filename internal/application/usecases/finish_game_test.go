package usecases

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinishGameUseCase(t *testing.T) {
	// Create repository
	gameRepo := persistence.NewInMemoryGameRepository()

	// Create a use case
	useCase := NewFinishGameUseCase(gameRepo)

	t.Run("successful game finish", func(t *testing.T) {
		// Arrange: create test data
		creatorID := int64(12345)

		// Create a game
		game := entities.NewGame("Test Game", creatorID)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)
		err = game.AddPlayer("Player 2")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = game.Start()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		// Act: Finish the game
		command := dto.FinishGameCommand{
			GameID:  string(game.ID()),
			AdminID: creatorID,
		}

		response, err := useCase.Execute(command)

		// Assert: check the results
		require.NoError(t, err)
		assert.Equal(t, string(game.ID()), response.GameID)
		assert.Equal(t, entities.GameStatusFinished, response.Status)
		assert.False(t, response.CurrencyAwarded) // Points are not awarded by finish

		// Checking that the game is finished
		updatedGame, err := gameRepo.FindByID(game.ID())
		require.NoError(t, err)
		assert.Equal(t, entities.GameStatusFinished, updatedGame.Status())
		assert.NotNil(t, updatedGame.EndedAt())
	})

	t.Run("error: not game creator", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = game.Start()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.FinishGameCommand{
			GameID:  string(game.ID()),
			AdminID: 99999, // other user
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only game creator")
	})

	t.Run("error: game not in progress", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		// We don't start the game!

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.FinishGameCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid game state")
	})
}
