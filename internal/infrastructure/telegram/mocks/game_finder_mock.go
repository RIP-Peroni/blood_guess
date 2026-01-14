package mocks

import (
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	"github.com/stretchr/testify/mock"
)

type MockGameFinder struct {
	mock.Mock
}

func (m *MockGameFinder) GameRepo() ports.GameRepository {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(ports.GameRepository)
}

func (m *MockGameFinder) FindActiveGame() (*entities.Game, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Game), args.Error(1)
}

func (m *MockGameFinder) FindLastGameByStatus(status entities.GameStatus) (*entities.Game, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Game), args.Error(1)
}

func (m *MockGameFinder) FindLatestGame() (*entities.Game, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Game), args.Error(1)
}

func (m *MockGameFinder) CanCreateNewGame() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *MockGameFinder) FindPlayerByName(game *entities.Game, playerName string) (*entities.PlayerSlot, error) {
	args := m.Called(game, playerName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PlayerSlot), args.Error(1)
}
