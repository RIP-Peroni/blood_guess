package usecases

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
)

type TestFixture struct {
	GameRepo       *persistence.InMemoryGameRepository
	UserRepo       *persistence.InMemoryUserRepository
	PredictionRepo *persistence.InMemoryPredictionRepository
	Game           *entities.Game
	User           *entities.User
}

func SetupTestFixture() *TestFixture {
	gameRepo := persistence.NewInMemoryGameRepository()
	userRepo := persistence.NewInMemoryUserRepository()
	predictionRepo := persistence.NewInMemoryPredictionRepository()

	game := entities.NewGame("Test Game", 12345)
	_ = game.AddPlayer("Player 1", "townsfolk")
	_ = game.AddPlayer("Player 2", "outsider")
	_ = game.OpenPredictions()
	_ = gameRepo.Save(game)

	user := entities.NewUser(67890, "test_user")
	_ = userRepo.Save(user)

	return &TestFixture{
		GameRepo:       gameRepo,
		UserRepo:       userRepo,
		PredictionRepo: predictionRepo,
		Game:           game,
		User:           user,
	}
}
