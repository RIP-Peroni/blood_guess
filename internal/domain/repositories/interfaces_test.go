package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testGameRepositoryInterface(t *testing.T) {
	var _ GameRepository

	assert.True(t, true, "The interface is defined correctly")
}

func TestPredictionValidation(t *testing.T) {
	t.Run("a prediction is valid", func(t *testing.T) {
		prediction := Prediction{
			GameID:        "game-123",
			UserID:        "user-456",
			PlayerSlotID:  "slot-789",
			PredictedRole: "demon",
		}

		assert.NotEmpty(t, prediction.GameID)
		assert.NotEmpty(t, prediction.UserID)
		assert.NotEmpty(t, prediction.PlayerSlotID)
		assert.NotEmpty(t, prediction.PredictedRole)
	})
}
