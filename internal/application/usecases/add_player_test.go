// Файл: ./internal/application/usecases/add_player_test.go
package usecases

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddPlayerUseCase(t *testing.T) {
	gameRepo := persistence.NewInMemoryGameRepository()
	useCase := NewAddPlayerUseCase(gameRepo)

	t.Run("successfully add player to game", func(t *testing.T) {
		// Create a game
		game := entities.NewGame("Test Game", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "Alice",
			AdminID:    12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)

		assert.Equal(t, 1, len(response.Players))
		assert.Equal(t, "Alice", response.Players[0].Name)
	})

	t.Run("error adding player with duplicate name", func(t *testing.T) {
		game := entities.NewGame("Test Game 2", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		// Adding the first player
		command1 := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "Bob",
			AdminID:    12345,
		}
		_, err = useCase.Execute(command1)
		require.NoError(t, err)

		// Trying to add a player with the same name
		command2 := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "Bob",
			AdminID:    12345,
		}
		_, err = useCase.Execute(command2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("error adding player to non-existent game", func(t *testing.T) {
		command := dto.AddPlayerCommand{
			GameID:     "non-existent",
			PlayerName: "Charlie",
			AdminID:    12345,
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("error adding player by non-creator", func(t *testing.T) {
		game := entities.NewGame("Test Game 3", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "David",
			AdminID:    99999, // Not a creator
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("error adding player to game not in created state", func(t *testing.T) {
		game := entities.NewGame("Test Game 4", 12345)
		// First, add a player so that you can open the predictions
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		// Now let's open the forecasts - it should work
		err = game.OpenPredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "Eve",
			AdminID:    12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not in 'created' state")
	})
	t.Run("successfully add player to game without role", func(t *testing.T) {
		// Create a game
		game := entities.NewGame("Test Game", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "Alice",
			AdminID:    12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)

		assert.Equal(t, 1, len(response.Players))
		assert.Equal(t, "Alice", response.Players[0].Name)
		assert.Equal(t, "", response.Players[0].RealRole)
	})

	t.Run("error adding player with duplicate name", func(t *testing.T) {
		game := entities.NewGame("Test Game 2", 12345)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		// Adding the first player
		command1 := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "Bob",
			AdminID:    12345,
		}
		_, err = useCase.Execute(command1)
		require.NoError(t, err)

		// Trying to add a player with the same name
		command2 := dto.AddPlayerCommand{
			GameID:     string(game.ID()),
			PlayerName: "Bob",
			AdminID:    12345,
		}
		_, err = useCase.Execute(command2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}
