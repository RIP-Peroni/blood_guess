package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicScoringRules(t *testing.T) {
	scoringService := NewBasicScoringRules()

	t.Run("scoring for demons", func(t *testing.T) {
		predictions := []Prediction{
			{PlayerSlotID: "slot1", PredictedRole: "demon"},
			{PlayerSlotID: "slot2", PredictedRole: "minion"},
			{PlayerSlotID: "slot3", PredictedRole: "townsfolk"},
		}

		realRoles := map[string]string{
			"slot1": "demon",     // угадал демона
			"slot2": "townsfolk", // не угадал приспешника
			"slot3": "townsfolk", // угадал мирного
		}

		scores := scoringService.Calculate(predictions, realRoles)

		// Ожидаем: +10 за демона, +0 за неверный прогноз, +2 за мирного
		assert.Equal(t, 12, scores)
	})

	t.Run("all roles have been predicted successfully", func(t *testing.T) {
		predictions := []Prediction{
			{PlayerSlotID: "slot1", PredictedRole: "demon"},
			{PlayerSlotID: "slot2", PredictedRole: "minion"},
			{PlayerSlotID: "slot3", PredictedRole: "townsfolk"},
		}

		realRoles := map[string]string{
			"slot1": "demon",
			"slot2": "minion",
			"slot3": "townsfolk",
		}

		scores := scoringService.Calculate(predictions, realRoles)

		// +10 за демона, +5 за приспешника, +2 за мирного = 17
		assert.Equal(t, 17, scores)
	})

	t.Run("nothing is predicted successfully", func(t *testing.T) {
		predictions := []Prediction{
			{PlayerSlotID: "slot1", PredictedRole: "demon"},
			{PlayerSlotID: "slot2", PredictedRole: "minion"},
		}

		realRoles := map[string]string{
			"slot1": "townsfolk",
			"slot2": "demon",
		}

		scores := scoringService.Calculate(predictions, realRoles)

		// Все прогнозы неверные = 0 очков
		assert.Equal(t, 0, scores)
	})

	t.Run("penalty for the demon's incorrect prediction", func(t *testing.T) {
		predictions := []Prediction{
			{PlayerSlotID: "slot1", PredictedRole: "demon"},     // неверно
			{PlayerSlotID: "slot2", PredictedRole: "townsfolk"}, // верно
		}

		realRoles := map[string]string{
			"slot1": "townsfolk",
			"slot2": "townsfolk",
		}

		scores := scoringService.Calculate(predictions, realRoles)

		// +2 за мирного, -3 за неверный демон = -1
		assert.Equal(t, -1, scores)
	})
}
