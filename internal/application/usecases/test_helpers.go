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
