package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"RIP-Peroni/blood_guess/internal/domain/entities"
)

func TestResponse_Structures(t *testing.T) {
	t.Run("GameResponse structure", func(t *testing.T) {
		response := GameResponse{
			ID:        "game-123",
			Name:      "Test Game",
			Status:    entities.GameStatusCreated,
			CreatorID: 123,
			CreatedAt: time.Now(),
		}
		assert.Equal(t, "game-123", response.ID)
		assert.Equal(t, "Test Game", response.Name)
		assert.Equal(t, entities.GameStatusCreated, response.Status)
	})

	t.Run("PredictionResponse structure", func(t *testing.T) {
		response := PredictionResponse{
			ID:            "pred-123",
			GameID:        "game-123",
			UserID:        "user-456",
			PlayerSlotID:  "slot-789",
			PredictedRole: "demon",
			Points:        10,
			PointsAwarded: true,
			CreatedAt:     time.Now(),
		}
		assert.Equal(t, "pred-123", response.ID)
		assert.Equal(t, "demon", response.PredictedRole)
		assert.Equal(t, 10, response.Points)
		assert.True(t, response.PointsAwarded)
	})

	t.Run("UserResponse structure", func(t *testing.T) {
		response := UserResponse{
			ID:         "user-123",
			TelegramID: 12345,
			Username:   "test_user",
			Balance:    100,
			CreatedAt:  time.Now(),
		}
		assert.Equal(t, "user-123", response.ID)
		assert.Equal(t, int64(12345), response.TelegramID)
		assert.Equal(t, "test_user", response.Username)
		assert.Equal(t, 100, response.Balance)
	})

	t.Run("CalculateResultsResponse structure", func(t *testing.T) {
		response := CalculateResultsResponse{
			GameID: "game-123",
			Scores: map[string]int{
				"user-1": 15,
				"user-2": -5,
			},
		}
		assert.Equal(t, "game-123", response.GameID)
		assert.Equal(t, 15, response.Scores["user-1"])
		assert.Equal(t, -5, response.Scores["user-2"])
	})
}
