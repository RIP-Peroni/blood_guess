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
	t.Run("successfully set real role after game", func(t *testing.T) {
		t.Log("=== Test: успешная установка реальной роли после игры ===")

		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewSetRealRoleUseCase(gameRepo)

		// Создаем игру и переводим в состояние in_progress
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = game.Start() // Переводим в in_progress
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		playerID := string(game.Players()[0].ID)
		t.Logf("Игра создана, статус: %s, игрок ID: %s", game.Status(), playerID)

		command := dto.SetRealRoleCommand{
			GameID:       string(game.ID()),
			PlayerSlotID: playerID,
			RealRole:     "demon", // Устанавливаем реальную роль - демон
			AdminID:      12345,
		}

		err = useCase.Execute(command)
		require.NoError(t, err, "Не удалось установить реальную роль")

		// Проверяем, что реальная роль установлена
		updatedGame, err := gameRepo.FindByID(game.ID())
		require.NoError(t, err)

		player, err := updatedGame.FindPlayerByID(entities.PlayerSlotID(playerID))
		require.NoError(t, err)

		assert.Equal(t, "demon", player.RealRole, "Реальная роль не установлена")
		assert.True(t, player.IsRealRoleSet, "Флаг IsRealRoleSet не установлен")
		t.Logf("Реальная роль успешно установлена: %s", player.RealRole)
	})

	t.Run("error when setting role before game starts", func(t *testing.T) {
		t.Log("=== Test: ошибка при установке реальной роли до начала игры ===")

		gameRepo := persistence.NewInMemoryGameRepository()
		useCase := NewSetRealRoleUseCase(gameRepo)

		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		// Не начинаем игру - статус остается created или predictions_open
		err = game.OpenPredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		playerID := string(game.Players()[0].ID)
		t.Logf("Статус игры: %s", game.Status())

		command := dto.SetRealRoleCommand{
			GameID:       string(game.ID()),
			PlayerSlotID: playerID,
			RealRole:     "demon",
			AdminID:      12345,
		}

		err = useCase.Execute(command)
		assert.Error(t, err, "Ожидалась ошибка при установке роли до начала игры")
		assert.Contains(t, err.Error(), "can only set real roles during or after the game",
			"Сообщение об ошибке должно указывать на недопустимое время установки роли")
		t.Logf("Получена ожидаемая ошибка: %v", err)
	})
}
