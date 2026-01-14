package usecases

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAwardPointsUseCase_WithUnsetRoles(t *testing.T) {
	// Создаем все репозитории
	gameRepo := persistence.NewInMemoryGameRepository()
	userRepo := persistence.NewInMemoryUserRepository()
	predictionRepo := persistence.NewInMemoryPredictionRepository()
	scoringService := services.NewBasicScoringRules()

	useCase := NewAwardPointsUseCase(
		gameRepo,
		userRepo,
		predictionRepo,
		scoringService,
	)

	// Создаем тестовые данные
	creatorID := int64(12345)
	userID := int64(67890)

	// Создаем игру
	game := entities.NewGame("Test Game", creatorID)
	err := game.AddPlayer("Player 1")
	require.NoError(t, err)
	err = game.AddPlayer("Player 2")
	require.NoError(t, err)
	err = game.AddPlayer("Player 3")
	require.NoError(t, err)

	// Проходим все стадии игры
	err = game.OpenPredictions()
	require.NoError(t, err)
	err = game.ClosePredictions()
	require.NoError(t, err)
	err = game.Start()
	require.NoError(t, err)
	err = game.Finish()
	require.NoError(t, err)

	// Устанавливаем только злые роли (не все)
	player1ID := game.Players()[0].ID
	err = game.SetPlayerRealRole(player1ID, "demon")
	require.NoError(t, err)

	err = gameRepo.Save(game)
	require.NoError(t, err)

	// Создаем пользователя
	user := entities.NewUser(userID, "testuser")
	err = userRepo.Save(user)
	require.NoError(t, err)

	// Создаем прогнозы
	player2ID := string(game.Players()[1].ID)

	// Прогноз 1: правильно угадал демона
	pred1, err := entities.NewPrediction(
		game.ID(),
		user.ID(),
		entities.PlayerSlotID(player1ID),
		"demon",
	)
	require.NoError(t, err)
	err = predictionRepo.Save(pred1)
	require.NoError(t, err)

	// Прогноз 2: неправильно угадал приспешника (игрок на самом деле мирный)
	pred2, err := entities.NewPrediction(
		game.ID(),
		user.ID(),
		entities.PlayerSlotID(player2ID),
		"minion",
	)
	require.NoError(t, err)
	err = predictionRepo.Save(pred2)
	require.NoError(t, err)

	// Прогноз 3: угадал, что игрок мирный (не делал прогноз на злую роль)
	// Не создаем прогноз для player3 - это значит, что пользователь считает его мирным

	// Выполняем начисление очков
	command := dto.AwardPointsCommand{
		GameID:  string(game.ID()),
		AdminID: creatorID,
	}

	response, err := useCase.Execute(command)
	require.NoError(t, err)

	// Проверяем результаты:
	// 1. За демона: +10 очков
	// 2. За неправильный прогноз на minion: -2 очка (игрок мирный)
	// 3. За игрока без прогноза: 0 очков
	// Итого: 10 - 2 = 8 очков
	expectedScore := 8 // 10 за демона, -2 за неправильный minion

	assert.Equal(t, expectedScore, response.UserScores[string(user.ID())])
}

func TestAwardPointsUseCase_NoRolesSet(t *testing.T) {
	// Создаем все репозитории
	gameRepo := persistence.NewInMemoryGameRepository()
	userRepo := persistence.NewInMemoryUserRepository()
	predictionRepo := persistence.NewInMemoryPredictionRepository()
	scoringService := services.NewBasicScoringRules()

	useCase := NewAwardPointsUseCase(
		gameRepo,
		userRepo,
		predictionRepo,
		scoringService,
	)

	// Создаем тестовые данные
	creatorID := int64(12345)

	// Создаем игру
	game := entities.NewGame("Test Game", creatorID)
	err := game.AddPlayer("Player 1")
	require.NoError(t, err)

	// Проходим все стадии игры
	err = game.OpenPredictions()
	require.NoError(t, err)
	err = game.ClosePredictions()
	require.NoError(t, err)
	err = game.Start()
	require.NoError(t, err)
	err = game.Finish()
	require.NoError(t, err)

	// НЕ устанавливаем никакие роли

	err = gameRepo.Save(game)
	require.NoError(t, err)

	// Выполняем начисление очков
	command := dto.AwardPointsCommand{
		GameID:  string(game.ID()),
		AdminID: creatorID,
	}

	_, err = useCase.Execute(command)
	// Должна быть ошибка, так как нет ни одной установленной роли
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "не установлено ни одной реальной роли")
}
