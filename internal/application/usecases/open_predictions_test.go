package usecases

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenPredictionsUseCase(t *testing.T) {
	// Create all repositories
	gameRepo := persistence.NewInMemoryGameRepository()

	// Create a use case
	useCase := NewOpenPredictionsUseCase(gameRepo)

	t.Run("successful opening of predictions", func(t *testing.T) {
		// Creating a game in the created state
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.OpenPredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)
		assert.Equal(t, string(game.ID()), response.ID)
		assert.Equal(t, entities.GameStatusPredictionsOpen, response.Status)
	})

	t.Run("Error opening predictions for a non-existent games", func(t *testing.T) {
		command := dto.OpenPredictionsCommand{
			GameID:  "non-existent-game",
			AdminID: 12345,
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("error when opening predictions without players", func(t *testing.T) {
		game := entities.NewGame("Empty Game", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.OpenPredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("Error opening predictions by someone other than the game creator", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345) // creator ID 12345
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.OpenPredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 99999, // other user
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("error opening predictions for a game in an invalid state", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		// Opening predictions
		err = game.OpenPredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.OpenPredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err) // Already open
	})

	t.Run("error opening predictions for a completed game", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		// Simulate the full game cycle
		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = game.Start()
		require.NoError(t, err)

		err = game.Finish()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.OpenPredictionsCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err) // Game over
	})
}
