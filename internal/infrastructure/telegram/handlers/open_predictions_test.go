package handlers

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/mocks"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockOpenPredictionsInput реализует ports.OpenPredictionsInput
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

// Проверяем, что мок реализует интерфейс
var _ ports.OpenPredictionsInput = (*MockOpenPredictionsInput)(nil)

func TestOpenPredictionsHandler_Handle_Success(t *testing.T) {
	// Arrange
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockOpenPredictionsInput)
	mockGameFinder := new(mocks.MockGameFinder)

	adminID := int64(67890)

	// Создаем тестовую игру
	game := entities.NewGame("Test Game", adminID)
	err := game.AddPlayer("Player 1")
	require.NoError(t, err)

	// Игра должна быть в статусе CREATED
	require.Equal(t, entities.GameStatusCreated, game.Status())

	// Настраиваем mock GameFinder
	mockGameFinder.On("FindLastGameByStatus", entities.GameStatusCreated).
		Return(game, nil).Once()

	// Настраиваем mock use case
	expectedCommand := dto.OpenPredictionsCommand{
		GameID:  string(game.ID()),
		AdminID: adminID,
	}

	expectedResponse := &dto.GameResponse{
		ID:        string(game.ID()),
		Name:      "Test Game",
		Status:    entities.GameStatusPredictionsOpen,
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

	mockUseCase.On("Execute", expectedCommand).Return(expectedResponse, nil).Once()

	// Настраиваем mock API для отправки сообщения
	mockAPI.On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).
		Run(func(args mock.Arguments) {
			msg := args.Get(0).(tgbotapi.MessageConfig)
			// Проверяем, что это HTML сообщение
			assert.Equal(t, tgbotapi.ModeHTML, msg.ParseMode)
			// Проверяем, что сообщение содержит ожидаемый текст
			assert.Contains(t, msg.Text, "Прогнозы открыты!")
		}).
		Return(tgbotapi.Message{}, nil).Once()

	// Act
	handler := NewOpenPredictionsHandler(mockAPI, mockUseCase, mockGameFinder)

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

	err = handler.Handle(update)

	// Assert
	assert.NoError(t, err)

	// Проверяем, что все моки были вызваны
	mockAPI.AssertExpectations(t)
	mockUseCase.AssertExpectations(t)
	mockGameFinder.AssertExpectations(t)
}

func TestOpenPredictionsHandler_Handle_NoPlayers(t *testing.T) {
	// Arrange
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockOpenPredictionsInput)
	mockGameFinder := new(mocks.MockGameFinder)

	adminID := int64(67890)

	// Создаем тестовую игру БЕЗ игроков
	game := entities.NewGame("Test Game", adminID)
	// Не добавляем игроков

	// Настраиваем mock GameFinder
	mockGameFinder.On("FindLastGameByStatus", entities.GameStatusCreated).
		Return(game, nil).Once()

	// Настраиваем mock API для отправки сообщения об ошибке
	mockAPI.On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).
		Run(func(args mock.Arguments) {
			msg := args.Get(0).(tgbotapi.MessageConfig)
			assert.Equal(t, tgbotapi.ModeHTML, msg.ParseMode)
			assert.Contains(t, msg.Text, "Нельзя открыть прогнозы для игры без игроков")
		}).
		Return(tgbotapi.Message{}, nil).Once()

	// Act
	handler := NewOpenPredictionsHandler(mockAPI, mockUseCase, mockGameFinder)

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

	// Assert
	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
	mockGameFinder.AssertExpectations(t)
	// Use case не должен вызываться
	mockUseCase.AssertNotCalled(t, "Execute")
}

func TestOpenPredictionsHandler_Handle_NotCreator(t *testing.T) {
	// Arrange
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockOpenPredictionsInput)
	mockGameFinder := new(mocks.MockGameFinder)

	creatorID := int64(12345)
	otherUserID := int64(67890)

	// Создаем тестовую игру с другим создателем
	game := entities.NewGame("Test Game", creatorID)
	err := game.AddPlayer("Player 1")
	require.NoError(t, err)

	// Настраиваем mock GameFinder
	mockGameFinder.On("FindLastGameByStatus", entities.GameStatusCreated).
		Return(game, nil).Once()

	// Настраиваем mock API для отправки сообщения об ошибке
	mockAPI.On("Send", mock.AnythingOfType("tgbotapi.MessageConfig")).
		Run(func(args mock.Arguments) {
			msg := args.Get(0).(tgbotapi.MessageConfig)
			assert.Equal(t, tgbotapi.ModeHTML, msg.ParseMode)
			assert.Contains(t, msg.Text, "Только создатель игры может открывать прогнозы")
		}).
		Return(tgbotapi.Message{}, nil).Once()

	// Act
	handler := NewOpenPredictionsHandler(mockAPI, mockUseCase, mockGameFinder)

	update := tgbotapi.Update{
		Message: &tgbotapi.Message{
			Chat: &tgbotapi.Chat{
				ID: 12345,
			},
			From: &tgbotapi.User{
				ID:        otherUserID,
				FirstName: "Other",
				LastName:  "User",
				UserName:  "otheruser",
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

	err = handler.Handle(update)

	// Assert
	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
	mockGameFinder.AssertExpectations(t)
	// Use case не должен вызываться
	mockUseCase.AssertNotCalled(t, "Execute")
}

func TestOpenPredictionsHandler_Command(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)
	mockUseCase := new(MockOpenPredictionsInput)
	mockGameFinder := new(mocks.MockGameFinder)

	handler := NewOpenPredictionsHandler(mockAPI, mockUseCase, mockGameFinder)

	assert.Equal(t, "openpred", handler.Command())
	assert.Equal(t, "Открыть прогнозы для последней созданной игры", handler.Description())
}
