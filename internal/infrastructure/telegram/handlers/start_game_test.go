package handlers

import (
	"testing"
	"time"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
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
		gameRepo := persistence.NewInMemoryGameRepository()

		// Create test game
		adminID := int64(12345)
		game := entities.NewGame("Test Game", adminID)
		err := game.AddPlayer("Player 1")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = game.ClosePredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		gameID := string(game.ID())

		// Setup mock expectations
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

		// Expect message to be sent
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

		// Create handler
		handler := NewStartGameHandler(mockAPI, mockUseCase, gameRepo)

		// Test update
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
				Text: "/startgame " + gameID,
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 10,
					},
				},
			},
		}

		// Execute
		err = handler.Handle(update)
		assert.NoError(t, err)

		// Verify mocks
		mockAPI.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("error: not creator", func(t *testing.T) {
		mockAPI := new(mocks.MockBotAPI)
		mockUseCase := new(MockStartGameInput)
		gameRepo := persistence.NewInMemoryGameRepository()

		// Create game with different creator
		creatorID := int64(12345)
		otherUserID := int64(99999)
		game := entities.NewGame("Test Game", creatorID)
		err := gameRepo.Save(game)
		require.NoError(t, err)

		gameID := string(game.ID())

		// Expect error message
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

		handler := NewStartGameHandler(mockAPI, mockUseCase, gameRepo)

		update := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 67890,
				},
				From: &tgbotapi.User{
					ID: otherUserID,
				},
				Text: "/startgame " + gameID,
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
	})

	t.Run("command without arguments shows usage", func(t *testing.T) {
		mockAPI := new(mocks.MockBotAPI)
		mockUseCase := new(MockStartGameInput)
		gameRepo := persistence.NewInMemoryGameRepository()

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

		handler := NewStartGameHandler(mockAPI, mockUseCase, gameRepo)

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
	})
}

func TestStartGameHandler_Command(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockStartGameInput)
	gameRepo := persistence.NewInMemoryGameRepository()

	handler := NewStartGameHandler(mockAPI, mockUseCase, gameRepo)

	assert.Equal(t, "startgame", handler.Command())
	assert.Equal(t, "Начать реальную игру (после закрытия прогнозов)", handler.Description())
}
