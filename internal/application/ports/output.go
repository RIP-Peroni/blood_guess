package ports

import "RIP-Peroni/blood_guess/internal/domain/entities"

type UserRepository interface {
	Save(user *entities.User) error
	FindById(id entities.UserID) (*entities.User, error)
	FindByTelegramID(telegramID int64) (*entities.User, error)
	Update(user *entities.User) error
}

type GameRepository interface {
	Save(game *entities.Game) error
	FindByID(id entities.GameID) (*entities.Game, error)
	FindActiveGames() ([]*entities.Game, error)
	FindByStatus(status entities.GameStatus) ([]*entities.Game, error)
	Update(game *entities.Game) error
}

type PredictionRepository interface {
	Save(prediction *entities.Prediction) error
	FindByID(id entities.PredictionID) (*entities.Prediction, error)
	FindByGameAndUser(gameID entities.GameID, userID entities.UserID) ([]*entities.Prediction, error)
	FindByGame(gameID entities.GameID) ([]*entities.Prediction, error)
	Update(prediction *entities.Prediction) error
	Delete(id entities.PredictionID) error
}
