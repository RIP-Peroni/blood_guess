package dto

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"time"
)

// GameResponse - reply with information about the game
type GameResponse struct {
	ID        string
	Name      string
	Status    entities.GameStatus
	CreatorID int64
	Players   []PlayerResponse
	CreatedAt time.Time
	StartedAt *time.Time
	EndedAt   *time.Time
}

// HasRealRolesSet checks if real roles are set for all players
func (r *GameResponse) HasRealRolesSet() bool {
	for _, player := range r.Players {
		if player.RealRole == "" {
			return false
		}
	}
	return true
}

// PlayerResponse - player information
type PlayerResponse struct {
	ID           string
	Name         string
	AssignedRole string
	RealRole     string
}

// HasRealRole checks if the real role is installed
func (r *PlayerResponse) HasRealRole() bool {
	return r.RealRole != ""
}

// PredictionResponse - answer about the prediction
type PredictionResponse struct {
	ID            string
	GameID        string
	UserID        string
	PlayerSlotID  string
	PredictedRole string
	Points        int  // 0 if it hasn't been calculated yet
	PointsAwarded bool // True if points have already been awarded
	CreatedAt     time.Time
}

// IsCorrect checks whether the prediction is correct
func (r *PredictionResponse) IsCorrect(realRole string) bool {
	return r.PredictedRole == realRole
}

// UserResponse - response with user information
type UserResponse struct {
	ID         string
	TelegramID int64
	Username   string
	Balance    int
	CreatedAt  time.Time
}

// CalculateResultsResponse - answer with scoring results
type CalculateResultsResponse struct {
	GameID string
	Scores map[string]int // UserID -> Points
}

// FromDomainGame converts domain entity Game into DTO
func FromDomainGame(game *entities.Game) *GameResponse {
	players := make([]PlayerResponse, 0, len(game.Players()))
	for _, player := range game.Players() {
		players = append(players, PlayerResponse{
			ID:           string(player.ID),
			Name:         player.Name,
			AssignedRole: player.AssignedRole,
			RealRole:     player.RealRole,
		})
	}

	return &GameResponse{
		ID:        string(game.ID()),
		Name:      game.Name(),
		Status:    game.Status(),
		CreatorID: game.CreatorID(),
		Players:   players,
		CreatedAt: game.CreatedAt(),
	}
}

// FromDomainPrediction converts domain entity Prediction into DTO
func FromDomainPrediction(prediction *entities.Prediction) *PredictionResponse {
	points, awarded := prediction.PointsAwarded()

	return &PredictionResponse{
		ID:            string(prediction.ID()),
		GameID:        string(prediction.GameID()),
		UserID:        string(prediction.UserID()),
		PlayerSlotID:  string(prediction.PlayerSlotID()),
		PredictedRole: prediction.PredictedRole(),
		Points:        points,
		PointsAwarded: awarded,
		CreatedAt:     prediction.CreatedAt(),
	}
}

// FromDomainUser converts domain entity User into DTO
func FromDomainUser(user *entities.User) *UserResponse {
	return &UserResponse{
		ID:         string(user.ID()),
		TelegramID: user.TelegramID(),
		Username:   user.Username(),
		Balance:    user.Balance(),
		CreatedAt:  user.CreatedAt(),
	}
}
