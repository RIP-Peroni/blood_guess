package usecases

import (
	"fmt"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

// ClosePredictionsUseCase implements use case "closing predictions"
type ClosePredictionsUseCase struct {
	gameRepo ports.GameRepository
}

func NewClosePredictionsUseCase(gameRepo ports.GameRepository) *ClosePredictionsUseCase {
	return &ClosePredictionsUseCase{
		gameRepo: gameRepo,
	}
}

// Execute executes closing of predictions
func (uc *ClosePredictionsUseCase) Execute(command dto.ClosePredictionsCommand) (*dto.GameResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	if game.CreatorID() != command.AdminID {
		return nil, ErrNotGameCreator
	}

	if game.Status() != entities.GameStatusPredictionsOpen {
		return nil, fmt.Errorf("%w: current status is %s", ErrPredictionsNotOpen, game.Status())
	}

	if err := game.ClosePredictions(); err != nil {
		return nil, fmt.Errorf("failed to close predictions: %w", err)
	}

	if err := uc.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(game), nil
}

// toResponse converts the domain entity into a response DTO
func (uc *ClosePredictionsUseCase) toResponse(game *entities.Game) *dto.GameResponse {
	players := make([]dto.PlayerResponse, 0, len(game.Players()))
	for _, player := range game.Players() {
		players = append(players, dto.PlayerResponse{
			ID:           string(player.ID),
			Name:         player.Name,
			AssignedRole: player.AssignedRole,
			RealRole:     player.RealRole,
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
