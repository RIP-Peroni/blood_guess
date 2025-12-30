package persistence

import (
	"errors"
	"sync"

	"RIP-Peroni/blood_guess/internal/domain/entities"
)

var (
	ErrGameNotFound       = errors.New("game not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrPredictionNotFound = errors.New("prediction not found")
	ErrNilEntity          = errors.New("entity cannot be nil")
)

// === GameRepository Implementation ===

type InMemoryGameRepository struct {
	games map[entities.GameID]*entities.Game
	mu    sync.RWMutex
}

func NewInMemoryGameRepository() *InMemoryGameRepository {
	return &InMemoryGameRepository{
		games: make(map[entities.GameID]*entities.Game),
	}
}

func (r *InMemoryGameRepository) Save(game *entities.Game) error {
	if game == nil {
		return ErrNilEntity
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.games[game.ID()] = game
	return nil
}

func (r *InMemoryGameRepository) FindByID(id entities.GameID) (*entities.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	game, exists := r.games[id]
	if !exists {
		return nil, ErrGameNotFound
	}
	return game, nil
}

func (r *InMemoryGameRepository) FindActiveGames() ([]*entities.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var activeGames []*entities.Game
	for _, game := range r.games {
		if game.Status() != entities.GameStatusFinished {
			activeGames = append(activeGames, game)
		}
	}
	return activeGames, nil
}

func (r *InMemoryGameRepository) FindByStatus(status entities.GameStatus) ([]*entities.Game, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var games []*entities.Game
	for _, game := range r.games {
		if game.Status() == status {
			games = append(games, game)
		}
	}
	return games, nil
}

func (r *InMemoryGameRepository) Update(game *entities.Game) error {
	if game == nil {
		return ErrNilEntity
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.games[game.ID()]; !exists {
		return ErrGameNotFound
	}

	r.games[game.ID()] = game
	return nil
}

// === UserRepository Implementation ===

type InMemoryUserRepository struct {
	users map[entities.UserID]*entities.User
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[entities.UserID]*entities.User),
	}
}

func (r *InMemoryUserRepository) Save(user *entities.User) error {
	if user == nil {
		return ErrNilEntity
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID()] = user
	return nil
}

func (r *InMemoryUserRepository) FindById(id entities.UserID) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) FindByTelegramID(telegramID int64) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.TelegramID() == telegramID {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (r *InMemoryUserRepository) Update(user *entities.User) error {
	if user == nil {
		return ErrNilEntity
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID()]; !exists {
		return ErrUserNotFound
	}

	r.users[user.ID()] = user
	return nil
}

// === PredictionRepository Implementation ===

type InMemoryPredictionRepository struct {
	predictions map[entities.PredictionID]*entities.Prediction
	mu          sync.RWMutex
}

func NewInMemoryPredictionRepository() *InMemoryPredictionRepository {
	return &InMemoryPredictionRepository{
		predictions: make(map[entities.PredictionID]*entities.Prediction),
	}
}

func (r *InMemoryPredictionRepository) Save(prediction *entities.Prediction) error {
	if prediction == nil {
		return ErrNilEntity
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.predictions[prediction.ID()] = prediction
	return nil
}

func (r *InMemoryPredictionRepository) FindByID(id entities.PredictionID) (*entities.Prediction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prediction, exists := r.predictions[id]
	if !exists {
		return nil, ErrPredictionNotFound
	}
	return prediction, nil
}

func (r *InMemoryPredictionRepository) FindByGameAndUser(gameID entities.GameID, userID entities.UserID) ([]*entities.Prediction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Prediction
	for _, prediction := range r.predictions {
		if prediction.GameID() == gameID && prediction.UserID() == userID {
			result = append(result, prediction)
		}
	}
	return result, nil
}

func (r *InMemoryPredictionRepository) FindByGame(gameID entities.GameID) ([]*entities.Prediction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entities.Prediction
	for _, prediction := range r.predictions {
		if prediction.GameID() == gameID {
			result = append(result, prediction)
		}
	}
	return result, nil
}

func (r *InMemoryPredictionRepository) Update(prediction *entities.Prediction) error {
	if prediction == nil {
		return ErrNilEntity
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.predictions[prediction.ID()]; !exists {
		return ErrPredictionNotFound
	}

	r.predictions[prediction.ID()] = prediction
	return nil
}

func (r *InMemoryPredictionRepository) Delete(id entities.PredictionID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.predictions[id]; !exists {
		return ErrPredictionNotFound
	}

	delete(r.predictions, id)
	return nil
}
