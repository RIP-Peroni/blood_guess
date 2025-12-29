package usecases

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitPredictionUseCase(t *testing.T) {
	// Create all repositories
	gameRepo := persistence.NewInMemoryGameRepository()
	userRepo := persistence.NewInMemoryUserRepository()
	predictionRepo := persistence.NewInMemoryPredictionRepository()

	// Create use case
	useCase := NewSubmitPredictionUseCase(
		gameRepo,
		userRepo,
		predictionRepo,
	)

	// Create test data
	creatorID := int64(12345)
	userID := int64(67890)
	game := entities.NewGame("Test Game", creatorID)

	// Add players to the game
	err := game.AddPlayer("Player 1", "townsfolk")
	require.NoError(t, err)
	err = game.AddPlayer("Player 2", "outsider")
	require.NoError(t, err)

	// Opening the predictions
	err = game.OpenPredictions()
	require.NoError(t, err)

	// Save the game
	err = gameRepo.Save(game)
	require.NoError(t, err)

	// Create a user
	user := entities.NewUser(userID, "test_user")
	err = userRepo.Save(user)
	require.NoError(t, err)

	t.Run("successful creation of a prediction", func(t *testing.T) {
		// Get the ID of the player's first slot
		playerSlotID := string(game.Players()[0].ID)

		command := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  playerSlotID,
			PredictedRole: "demon",
		}

		response, err := useCase.Execute(command)
		require.NoError(t, err)
		assert.Equal(t, string(game.ID()), response.GameID)
		assert.Equal(t, string(user.ID()), response.UserID)
		assert.Equal(t, playerSlotID, response.PlayerSlotID)
		assert.Equal(t, "demon", response.PredictedRole)
		assert.False(t, response.PointsAwarded)
	})

	t.Run("error when creating a second prediction for the same slot", func(t *testing.T) {
		playerSlotID := string(game.Players()[0].ID)

		command := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  playerSlotID,
			PredictedRole: "minion",
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already predicted")
	})

	t.Run("prediction error for a non-existent game", func(t *testing.T) {
		command := dto.SubmitPredictionCommand{
			GameID:        "non-existent-game",
			UserID:        string(user.ID()),
			PlayerSlotID:  "slot-123",
			PredictedRole: "demon",
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("prediction error for a non-existent user", func(t *testing.T) {
		playerSlotID := string(game.Players()[1].ID)

		command := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        "non-existent-user",
			PlayerSlotID:  playerSlotID,
			PredictedRole: "demon",
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("prediction error for a non-existent player slot", func(t *testing.T) {
		command := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  "non-existent-slot",
			PredictedRole: "demon",
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("a prediction error when the game does not accept predictions", func(t *testing.T) {
		// Create a new game, but don't open the predictions
		closedGame := entities.NewGame("Closed Game", creatorID)
		err := closedGame.AddPlayer("Player 3", "townsfolk")
		require.NoError(t, err)

		err = gameRepo.Save(closedGame)
		require.NoError(t, err)

		playerSlotID := string(closedGame.Players()[0].ID)

		command := dto.SubmitPredictionCommand{
			GameID:        string(closedGame.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  playerSlotID,
			PredictedRole: "demon",
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
	})

	t.Run("invalid role error", func(t *testing.T) {
		playerSlotID := string(game.Players()[1].ID)

		command := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  playerSlotID,
			PredictedRole: "invalid_role",
		}

		_, err := useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
	})

	t.Run("Can create predictions for different slots", func(t *testing.T) {
		// Let's create a new user for this test.
		newUser := entities.NewUser(99999, "another_user")
		err := userRepo.Save(newUser)
		require.NoError(t, err)

		// Prediction for the first slot
		playerSlotID1 := string(game.Players()[0].ID)
		command1 := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(newUser.ID()),
			PlayerSlotID:  playerSlotID1,
			PredictedRole: "demon",
		}

		response1, err := useCase.Execute(command1)
		require.NoError(t, err)
		assert.Equal(t, playerSlotID1, response1.PlayerSlotID)

		// Prediction for the second slot
		playerSlotID2 := string(game.Players()[1].ID)
		command2 := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(newUser.ID()),
			PlayerSlotID:  playerSlotID2,
			PredictedRole: "minion",
		}

		response2, err := useCase.Execute(command2)
		require.NoError(t, err)
		assert.Equal(t, playerSlotID2, response2.PlayerSlotID)

		// We check that the user now has two predictions
		predictions, err := predictionRepo.FindByGameAndUser(
			game.ID(),
			newUser.ID(),
		)
		require.NoError(t, err)
		assert.Len(t, predictions, 2)
	})
}
