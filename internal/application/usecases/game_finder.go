package usecases

import (
	"errors"
	"sort"
	"time"

	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

var (
	ErrNoActiveGames       = errors.New("нет активных игр")
	ErrNoGamesWithStatus   = errors.New("нет игр с указанным статусом")
	ErrActiveGameExists    = errors.New("уже есть активная игра")
	ErrMultipleActiveGames = errors.New("обнаружено несколько активных игр")
)

type GameFinderInterface interface {
	GameRepo() ports.GameRepository
	FindActiveGame() (*entities.Game, error)
	FindLastGameByStatus(status entities.GameStatus) (*entities.Game, error)
	FindLatestGame() (*entities.Game, error)
	CanCreateNewGame() (bool, error)
	FindPlayerByName(game *entities.Game, playerName string) (*entities.PlayerSlot, error)
}

// GameFinder предоставляет утилиты для поиска игр
type GameFinder struct {
	gameRepo ports.GameRepository
}

// GameRepo возвращает репозиторий игр
func (f *GameFinder) GameRepo() ports.GameRepository {
	return f.gameRepo
}

func NewGameFinder(gameRepo ports.GameRepository) *GameFinder {
	return &GameFinder{
		gameRepo: gameRepo,
	}
}

// FindActiveGame ищет активную игру (не FINISHED)
func (f *GameFinder) FindActiveGame() (*entities.Game, error) {
	activeGames, err := f.gameRepo.FindActiveGames()
	if err != nil {
		return nil, err
	}

	if len(activeGames) == 0 {
		return nil, ErrNoActiveGames
	}

	if len(activeGames) > 1 {
		return nil, ErrMultipleActiveGames
	}

	return activeGames[0], nil
}

// FindLastGameByStatus находит последнюю игру с указанным статусом
func (f *GameFinder) FindLastGameByStatus(status entities.GameStatus) (*entities.Game, error) {
	games, err := f.gameRepo.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	if len(games) == 0 {
		return nil, ErrNoGamesWithStatus
	}

	// Сортируем по времени создания (самая новая первая)
	sort.Slice(games, func(i, j int) bool {
		return games[i].CreatedAt().After(games[j].CreatedAt())
	})

	return games[0], nil
}

// FindLatestGame находит самую новую игру (любого статуса)
func (f *GameFinder) FindLatestGame() (*entities.Game, error) {
	allGames := []entities.GameStatus{
		entities.GameStatusCreated,
		entities.GameStatusPredictionsOpen,
		entities.GameStatusPredictionsClosed,
		entities.GameStatusInProgress,
		entities.GameStatusFinished,
	}

	var latestGame *entities.Game
	var latestTime time.Time

	for _, status := range allGames {
		games, err := f.gameRepo.FindByStatus(status)
		if err != nil {
			continue
		}

		for _, game := range games {
			if game.CreatedAt().After(latestTime) {
				latestGame = game
				latestTime = game.CreatedAt()
			}
		}
	}

	if latestGame == nil {
		return nil, ErrNoActiveGames
	}

	return latestGame, nil
}

// CanCreateNewGame проверяет, можно ли создать новую игру
func (f *GameFinder) CanCreateNewGame() (bool, error) {
	activeGames, err := f.gameRepo.FindActiveGames()
	if err != nil {
		return false, err
	}

	// Ищем игры, которые не завершены
	for _, game := range activeGames {
		if game.Status() != entities.GameStatusFinished {
			return false, ErrActiveGameExists
		}
	}

	return true, nil
}

// FindPlayerByName ищет игрока по имени в указанной игре
func (f *GameFinder) FindPlayerByName(game *entities.Game, playerName string) (*entities.PlayerSlot, error) {
	for _, player := range game.Players() {
		if player.Name == playerName {
			return &player, nil
		}
	}

	return nil, errors.New("игрок не найден")
}
