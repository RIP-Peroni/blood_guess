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
}

// PlayerResponse - player information
type PlayerResponse struct {
	ID           string
	Name         string
	AssignedRole string
	RealRole     string
}

// PredictionResponse - answer about the prediction
type PredictionResponse struct {
	ID            string
	GameID        string
	UserID        string
	PlayerSlotID  string
	PredictedRole string
	Points        int  // 0 if not yet calculated
	PointsAwarded bool // True if points have already been awarded
	CreatedAt     time.Time
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

type FinishGameResponse struct {
	GameID          string
	Name            string
	Status          entities.GameStatus
	PlayerResults   []PlayerResultResponse
	UserScores      map[string]int // userID -> totalScore
	CurrencyAwarded bool
	EndedAt         time.Time
}

type PlayerResultResponse struct {
	PlayerID     string
	PlayerName   string
	AssignedRole string
	RealRole     string
	Predictions  []PredictionResultResponse
}

type PredictionResultResponse struct {
	UserID        string
	Username      string
	PredictedRole string
	Points        int
	IsCorrect     bool
}
