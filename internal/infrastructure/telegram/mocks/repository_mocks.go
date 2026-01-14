package mocks

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) Save(game *entities.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

func (m *MockGameRepository) FindByID(id entities.GameID) (*entities.Game, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Game), args.Error(1)
}

func (m *MockGameRepository) FindActiveGames() ([]*entities.Game, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Game), args.Error(1)
}

func (m *MockGameRepository) FindByStatus(status entities.GameStatus) ([]*entities.Game, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Game), args.Error(1)
}

func (m *MockGameRepository) Update(game *entities.Game) error {
	args := m.Called(game)
	return args.Error(0)
}
