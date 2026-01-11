package usecases

import (
	"testing"
	"time"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyPlayersUseCase(t *testing.T) {
	t.Run("successfully copy players from previous game", func(t *testing.T) {
		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewCopyPlayersUseCase(gameRepo)

		previousGame := entities.NewGame("Previous Game", 12345)
		err := previousGame.AddPlayer("Вася")
		require.NoError(t, err)
		err = previousGame.AddPlayer("Петя")
		require.NoError(t, err)
		err = gameRepo.Save(previousGame)
		require.NoError(t, err)

		newGame := entities.NewGame("New Game", 12345)
		err = gameRepo.Save(newGame)
		require.NoError(t, err)

		command := dto.CopyPlayersCommand{
			TargetGameID: string(newGame.ID()),
			AdminID:      12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)
		assert.Len(t, response.Players, 2)
		assert.Equal(t, "Вася", response.Players[0].Name)
		assert.Equal(t, "Петя", response.Players[1].Name)
	})

	t.Run("copy from most recent game", func(t *testing.T) {
		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewCopyPlayersUseCase(gameRepo)

		oldGame := entities.NewGame("Old Game", 12345)
		err := oldGame.AddPlayer("Старый игрок")
		require.NoError(t, err)
		err = gameRepo.Save(oldGame)
		require.NoError(t, err)

		// Ждем немного, чтобы время отличалось
		time.Sleep(1 * time.Millisecond)

		recentGame := entities.NewGame("Recent Game", 12345)
		err = recentGame.AddPlayer("Новый игрок")
		require.NoError(t, err)
		err = gameRepo.Save(recentGame)
		require.NoError(t, err)

		newGame := entities.NewGame("New Game", 12345)
		err = gameRepo.Save(newGame)
		require.NoError(t, err)

		command := dto.CopyPlayersCommand{
			TargetGameID: string(newGame.ID()),
			AdminID:      12345,
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)
		assert.Len(t, response.Players, 1)
		assert.Equal(t, "Новый игрок", response.Players[0].Name) // Копируется из самой последней игры
	})

	t.Run("error when no previous games", func(t *testing.T) {
		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewCopyPlayersUseCase(gameRepo)

		// Создаем новую игру без предыдущих
		newGame := entities.NewGame("New Game", 12345)
		err := gameRepo.Save(newGame)
		require.NoError(t, err)

		command := dto.CopyPlayersCommand{
			TargetGameID: string(newGame.ID()),
			AdminID:      12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no previous games found")
	})

	t.Run("error when not creator", func(t *testing.T) {
		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewCopyPlayersUseCase(gameRepo)

		previousGame := entities.NewGame("Previous Game", 12345)
		err := previousGame.AddPlayer("Игрок")
		require.NoError(t, err)
		err = gameRepo.Save(previousGame)
		require.NoError(t, err)

		newGame := entities.NewGame("New Game", 99999) // Другой создатель
		err = gameRepo.Save(newGame)
		require.NoError(t, err)

		command := dto.CopyPlayersCommand{
			TargetGameID: string(newGame.ID()),
			AdminID:      12345, // Не создатель новой игры
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only game creator")
	})

	t.Run("error when game not in created state", func(t *testing.T) {
		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewCopyPlayersUseCase(gameRepo)

		// Создаем предыдущую игру
		previousGame := entities.NewGame("Previous Game", 12345)
		err := previousGame.AddPlayer("Игрок")
		require.NoError(t, err)
		err = gameRepo.Save(previousGame)
		require.NoError(t, err)

		newGame := entities.NewGame("New Game", 12345)
		err = newGame.AddPlayer("Игрок 2")
		require.NoError(t, err)
		err = newGame.OpenPredictions()
		require.NoError(t, err)
		err = gameRepo.Save(newGame)
		require.NoError(t, err)

		command := dto.CopyPlayersCommand{
			TargetGameID: string(newGame.ID()),
			AdminID:      12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid game state")
	})
}
