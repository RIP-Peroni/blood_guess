package persistence_test

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryPredictionRepository(t *testing.T) {
	t.Run("saving and getting the predictions", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		pred, _ := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)

		err := repo.Save(pred)
		require.NoError(t, err)

		found, err := repo.FindByID(pred.ID())
		require.NoError(t, err)
		assert.Equal(t, pred.ID(), found.ID())
		assert.Equal(t, pred.GameID(), found.GameID())
		assert.Equal(t, pred.PredictedRole(), found.PredictedRole())
	})

	t.Run("search for predictions by game and user", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		gameID := entities.GameID("game-123")
		userID := entities.UserID("user-456")

		pred1, _ := entities.NewPrediction(gameID, userID, "slot-1", "demon")
		pred2, _ := entities.NewPrediction(gameID, userID, "slot-2", "minion")

		err := repo.Save(pred1)
		require.NoError(t, err)
		err = repo.Save(pred2)
		require.NoError(t, err)

		predictions, err := repo.FindByGameAndUser(gameID, userID)
		require.NoError(t, err)
		assert.Len(t, predictions, 2)
	})

	t.Run("predictions update", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		pred, _ := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)
		err := repo.Save(pred)
		require.NoError(t, err)

		// Award points
		err = pred.AwardPoints(10)
		require.NoError(t, err)

		err = repo.Update(pred)
		require.NoError(t, err)

		updated, err := repo.FindByID(pred.ID())
		require.NoError(t, err)

		points, awarded := updated.PointsAwarded()
		assert.True(t, awarded)
		assert.Equal(t, 10, points)
	})

	t.Run("delete predictions", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		pred, _ := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)
		err := repo.Save(pred)
		require.NoError(t, err)

		err = repo.Delete(pred.ID())
		require.NoError(t, err)

		_, err = repo.FindByID(pred.ID())
		assert.Error(t, err)
	})

	t.Run("search for a non-existent predictions", func(t *testing.T) {
		repo := persistence.NewInMemoryPredictionRepository()
		_, err := repo.FindByID("non-existent")
		assert.Error(t, err)

		_, err = repo.FindByGameAndUser("game-123", "user-999")
		require.NoError(t, err)
	})
}
