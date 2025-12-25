package ports

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
)

// CreateGameInput - port for game creation
type CreateGameInput interface {
	Execute(command dto.CreateGameCommand) (*dto.GameResponse, error)
}

// SubmitPredictionInput - port for prediction sending
type SubmitPredictionInput interface {
	Execute(command dto.SubmitPredictionCommand) (*dto.PredictionResponse, error)
}

// OpenPredictionsInput - port for predictions opening
type OpenPredictionsInput interface {
	Execute(command dto.OpenPredictionsCommand) (*dto.GameResponse, error)
}

// CalculateResultsInput - port for scoring results
type CalculateResultsInput interface {
	Execute(command dto.FinishGameCommand) (map[string]int, error) //[userID]points
}
