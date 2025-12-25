package persistence

// UnitOfWork provides access to all repositories
type UnitOfWork struct {
	GameRepo       *InMemoryGameRepository
	UserRepo       *InMemoryUserRepository
	PredictionRepo *InMemoryPredictionRepository
}

func NewUnitOfWork() *UnitOfWork {
	return &UnitOfWork{
		GameRepo:       NewInMemoryGameRepository(),
		UserRepo:       NewInMemoryUserRepository(),
		PredictionRepo: NewInMemoryPredictionRepository(),
	}
}
