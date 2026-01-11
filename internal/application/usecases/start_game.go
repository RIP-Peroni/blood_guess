package usecases

import (
	"fmt"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

// StartGameUseCase implements the "start game" use case
type StartGameUseCase struct {
	gameRepo ports.GameRepository
}

func NewStartGameUseCase(gameRepo ports.GameRepository) *StartGameUseCase {
	return &StartGameUseCase{
		gameRepo: gameRepo,
	}
}

// Execute starts the real game
func (uc *StartGameUseCase) Execute(command dto.StartGameCommand) (*dto.GameResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	// Check if user is the game creator
	if game.CreatorID() != command.AdminID {
		return nil, ErrNotGameCreator
	}

	// Check if game is in correct state (predictions closed)
	if game.Status() != entities.GameStatusPredictionsClosed {
		return nil, fmt.Errorf("%w: current status is %s", ErrInvalidGameState, game.Status())
	}

	// Start the game
	if err := game.Start(); err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}

	// Save the updated game
	if err := uc.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(game), nil
}

// toResponse converts domain entity to response DTO
func (uc *StartGameUseCase) toResponse(game *entities.Game) *dto.GameResponse {
	players := make([]dto.PlayerResponse, 0, len(game.Players()))
	for _, player := range game.Players() {
		players = append(players, dto.PlayerResponse{
			ID:       string(player.ID),
			Name:     player.Name,
			RealRole: player.RealRole,
		})
	}

	return &dto.GameResponse{
		ID:        string(game.ID()),
		Name:      game.Name(),
		Status:    game.Status(),
		CreatorID: game.CreatorID(),
		Players:   players,
		CreatedAt: game.CreatedAt(),
	}
}
