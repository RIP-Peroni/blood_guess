package usecases

import "errors"

// Common mistakes for all use cases
var (
	ErrGameNotFound                = errors.New("game not found")
	ErrUserNotFound                = errors.New("user not found")
	ErrNotGameCreator              = errors.New("only game creator can perform this action")
	ErrGameNotAcceptingPredictions = errors.New("game is not accepting predictions")
	ErrAlreadyPredicted            = errors.New("user already predicted for this player slot")
	ErrPlayerSlotNotFound          = errors.New("player slot not found")
	ErrInvalidGameState            = errors.New("invalid game state for this operation")
	ErrNoPlayersInGame             = errors.New("game has no players")
	ErrPredictionsNotOpen          = errors.New("predictions are not open")
)
