package usecases

import (
	"testing"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinishGameUseCase(t *testing.T) {
	// Create all repositories
	gameRepo := persistence.NewInMemoryGameRepository()
	userRepo := persistence.NewInMemoryUserRepository()
	predictionRepo := persistence.NewInMemoryPredictionRepository()

	// Create a scoring service
	scoringService := services.NewBasicScoringRules()

	// Create a use case
	useCase := NewFinishGameUseCase(gameRepo, userRepo, predictionRepo, scoringService)

	t.Run("successful game finish with points calculation", func(t *testing.T) {
		// Arrange: create test data
		creatorID := int64(12345)
		user1ID := int64(67890)
		user2ID := int64(99999)

		// Create a game
		game := entities.NewGame("Test Game", creatorID)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)
		err = game.AddPlayer("Player 2")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = game.Start()
		require.NoError(t, err)

		// Set up real roles
		player1ID := game.Players()[0].ID
		player2ID := game.Players()[1].ID
		err = game.SetPlayerRealRole(player1ID, "demon") // Player 1 turned out to be a demon
		require.NoError(t, err)
		err = game.SetPlayerRealRole(player2ID, "townsfolk") // Player 2 turned out to be a townsfolk
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		// Create users
		user1 := entities.NewUser(user1ID, "user1")
		err = userRepo.Save(user1)
		require.NoError(t, err)

		user2 := entities.NewUser(user2ID, "user2")
		err = userRepo.Save(user2)
		require.NoError(t, err)

		// Generating predictions
		// User 1: correctly guessed the demon, incorrectly guessed the townsperson as minion
		pred1, err := entities.NewPrediction(
			game.ID(),
			user1.ID(),
			player1ID,
			"demon", // guessed the demon correctly
		)
		require.NoError(t, err)
		err = predictionRepo.Save(pred1)
		require.NoError(t, err)

		pred2, err := entities.NewPrediction(
			game.ID(),
			user1.ID(),
			player2ID,
			"minion", // Made a mistake, thought he was a minion, but it turned out to be a townsfolk
		)
		require.NoError(t, err)
		err = predictionRepo.Save(pred2)
		require.NoError(t, err)

		// User 2: Guessed everything wrong (только злые роли!)
		pred3, err := entities.NewPrediction(
			game.ID(),
			user2.ID(),
			player1ID,
			"minion", // mistake - thought minion, but turned out to be a demon
		)
		require.NoError(t, err)
		err = predictionRepo.Save(pred3)
		require.NoError(t, err)

		pred4, err := entities.NewPrediction(
			game.ID(),
			user2.ID(),
			player2ID,
			"demon", // mistake - thought demon, but turned out to be a townsfolk
		)
		require.NoError(t, err)
		err = predictionRepo.Save(pred4)
		require.NoError(t, err)

		// Act: End the game
		command := dto.FinishGameCommand{
			GameID:  string(game.ID()),
			AdminID: creatorID,
		}

		response, err := useCase.Execute(command)

		// Assert: check the results
		require.NoError(t, err)
		assert.Equal(t, string(game.ID()), response.GameID)
		assert.Equal(t, entities.GameStatusFinished, response.Status)
		assert.True(t, response.CurrencyAwarded)

		// Checking the points awarded
		assert.Equal(t, 8, response.UserScores[string(user1.ID())])  // +10 per demon, -2 per minion = 8
		assert.Equal(t, -5, response.UserScores[string(user2.ID())]) // -2 for wrong minion, -3 for wrong demon = -5

		updatedUser1, err := userRepo.FindById(user1.ID())
		require.NoError(t, err)
		assert.Equal(t, 8, updatedUser1.Balance()) // 8 points

		updatedUser2, err := userRepo.FindById(user2.ID())
		require.NoError(t, err)
		assert.Equal(t, 0, updatedUser2.Balance()) // -5, but the balance cannot be negative, so 0

		// Check that the predictions received points
		updatedPred1, err := predictionRepo.FindByID(pred1.ID())
		require.NoError(t, err)
		points1, awarded1 := updatedPred1.PointsAwarded()
		assert.True(t, awarded1)
		assert.Equal(t, 10, points1)

		// Checking that the game is finished
		updatedGame, err := gameRepo.FindByID(game.ID())
		require.NoError(t, err)
		assert.Equal(t, entities.GameStatusFinished, updatedGame.Status())
		assert.NotNil(t, updatedGame.EndedAt())
	})

	t.Run("error: not all real roles set", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = game.Start()
		require.NoError(t, err)

		// We do not set real roles!

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.FinishGameCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "real roles")
	})

	t.Run("error: not game creator", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = game.Start()
		require.NoError(t, err)

		err = game.SetPlayerRealRole(game.Players()[0].ID, "demon")
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.FinishGameCommand{
			GameID:  string(game.ID()),
			AdminID: 99999, // other user
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only game creator")
	})

	t.Run("error: game not in progress", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		// We don't start the game!

		err = gameRepo.Save(game)
		require.NoError(t, err)

		command := dto.FinishGameCommand{
			GameID:  string(game.ID()),
			AdminID: 12345,
		}

		_, err = useCase.Execute(command)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid game state")
	})
}
