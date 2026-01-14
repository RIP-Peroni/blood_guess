package usecases

import (
	"fmt"
	"time"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

// FinishGameUseCase implements the "finish game" use case
type FinishGameUseCase struct {
	gameRepo ports.GameRepository
}

func NewFinishGameUseCase(
	gameRepo ports.GameRepository,
) *FinishGameUseCase {
	return &FinishGameUseCase{
		gameRepo: gameRepo,
	}
}

// Execute finishes the game (simply changes status to FINISHED)
func (uc *FinishGameUseCase) Execute(command dto.FinishGameCommand) (*dto.FinishGameResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	// Get the game
	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	// Check if user is the game creator
	if game.CreatorID() != command.AdminID {
		return nil, ErrNotGameCreator
	}

	// Check if game is in correct state (in progress)
	if game.Status() != entities.GameStatusInProgress {
		return nil, fmt.Errorf("%w: current status is %s", ErrInvalidGameState, game.Status())
	}

	// Finish the game (simply change status)
	if err := game.Finish(); err != nil {
		return nil, fmt.Errorf("failed to finish game: %w", err)
	}

	// Save the updated game
	if err := uc.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(game), nil
}

// toResponse converts domain entity to response DTO
func (uc *FinishGameUseCase) toResponse(game *entities.Game) *dto.FinishGameResponse {
	// Get endedAt time
	endedAt := time.Now()
	if game.EndedAt() != nil {
		endedAt = *game.EndedAt()
	}

	return &dto.FinishGameResponse{
		GameID:          string(game.ID()),
		Name:            game.Name(),
		Status:          game.Status(),
		PlayerResults:   []dto.PlayerResultResponse{},
		UserScores:      make(map[string]int),
		CurrencyAwarded: false,
		EndedAt:         endedAt,
	}
}
