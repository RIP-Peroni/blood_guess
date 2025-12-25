package entities

import (
	"RIP-Peroni/blood_guess/internal/domain/value_objects"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrediction_Entity(t *testing.T) {
	t.Run("creating prediction", func(t *testing.T) {
		gameID := GameID("game-123")
		userID := UserID("user-456")
		slotID := PlayerSlotID("player-789")

		prediction, _ := NewPrediction(gameID, userID, slotID, "demon")

		assert.Equal(t, gameID, prediction.GameID())
		assert.Equal(t, userID, prediction.UserID())
		assert.Equal(t, slotID, prediction.PlayerSlotID())
		assert.Equal(t, value_objects.Role("demon"), prediction.PredictedRole())
		assert.Nil(t, prediction.pointsAwarded)
		assert.WithinDuration(t, time.Now(), prediction.CreatedAt(), time.Second)
	})
	t.Run("verification of the prediction's correctness", func(t *testing.T) {
		prediction, _ := NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)

		assert.True(t, prediction.IsCorrect("demon"))
		assert.False(t, prediction.IsCorrect("minion"))
		assert.False(t, prediction.IsCorrect(""))
	})
	t.Run("scoring", func(t *testing.T) {
		prediction, _ := NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)

		pointsAwarded, ok := prediction.PointsAwarded()
		assert.False(t, ok)
		assert.Equal(t, 0, pointsAwarded)

		err := prediction.AwardPoints(10)
		assert.NoError(t, err)
		pointsAwarded, ok = prediction.PointsAwarded()
		assert.True(t, ok)
		assert.Equal(t, 10, pointsAwarded)

		err = prediction.AwardPoints(-5)
		assert.Error(t, err)
	})
}
