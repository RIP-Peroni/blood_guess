// internal/domain/services/basic_scoring_test.go
package services_test

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"

	"github.com/stretchr/testify/assert"
)

func TestBasicScoringRules_CalculatePointsForUser(t *testing.T) {
	scoringService := services.NewBasicScoringRules()

	t.Run("demon scoring", func(t *testing.T) {
		p1, _ := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		p2, _ := entities.NewPrediction("game-123", "user-1", "slot2", "minion")
		p3, _ := entities.NewPrediction("game-123", "user-1", "slot3", "townsfolk")
		predictions := []*entities.Prediction{
			p1,
			p2,
			p3,
		}

		realRoles := map[string]string{
			"slot1": "demon",     // guessed the demon
			"slot2": "townsfolk", // didn't guess the minion
			"slot3": "townsfolk", // guessed peaceful
		}

		totalScore := scoringService.CalculatePointsForUser(predictions, realRoles)

		// Expected: +10 for a demon, -2 for a minion's incorrect prediction, +2 for a peaceful one = 10
		assert.Equal(t, 10, totalScore)
	})

	t.Run("all roles are guessed", func(t *testing.T) {
		p1, _ := entities.NewPrediction(entities.GameID("game-123"), entities.UserID("user-1"), entities.PlayerSlotID("slot1"), "demon")
		p2, _ := entities.NewPrediction(entities.GameID("game-123"), entities.UserID("user-1"), entities.PlayerSlotID("slot2"), "minion")
		p3, _ := entities.NewPrediction(entities.GameID("game-123"), entities.UserID("user-1"), entities.PlayerSlotID("slot3"), "townsfolk")
		predictions := []*entities.Prediction{
			p1,
			p2,
			p3,
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

	t.Run("score calculation for each prediction", func(t *testing.T) {
		p1, _ := entities.NewPrediction("game-123", "user-1", "slot1", "demon")
		p2, _ := entities.NewPrediction("game-123", "user-1", "slot2", "minion")
		p3, _ := entities.NewPrediction("game-123", "user-1", "slot3", "townsfolk")

		predictions := []*entities.Prediction{p1, p2, p3}

		realRoles := map[string]string{
			"slot1": "demon",
			"slot2": "townsfolk",
			"slot3": "townsfolk",
		}

		pointsPerPrediction := scoringService.CalculatePointsForEachPrediction(predictions, realRoles)

		assert.Equal(t, 10, pointsPerPrediction[p1.ID()])
		assert.Equal(t, -2, pointsPerPrediction[p2.ID()])
		assert.Equal(t, 2, pointsPerPrediction[p3.ID()])
	})
}
