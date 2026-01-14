package ports

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
)

// CreateGameInput - port for creating a game
type CreateGameInput interface {
	Execute(command dto.CreateGameCommand) (*dto.GameResponse, error)
}

// AddPlayerInput - port for adding a player
type AddPlayerInput interface {
	Execute(command dto.AddPlayerCommand) (*dto.GameResponse, error)
}

// OpenPredictionsInput - port for opening predictions
type OpenPredictionsInput interface {
	Execute(command dto.OpenPredictionsCommand) (*dto.GameResponse, error)
}

// ClosePredictionsInput - port for closing predictions
type ClosePredictionsInput interface {
	Execute(command dto.ClosePredictionsCommand) (*dto.GameResponse, error)
}

// SubmitPredictionInput - port for sending the prediction
type SubmitPredictionInput interface {
	Execute(command dto.SubmitPredictionCommand) (*dto.PredictionResponse, error)
}

// CalculateResultsInput - port for calculating results
type CalculateResultsInput interface {
	Execute(command dto.FinishGameCommand) (map[string]int, error) //[userID]points
}

// StartGameInput port for starting a game
type StartGameInput interface {
	Execute(command dto.StartGameCommand) (*dto.GameResponse, error)
}

// FinishGameInput port for finishing a game
type FinishGameInput interface {
	Execute(command dto.FinishGameCommand) (*dto.FinishGameResponse, error)
}

// AddPlayersInput - port for adding multiple players
type AddPlayersInput interface {
	Execute(command dto.AddPlayersCommand) (*dto.GameResponse, error)
}

// CopyPlayersInput - port for copying players from a previous game
type CopyPlayersInput interface {
	Execute(command dto.CopyPlayersCommand) (*dto.GameResponse, error)
}

type SetRealRoleInput interface {
	Execute(command dto.SetRealRoleCommand) error
}

// AwardPointsInput - port for awarding points after roles are set
type AwardPointsInput interface {
	Execute(command dto.AwardPointsCommand) (*dto.FinishGameResponse, error)
}
