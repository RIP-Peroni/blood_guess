package handlers

import (
	"testing"
	"time"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/usecases"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/mocks"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock для StartGameInput
type MockStartGameInput struct {
	mock.Mock
}

func (m *MockStartGameInput) Execute(command dto.StartGameCommand) (*dto.GameResponse, error) {
	args := m.Called(command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GameResponse), args.Error(1)
}

func TestStartGameHandler_Handle(t *testing.T) {
	t.Run("successful game start", func(t *testing.T) {
		mockAPI := new(mocks.MockBotAPI)
		mockUseCase := new(MockStartGameInput)
		mockGameFinder := new(mocks.MockGameFinder)

		// Создаем тестовую игру
		adminID := int64(12345)
		game := entities.NewGame("Test Game", adminID)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)
		err = game.OpenPredictions()
		require.NoError(t, err)
		err = game.ClosePredictions()
		require.NoError(t, err)

		gameID := string(game.ID())

		// Настраиваем mock GameFinder
		mockGameFinder.On("FindLastGameByStatus", entities.GameStatusPredictionsClosed).
			Return(game, nil)

		// Настраиваем mock use case
		expectedCommand := dto.StartGameCommand{
			GameID:  gameID,
			AdminID: adminID,
		}

		expectedResponse := &dto.GameResponse{
			ID:        gameID,
			Name:      "Test Game",
			Status:    entities.GameStatusInProgress,
			CreatorID: adminID,
			Players: []dto.PlayerResponse{
				{
					ID:       string(game.Players()[0].ID),
					Name:     "Player 1",
					RealRole: "",
				},
			},
			CreatedAt: time.Now(),
		}

		mockUseCase.On("Execute", expectedCommand).Return(expectedResponse, nil)

		// Ожидаем отправку сообщения
		expectedMessage := mock.MatchedBy(func(c tgbotapi.Chattable) bool {
			msg, ok := c.(tgbotapi.MessageConfig)
			if !ok {
				return false
			}
			return msg.ChatID == 67890 &&
				msg.ParseMode == tgbotapi.ModeHTML &&
				len(msg.Text) > 0
		})

		mockAPI.On("Send", expectedMessage).Return(tgbotapi.Message{}, nil)

		// Создаем хендлер
		handler := NewStartGameHandler(mockAPI, mockUseCase, mockGameFinder)

		// Тестируем update
		update := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 67890,
				},
				From: &tgbotapi.User{
					ID:        adminID,
					FirstName: "Admin",
					UserName:  "admin_user",
				},
				Text: "/startgame",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 10,
					},
				},
			},
		}

		// Выполняем
		err = handler.Handle(update)
		assert.NoError(t, err)

		// Проверяем моки
		mockAPI.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
		mockGameFinder.AssertExpectations(t)
	})

	t.Run("error: not creator", func(t *testing.T) {
		mockAPI := new(mocks.MockBotAPI)
		mockUseCase := new(MockStartGameInput)
		mockGameFinder := new(mocks.MockGameFinder)

		// Создаем игру с другим создателем
		creatorID := int64(12345)
		otherUserID := int64(99999)
		game := entities.NewGame("Test Game", creatorID)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)
		err = game.OpenPredictions()
		require.NoError(t, err)
		err = game.ClosePredictions()
		require.NoError(t, err)

		// Настраиваем mock GameFinder
		mockGameFinder.On("FindLastGameByStatus", entities.GameStatusPredictionsClosed).
			Return(game, nil)

		// Ожидаем сообщение об ошибке
		expectedMessage := mock.MatchedBy(func(c tgbotapi.Chattable) bool {
			msg, ok := c.(tgbotapi.MessageConfig)
			if !ok {
				return false
			}
			return msg.ChatID == 67890 &&
				msg.ParseMode == tgbotapi.ModeHTML &&
				len(msg.Text) > 0
		})

		mockAPI.On("Send", expectedMessage).Return(tgbotapi.Message{}, nil)

		handler := NewStartGameHandler(mockAPI, mockUseCase, mockGameFinder)

		update := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 67890,
				},
				From: &tgbotapi.User{
					ID: otherUserID,
				},
				Text: "/startgame",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 10,
					},
				},
			},
		}

		err = handler.Handle(update)
		assert.NoError(t, err)
		mockAPI.AssertExpectations(t)
		mockGameFinder.AssertExpectations(t)
		// Use case не должен вызываться
		mockUseCase.AssertNotCalled(t, "Execute", mock.Anything)
	})

	t.Run("command without arguments shows usage", func(t *testing.T) {
		mockAPI := new(mocks.MockBotAPI)
		mockUseCase := new(MockStartGameInput)
		mockGameFinder := new(mocks.MockGameFinder)

		// Настраиваем mock GameFinder - нет игры
		mockGameFinder.On("FindLastGameByStatus", entities.GameStatusPredictionsClosed).
			Return(nil, usecases.ErrNoGamesWithStatus)

		expectedMessage := mock.MatchedBy(func(c tgbotapi.Chattable) bool {
			msg, ok := c.(tgbotapi.MessageConfig)
			if !ok {
				return false
			}
			return msg.ChatID == 67890 &&
				msg.ParseMode == tgbotapi.ModeHTML &&
				len(msg.Text) > 0
		})

		mockAPI.On("Send", expectedMessage).Return(tgbotapi.Message{}, nil)

		handler := NewStartGameHandler(mockAPI, mockUseCase, mockGameFinder)

		update := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 67890,
				},
				From: &tgbotapi.User{
					ID: 12345,
				},
				Text: "/startgame",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 10,
					},
				},
			},
		}

		err := handler.Handle(update)
		assert.NoError(t, err)
		mockAPI.AssertExpectations(t)
		mockGameFinder.AssertExpectations(t)
	})
}

func TestStartGameHandler_Command(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockStartGameInput)
	mockGameFinder := new(mocks.MockGameFinder)

	handler := NewStartGameHandler(mockAPI, mockUseCase, mockGameFinder)

	assert.Equal(t, "startgame", handler.Command())
	assert.Equal(t, "Начать реальную игру (после закрытия прогнозов)", handler.Description())
}
