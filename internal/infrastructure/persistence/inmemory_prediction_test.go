package persistence_test

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryPredictionRepository(t *testing.T) {
	t.Run("saving and getting the predictions", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		prediction, _ := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)

		err := repo.Save(prediction)
		require.NoError(t, err)

		found, err := repo.FindByID(prediction.ID())
		require.NoError(t, err)
		assert.Equal(t, prediction.ID(), found.ID())
		assert.Equal(t, prediction.GameID(), found.GameID())
		assert.Equal(t, prediction.PredictedRole(), found.PredictedRole())
	})

	t.Run("search for predictions by game and user", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		gameID := entities.GameID("game-123")
		userID := entities.UserID("user-456")

		prediction1, _ := entities.NewPrediction(gameID, userID, "slot-1", "demon")
		prediction2, _ := entities.NewPrediction(gameID, userID, "slot-2", "minion")

		_ = repo.Save(prediction1)
		_ = repo.Save(prediction2)

		predictions, err := repo.FindByGameAndUser(gameID, userID)
		require.NoError(t, err)
		fmt.Printf("predictions: %+v\n", predictions)
		assert.Len(t, predictions, 2)
	})

	t.Run("predictions update", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		prediction, _ := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)
		_ = repo.Save(prediction)

		//Award points
		err := prediction.AwardPoints(10)
		require.NoError(t, err)

		err = repo.Update(prediction)
		require.NoError(t, err)

		updated, err := repo.FindByID(prediction.ID())
		require.NoError(t, err)

		points, awarded := updated.PointsAwarded()
		assert.True(t, awarded)
		assert.Equal(t, 10, points)
	})

	t.Run("delete predictions", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		prediction, _ := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)
		_ = repo.Save(prediction)

		err := repo.Delete(prediction.ID())
		require.NoError(t, err)

		_, err = repo.FindByID(prediction.ID())
		assert.Error(t, err)
	})

	t.Run("search for a non-existent predictions", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		_, err := repo.FindByID(entities.PredictionID("non-existent"))
		assert.Error(t, err)

		_, err = repo.FindByGameAndUser("game-123", "user-999")
		require.NoError(t, err) // Returns an empty list, not an error
	})
}
