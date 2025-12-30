package usecases

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartGameUseCase(t *testing.T) {
	gameRepo := persistence.NewInMemoryGameRepository()
	useCase := NewStartGameUseCase(gameRepo)

	t.Run("successful game start", func(t *testing.T) {
		// Arrange
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1", "townsfolk")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		// Act
		command := dto.StartGameCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		response, err := useCase.Execute(command)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, string(game.ID()), response.ID)
		assert.Equal(t, entities.GameStatusInProgress, response.Status)

		// Verify game was updated
		updatedGame, err := gameRepo.FindByID(game.ID())
		require.NoError(t, err)
		assert.Equal(t, entities.GameStatusInProgress, updatedGame.Status())
		assert.NotNil(t, updatedGame.StartedAt())
	})

	t.Run("error: not game creator", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1", "townsfolk")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.StartGameCommand{
			GameID:  string(game.ID()),
			AdminID: 99999, // different user
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only game creator")
	})

	t.Run("error: predictions not closed", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1", "townsfolk")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)
		// Don't close predictions

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.StartGameCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid game state")
	})

	t.Run("error: game already in progress", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1", "townsfolk")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		// Manually start the game
		err = game.Start()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.StartGameCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid game state")
	})

	t.Run("error: game not found", func(t *testing.T) {
		command := dto.StartGameCommand{
			GameID:  "non-existent-game",
			AdminID: 12345,
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "game not found")
	})

	t.Run("error: invalid command", func(t *testing.T) {
		// Empty game ID
		command1 := dto.StartGameCommand{
			GameID:  "",
			AdminID: 12345,
		}
		_, err := useCase.Execute(command1)
		assert.Error(t, err)

		// Invalid admin ID
		command2 := dto.StartGameCommand{
			GameID:  "game-123",
			AdminID: 0,
		}
		_, err = useCase.Execute(command2)
		assert.Error(t, err)
	})
}
