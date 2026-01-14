package handlers

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/mocks"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для OpenPredictionsInput
type MockOpenPredictionsInput struct {
	mock.Mock
}

func (m *MockOpenPredictionsInput) Execute(command dto.OpenPredictionsCommand) (*dto.GameResponse, error) {
	args := m.Called(command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GameResponse), args.Error(1)
}

// Mock for GameRepository that supports the required methods
type MockGameRepositoryForOpenPredictions struct {
	mock.Mock
}

func (m *MockGameRepositoryForOpenPredictions) Save(game *entities.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

func (m *MockGameRepositoryForOpenPredictions) FindByID(id entities.GameID) (*entities.Game, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Game), args.Error(1)
}

func (m *MockGameRepositoryForOpenPredictions) FindActiveGames() ([]*entities.Game, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Game), args.Error(1)
}

func (m *MockGameRepositoryForOpenPredictions) FindByStatus(status entities.GameStatus) ([]*entities.Game, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entities.Game), args.Error(1)
}

func (m *MockGameRepositoryForOpenPredictions) Update(game *entities.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

func TestOpenPredictionsHandler_Handle(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockOpenPredictionsInput)
	mockGameFinder := new(mocks.MockGameFinder)

	// Тест 1: Успешное открытие прогнозов
	t.Run("successful opening of predictions", func(t *testing.T) {
		adminID := int64(67890)

		// Создаем тестовую игру
		game := entities.NewGame("Test Game", adminID)
		gameID := string(game.ID())

		// Настраиваем mock GameFinder
		mockGameFinder.On("FindLastGameByStatus", entities.GameStatusCreated).
			Return(game, nil)

		// Настраиваем mock use case
		expectedCommand := dto.OpenPredictionsCommand{
			GameID:  gameID,
			AdminID: adminID,
		}

		expectedResponse := &dto.GameResponse{
			ID:        gameID,
			Name:      "Test Game",
			Status:    entities.GameStatusPredictionsOpen,
			CreatorID: adminID,
			Players:   []dto.PlayerResponse{},
			CreatedAt: time.Now(),
		}

		mockUseCase.On("Execute", expectedCommand).Return(expectedResponse, nil)

		// Ожидаем отправку сообщения
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

		// Создаем хендлер
		handler := NewOpenPredictionsHandler(mockAPI, mockUseCase, mockGameFinder)

		// Тестируем update
		update := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 12345,
				},
				From: &tgbotapi.User{
					ID:        adminID,
					FirstName: "Test",
					LastName:  "User",
					UserName:  "testuser",
				},
				Text: "/openpred",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 8,
					},
				},
			},
		}

		err := handler.Handle(update)

		assert.NoError(t, err)
		mockAPI.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
		mockGameFinder.AssertExpectations(t)
	})

	// Тест 2: Команда без аргументов показывает использование
	t.Run("command without arguments shows usage", func(t *testing.T) {
		mockAPI := new(mocks.MockBotAPI)
		mockUseCase := new(MockOpenPredictionsInput)
		mockGameFinder := new(mocks.MockGameFinder)

		// Создаем тестовую игру
		game := entities.NewGame("Test Game", 67890)
		_ = game.AddPlayer("Player 1")

		// Настраиваем mock GameFinder
		mockGameFinder.On("FindLastGameByStatus", entities.GameStatusCreated).
			Return(game, nil)

		// Ожидаем отправку сообщения
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

		handler := NewOpenPredictionsHandler(mockAPI, mockUseCase, mockGameFinder)

		update := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 12345,
				},
				From: &tgbotapi.User{
					ID:        67890,
					FirstName: "Test",
					LastName:  "User",
					UserName:  "testuser",
				},
				Text: "/openpred",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 8,
					},
				},
			},
		}

		err := handler.Handle(update)
		assert.NoError(t, err)
		mockAPI.AssertExpectations(t)
	})
}

func TestOpenPredictionsHandler_Command(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockOpenPredictionsInput)
	mockGameFinder := new(mocks.MockGameFinder)

	handler := NewOpenPredictionsHandler(mockAPI, mockUseCase, mockGameFinder)

	assert.Equal(t, "openpred", handler.Command())
	assert.Equal(t, "Открыть прогнозы для последней созданной игры", handler.Description())
}
