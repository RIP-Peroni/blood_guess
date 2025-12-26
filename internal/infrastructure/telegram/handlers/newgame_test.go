// Файл: ./internal/infrastructure/telegram/handlers/newgame_test.go
package handlers

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/mocks"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock для CreateGameInput
type MockCreateGameInput struct {
	mock.Mock
}

func (m *MockCreateGameInput) Execute(command dto.CreateGameCommand) (*dto.GameResponse, error) {
	args := m.Called(command)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.GameResponse), args.Error(1)
}

func TestNewGameHandler_Handle(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockCreateGameInput)

	// Set up expectations for a successful case
	// Note: CreatorID must be 67890 (as in update.Message.From.ID)
	expectedCommand := dto.CreateGameCommand{
		Name:      "Тестовая игра",
		CreatorID: 67890,
	}

	expectedResponse := &dto.GameResponse{
		ID:        "game-123",
		Name:      "Тестовая игра",
		Status:    "created",
		CreatorID: 67890,
		Players:   []dto.PlayerResponse{},
		CreatedAt: time.Now(),
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

	handler := NewNewGameHandler(mockAPI, mockUseCase)

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
			Text: "/newgame Тестовая игра",
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
	assert.Equal(t, "newgame", handler.Command())
	assert.Equal(t, "Создать новую игру", handler.Description())
}
