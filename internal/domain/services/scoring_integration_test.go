package services_test

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"

	"github.com/stretchr/testify/assert"
)

func TestScoringIntegration_RealWorldScenario(t *testing.T) {
	scoringService := services.NewBasicScoringRules()

	// Сценарий: 5 игроков, 2 пользователя делают прогнозы
	gameID := "game-123"

	// Создаём прогнозы пользователя 1
	p1, _ := entities.NewPrediction(entities.GameID(gameID), "user-1", "player-1", "demon")  // Правильно
	p2, _ := entities.NewPrediction(entities.GameID(gameID), "user-1", "player-2", "minion") // Правильно
	p3, _ := entities.NewPrediction(entities.GameID(gameID), "user-1", "player-3", "demon")  // Неправильно (мирный)

	// Создаём прогнозы пользователя 2
	p4, _ := entities.NewPrediction(entities.GameID(gameID), "user-2", "player-1", "minion") // Неправильно (demon)
	p5, _ := entities.NewPrediction(entities.GameID(gameID), "user-2", "player-2", "demon")  // Неправильно (minion)
	p6, _ := entities.NewPrediction(entities.GameID(gameID), "user-2", "player-4", "minion") // Неправильно (мирный)

	predictionsUser1 := []*entities.Prediction{p1, p2, p3}
	predictionsUser2 := []*entities.Prediction{p4, p5, p6}

	// Реальные роли после игры:
	// player-1: demon
	// player-2: minion
	// player-3: townsfolk
	// player-4: outsider
	// player-5: townsfolk (никто не предсказывал)
	realRoles := map[string]string{
		"player-1": "demon",
		"player-2": "minion",
		"player-3": "townsfolk",
		"player-4": "outsider",
		"player-5": "townsfolk",
	}

	// Подсчитываем очки для пользователя 1:
	// p1: предсказал demon для player-1, реальность demon = +10
	// p2: предсказал minion для player-2, реальность minion = +5
	// p3: предсказал demon для player-3, реальность townsfolk = -3
	// Итого: 10 + 5 - 3 = 12
	scoreUser1 := scoringService.CalculatePointsForUser(predictionsUser1, realRoles)
	assert.Equal(t, 12, scoreUser1, "Неправильный счёт для пользователя 1")

	// Подсчитываем очки для пользователя 2:
	// p4: предсказал minion для player-1, реальность demon = -2 (minion penalty)
	// p5: предсказал demon для player-2, реальность minion = -3 (demon penalty)
	// p6: предсказал minion для player-4, реальность outsider = -2 (minion penalty)
	// Итого: -2 - 3 - 2 = -7
	scoreUser2 := scoringService.CalculatePointsForUser(predictionsUser2, realRoles)
	assert.Equal(t, -7, scoreUser2, "Неправильный счёт для пользователя 2")
}
