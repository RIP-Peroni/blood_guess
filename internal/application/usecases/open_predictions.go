package usecases

import (
	"fmt"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

// OpenPredictionsUseCase implements use case "opening predictions"
type OpenPredictionsUseCase struct {
	gameRepo ports.GameRepository
}

func NewOpenPredictionsUseCase(gameRepo ports.GameRepository) *OpenPredictionsUseCase {
	return &OpenPredictionsUseCase{
		gameRepo: gameRepo,
	}
}

// Execute executes opening of predictions
func (uc *OpenPredictionsUseCase) Execute(command dto.OpenPredictionsCommand) (*dto.GameResponse, error) {
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

	if game.Status() != entities.GameStatusCreated {
		return nil, fmt.Errorf("%w: current status is %s", ErrInvalidGameState, game.Status())
	}

	if len(game.Players()) == 0 {
		return nil, ErrNoPlayersInGame
	}

	if err := game.OpenPredictions(); err != nil {
		return nil, fmt.Errorf("failed to open predictions: %w", err)
	}

	if err := uc.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(game), nil
}

// toResponse converts the domain entity into a response DTO
func (uc *OpenPredictionsUseCase) toResponse(game *entities.Game) *dto.GameResponse {
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
