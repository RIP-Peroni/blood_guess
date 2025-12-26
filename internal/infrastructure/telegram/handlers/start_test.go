package handlers

import (
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/mocks"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartHandler_Handle(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)

	expectedMessage := mock.MatchedBy(func(c tgbotapi.Chattable) bool {
		msg, ok := c.(tgbotapi.MessageConfig)
		if !ok {
			return false
		}
		return msg.ChatID == 12345 &&
			msg.ParseMode == tgbotapi.ModeMarkdownV2 &&
			len(msg.Text) > 0
	})

	mockAPI.On("Send", expectedMessage).Return(tgbotapi.Message{}, nil)

	handler := NewStartHandler(mockAPI)

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
		},
	}

	err := handler.Handle(update)

	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
	assert.Equal(t, "start", handler.Command())
	assert.Equal(t, "Начать работу с ботом", handler.Description())
}
