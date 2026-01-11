package services_test

import (
	"RIP-Peroni/blood_guess/internal/domain/constants"
	"testing"

	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"

	"github.com/stretchr/testify/assert"
)

func TestBasicScoringRules_CalculatePointsForUser(t *testing.T) {
	scoringService := services.NewBasicScoringRules()

	t.Run("correctly guessed demon and minion", func(t *testing.T) {
		p1, _ := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		p2, _ := entities.NewPrediction("game-123", "user-1", "slot2", "minion")
		predictions := []*entities.Prediction{p1, p2}

		realRoles := map[string]string{
			"slot1": "demon",  // Угадал демона: +10
			"slot2": "minion", // Угадал приспешника: +5
		}

		totalScore := scoringService.CalculatePointsForUser(predictions, realRoles)
		expected := constants.PointsForDemon + constants.PointsForMinion
		assert.Equal(t, expected, totalScore)
	})

	t.Run("incorrect predictions for evil roles", func(t *testing.T) {
		p1, _ := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		p2, _ := entities.NewPrediction("game-123", "user-1", "slot2", "minion")
		predictions := []*entities.Prediction{p1, p2}

		realRoles := map[string]string{
			"slot1": "minion", // Предсказал demon, реальность minion: -3
			"slot2": "demon",  // Предсказал minion, реальность demon: -2
		}

		totalScore := scoringService.CalculatePointsForUser(predictions, realRoles)
		expected := -(constants.PenaltyForDemon + constants.PenaltyForMinion)
		assert.Equal(t, expected, totalScore)
	})

	t.Run("incorrect predictions when real role is peaceful", func(t *testing.T) {
		p1, _ := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		p2, _ := entities.NewPrediction("game-123", "user-1", "slot2", "minion")
		predictions := []*entities.Prediction{p1, p2}

		realRoles := map[string]string{
			"slot1": "townsfolk", // Предсказал demon, реальность townsfolk: -3
			"slot2": "outsider",  // Предсказал minion, реальность outsider: -2
		}

		totalScore := scoringService.CalculatePointsForUser(predictions, realRoles)
		assert.Equal(t, -5, totalScore) // -3 - 2 = -5
	})

	t.Run("mixed correct and incorrect predictions", func(t *testing.T) {
		p1, _ := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		p2, _ := entities.NewPrediction("game-123", "user-1", "slot2", "minion")
		p3, _ := entities.NewPrediction("game-123", "user-1", "slot3", "demon")
		predictions := []*entities.Prediction{p1, p2, p3}

		realRoles := map[string]string{
			"slot1": "demon",     // Угадал демона: +10
			"slot2": "townsfolk", // Предсказал minion, реальность townsfolk: -2
			"slot3": "minion",    // Предсказал demon, реальность minion: -3
		}

		totalScore := scoringService.CalculatePointsForUser(predictions, realRoles)
		assert.Equal(t, 5, totalScore) // 10 - 2 - 3 = 5
	})
}

func TestBasicScoringRules_CalculatePointsForEachPrediction(t *testing.T) {
	scoringService := services.NewBasicScoringRules()

	t.Run("score calculation for each prediction", func(t *testing.T) {
		p1, _ := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		p2, _ := entities.NewPrediction("game-123", "user-1", "slot2", "minion")

		predictions := []*entities.Prediction{p1, p2}

		realRoles := map[string]string{
			"slot1": "demon",     // Угадал демона: +10
			"slot2": "townsfolk", // Предсказал minion, реальность townsfolk: -2
		}

		pointsPerPrediction := scoringService.CalculatePointsForEachPrediction(predictions, realRoles)

		assert.Equal(t, 10, pointsPerPrediction[p1.ID()])
		assert.Equal(t, -2, pointsPerPrediction[p2.ID()])
	})
}
