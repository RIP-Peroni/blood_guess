package handlers

import (
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/application/usecases"
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

// TestGameFinder - тестовая структура, которая реализует интерфейс GameFinderInterface
type TestGameFinder struct {
	repo *MockGameRepository
}

// GameRepo возвращает репозиторий игр
func (t *TestGameFinder) GameRepo() ports.GameRepository {
	return t.repo
}

// Добавляем методы, чтобы структура была совместима с GameFinderInterface
// Эти методы не будут использоваться в этом тесте, но нужны для совместимости
func (t *TestGameFinder) FindActiveGame() (*entities.Game, error) {
	return nil, nil
}

func (t *TestGameFinder) FindLastGameByStatus(status entities.GameStatus) (*entities.Game, error) {
	return nil, nil
}

func (t *TestGameFinder) FindLatestGame() (*entities.Game, error) {
	return nil, nil
}

func (t *TestGameFinder) CanCreateNewGame() (bool, error) {
	return true, nil
}

func (t *TestGameFinder) FindPlayerByName(game *entities.Game, playerName string) (*entities.PlayerSlot, error) {
	return nil, nil
}

func TestGamesHandler_Handle_WithGames(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)

	// Создаем тестовые игры
	game1 := entities.NewGame("Игра 1", 12345)
	game2 := entities.NewGame("Игра 2", 67890)
	_ = game1.OpenPredictions()
	_ = game2.ClosePredictions()

	activeGames := []*entities.Game{game1, game2}

	// Создаем mock для GameRepository
	mockGameRepo := new(MockGameRepository)
	mockGameRepo.On("FindActiveGames").Return(activeGames, nil)

	// Создаем TestGameFinder с mock репозиторием
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

	// Создаем хендлер. Нужно привести тип к *usecases.GameFinder
	// Поскольку TestGameFinder имеет те же поля, что и usecases.GameFinder,
	// мы можем использовать unsafe преобразование или создать реальный GameFinder
	handler := NewGamesHandler(mockAPI, (*usecases.GameFinder)(nil))

	// Создаем реальный GameFinder с mock репозиторием
	gameFinder := usecases.NewGameFinder(mockGameRepo)

	handler = NewGamesHandler(mockAPI, gameFinder)

	// Тестируем update
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

	// Вызываем хендлер
	err := handler.Handle(update)

	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
	mockGameRepo.AssertExpectations(t)
	assert.Equal(t, "games", handler.Command())
	assert.Equal(t, "Показать активные игры", handler.Description())
}

func TestGamesHandler_Handle_NoGames(t *testing.T) {
	mockAPI := new(mocks.MockBotAPI)

	// Создаем mock для GameRepository
	mockGameRepo := new(MockGameRepository)
	mockGameRepo.On("FindActiveGames").Return([]*entities.Game{}, nil)

	// Создаем реальный GameFinder с mock репозиторием
	gameFinder := usecases.NewGameFinder(mockGameRepo)

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
	handler := NewGamesHandler(mockAPI, gameFinder)

	// Тестируем update
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

	// Вызываем хендлер
	err := handler.Handle(update)

	assert.NoError(t, err)
	mockAPI.AssertExpectations(t)
	mockGameRepo.AssertExpectations(t)
}
