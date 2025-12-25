package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand_Structures(t *testing.T) {
	t.Run("CreateGameCommand structure", func(t *testing.T) {
		cmd := CreateGameCommand{
			Name:      "Test",
			CreatorID: 123,
		}
		assert.Equal(t, "Test", cmd.Name)
		assert.Equal(t, int64(123), cmd.CreatorID)
	})
}

func TestCommand_Completeness(t *testing.T) {
	t.Run("SubmitPredictionCommand contains all the required fields", func(t *testing.T) {
		cmd := SubmitPredictionCommand{
			GameID:        "game-123",
			UserID:        "user-456",
			PlayerSlotID:  "slot-789",
			PredictedRole: "demon",
		}

		assert.NotEmpty(t, cmd.GameID)
		assert.NotEmpty(t, cmd.UserID)
		assert.NotEmpty(t, cmd.PlayerSlotID)
		assert.NotEmpty(t, cmd.PredictedRole)
	})

	t.Run("FinishGameCommand contains real roles", func(t *testing.T) {
		cmd := FinishGameCommand{
			GameID:  "game-123",
			AdminID: 12345,
			RealRoles: map[string]string{
				"slot-1": "demon",
				"slot-2": "minion",
			},
		}

		assert.Equal(t, "game-123", cmd.GameID)
		assert.Equal(t, int64(12345), cmd.AdminID)
		assert.Len(t, cmd.RealRoles, 2)
		assert.Equal(t, "demon", cmd.RealRoles["slot-1"])
	})
}
