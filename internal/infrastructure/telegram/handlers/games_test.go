package handlers

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/mocks"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для GameRepository
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

func TestGamesHandler_Handle_WithGames(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockRepo := new(MockGameRepository)

	// Create test games
	game1 := entities.NewGame("Игра 1", 12345)
	game2 := entities.NewGame("Игра 2", 67890)
	_ = game1.OpenPredictions()
	_ = game2.ClosePredictions()

	activeGames := []*entities.Game{game1, game2}

	// Set up a mock repository
	mockRepo.On("FindActiveGames").Return(activeGames, nil)

	// Wait for the message to be sent
	expectedMessage := mock.MatchedBy(func(c tgbotapi.Chattable) bool {
		msg, ok := c.(tgbotapi.MessageConfig)
		if !ok {
			return false
		}
		return msg.ChatID == 12345 &&
			msg.ParseMode == tgbotapi.ModeHTML &&
			len(msg.Text) > 0
	})

	mockAPI.On("Send", expectedMessage).Return(tgbotapi.Message{}, nil)

	// Create a handler
	handler := NewGamesHandler(mockAPI, mockRepo)

	// Test update
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 12345,
			},
			From: &tgbotapi.User{
				ID: 67890,
			},
		},
	}

	// Call the handler
	err := handler.Handle(update)

	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	assert.Equal(t, "games", handler.Command())
	assert.Equal(t, "Показать активные игры", handler.Description())
}

func TestGamesHandler_Handle_NoGames(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockRepo := new(MockGameRepository)

	// Setting up a mock: no active games
	mockRepo.On("FindActiveGames").Return([]*entities.Game{}, nil)

	// Wait for the message to be sent
	expectedMessage := mock.MatchedBy(func(c tgbotapi.Chattable) bool {
		msg, ok := c.(tgbotapi.MessageConfig)
		if !ok {
			return false
		}
		return msg.ChatID == 12345 &&
			msg.ParseMode == tgbotapi.ModeHTML &&
			len(msg.Text) > 0
	})

	mockAPI.On("Send", expectedMessage).Return(tgbotapi.Message{}, nil)

	// Create a handler
	handler := NewGamesHandler(mockAPI, mockRepo)

	// Test update
	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 12345,
			},
			From: &tgbotapi.User{
				ID: 67890,
			},
		},
	}

	// Call handler
	err := handler.Handle(update)

	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
