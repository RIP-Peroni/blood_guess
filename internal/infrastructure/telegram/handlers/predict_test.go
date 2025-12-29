package handlers

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/persistence"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/mocks"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock для SubmitPredictionInput
type MockSubmitPredictionInput struct {
	mock.Mock
}

func (m *MockSubmitPredictionInput) Execute(command dto.SubmitPredictionCommand) (*dto.PredictionResponse, error) {
	args := m.Called(command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PredictionResponse), args.Error(1)
}

func TestPredictHandler_Handle(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockSubmitPredictionInput)
	gameRepo := persistence.NewInMemoryGameRepository()
	userRepo := persistence.NewInMemoryUserRepository()

	t.Run("successful prediction submission", func(t *testing.T) {
		game := entities.NewGame("Test Game", 12345)
		err := game.AddPlayer("Alice", "townsfolk")
		require.NoError(t, err)

		err = game.OpenPredictions()
		require.NoError(t, err)

		err = gameRepo.Save(game)
		require.NoError(t, err)

		playerID := string(game.Players()[0].ID)

		user := entities.NewUser(67890, "testuser")
		err = userRepo.Save(user)
		require.NoError(t, err)

		expectedCommand := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  playerID,
			PredictedRole: "demon",
		}

		expectedResponse := &dto.PredictionResponse{
			ID:            "pred-123",
			GameID:        string(game.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  playerID,
			PredictedRole: "demon",
			Points:        0,
			PointsAwarded: false,
			CreatedAt:     time.Now(),
		}

		mockUseCase.On("Execute", expectedCommand).Return(expectedResponse, nil)

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

		handler := NewPredictHandler(mockAPI, mockUseCase, gameRepo, userRepo)

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
				Text: "/predict " + string(game.ID()) + " " + playerID + " demon",
				Entities: []tgbotapi.MessageEntity{
					{
						Type:   "bot_command",
						Offset: 0,
						Length: 8,
					},
				},
			},
		}

		err = handler.Handle(update)
		assert.NoError(t, err)
		mockAPI.AssertExpectations(t)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("command without arguments shows usage", func(t *testing.T) {
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

		handler := NewPredictHandler(mockAPI, mockUseCase, gameRepo, userRepo)

		update := tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 12345,
				},
				From: &tgbotapi.User{
					ID: 67890,
				},
				Text: "/predict",
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

func TestPredictHandler_Command(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockSubmitPredictionInput)
	gameRepo := persistence.NewInMemoryGameRepository()
	userRepo := persistence.NewInMemoryUserRepository()

	handler := NewPredictHandler(mockAPI, mockUseCase, gameRepo, userRepo)

	assert.Equal(t, "predict", handler.Command())
	assert.Equal(t, "Сделать прогноз на роль игрока в игре", handler.Description())
}
