package services_test

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"

	"github.com/stretchr/testify/assert"
)

func TestBasicScoringRules_CalculatePointsForUser(t *testing.T) {
	scoringService := services.NewBasicScoringRules()

	t.Run("подсчёт очков за демонов", func(t *testing.T) {
		predictions := []entities.Prediction{
			*entities.NewPrediction("game-123", "user-1", "slot1", "demon"),
			*entities.NewPrediction("game-123", "user-1", "slot2", "minion"),
			*entities.NewPrediction("game-123", "user-1", "slot3", "townsfolk"),
		}

		realRoles := map[string]string{
			"slot1": "demon",     // угадал демона
			"slot2": "townsfolk", // не угадал приспешника
			"slot3": "townsfolk", // угадал мирного
		}

		totalScore := scoringService.CalculatePointsForUser(predictions, realRoles)

		// Expected: +10 for a demon, -2 for a minion's incorrect prediction, +2 for a peaceful one = 10
		assert.Equal(t, 10, totalScore)
	})

	t.Run("все роли угаданы", func(t *testing.T) {
		predictions := []entities.Prediction{
			*entities.NewPrediction("game-123", "user-1", "slot1", "demon"),
			*entities.NewPrediction("game-123", "user-1", "slot2", "minion"),
			*entities.NewPrediction("game-123", "user-1", "slot3", "townsfolk"),
		}

		realRoles := map[string]string{
			"slot1": "demon",
			"slot2": "minion",
			"slot3": "townsfolk",
		}

		totalScore := scoringService.CalculatePointsForUser(predictions, realRoles)

		// +10 for demon, +5 for minion, +2 for townsfolk = 17
		assert.Equal(t, 17, totalScore)
	})
}

func TestBasicScoringRules_CalculatePointsForEachPrediction(t *testing.T) {
	scoringService := services.NewBasicScoringRules()

	t.Run("расчёт очков для каждого прогноза", func(t *testing.T) {
		prediction1 := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		prediction2 := entities.NewPrediction("game-123", "user-1", "slot2", "minion")
		prediction3 := entities.NewPrediction("game-123", "user-1", "slot3", "townsfolk")

		predictions := []entities.Prediction{*prediction1, *prediction2, *prediction3}

		realRoles := map[string]string{
			"slot1": "demon",
			"slot2": "townsfolk",
			"slot3": "townsfolk",
		}

		pointsPerPrediction := scoringService.CalculatePointsForEachPrediction(predictions, realRoles)

		assert.Equal(t, 10, pointsPerPrediction[prediction1.ID()])
		assert.Equal(t, -2, pointsPerPrediction[prediction2.ID()])
		assert.Equal(t, 2, pointsPerPrediction[prediction3.ID()])
	})
}

func TestPrediction_PointsAwarded(t *testing.T) {
	t.Run("прогноз без начисленных очков", func(t *testing.T) {
		prediction := entities.NewPrediction("game-123", "user-456", "slot-789", "demon")

		points, awarded := prediction.PointsAwarded()
		assert.False(t, awarded)
		assert.Equal(t, 0, points)
		assert.False(t, prediction.HasPointsAwarded())
	})

	t.Run("начисление очков прогнозу", func(t *testing.T) {
		prediction := entities.NewPrediction("game-123", "user-456", "slot-789", "demon")

		err := prediction.AwardPoints(10)
		assert.NoError(t, err)

		points, awarded := prediction.PointsAwarded()
		assert.True(t, awarded)
		assert.Equal(t, 10, points)
		assert.True(t, prediction.HasPointsAwarded())
	})

	t.Run("нельзя начислить очки дважды", func(t *testing.T) {
		prediction := entities.NewPrediction("game-123", "user-456", "slot-789", "demon")

		err := prediction.AwardPoints(10)
		assert.NoError(t, err)

		err = prediction.AwardPoints(5)
		assert.Error(t, err)

		points, awarded := prediction.PointsAwarded()
		assert.True(t, awarded)
		assert.Equal(t, 10, points) // The first points awarded remain
	})
}
