package usecases

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetRealRoleUseCase(t *testing.T) {
	t.Run("successfully set evil role after game", func(t *testing.T) {
		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewSetRealRoleUseCase(gameRepo)

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

		playerID := string(game.Players()[0].ID)

		command := dto.SetRealRoleCommand{
			GameID:       string(game.ID()),
			PlayerSlotID: playerID,
			RealRole:     "demon",
			AdminID:      12345,
		}

		err = useCase.Execute(command)
		require.NoError(t, err)
	})

	t.Run("error when setting peaceful role", func(t *testing.T) {
		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewSetRealRoleUseCase(gameRepo)

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

		playerID := string(game.Players()[0].ID)

		command := dto.SetRealRoleCommand{
			GameID:       string(game.ID()),
			PlayerSlotID: playerID,
			RealRole:     "townsfolk", // Мирная роль - должна вызывать ошибку
			AdminID:      12345,
		}

		err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "для установки доступны только злые роли")
	})
}
