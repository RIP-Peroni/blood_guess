package usecases

import (
	"errors"
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

// AddPlayersUseCase implements the addition of multiple players
type AddPlayersUseCase struct {
	gameRepo ports.GameRepository
}

func NewAddPlayersUseCase(gameRepo ports.GameRepository) *AddPlayersUseCase {
	return &AddPlayersUseCase{
		gameRepo: gameRepo,
	}
}

// Execute performs the addition of multiple players
func (uc *AddPlayersUseCase) Execute(command dto.AddPlayersCommand) (*dto.GameResponse, error) {
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

	trimmedNames := make([]string, len(command.PlayerNames))
	for i, name := range command.PlayerNames {
		trimmedNames[i] = strings.TrimSpace(name)
	}

	for _, name := range trimmedNames {
		if err := game.AddPlayer(name); err != nil {
			if errors.Is(err, entities.ErrPlayerAlreadyAdded) {
				return nil, fmt.Errorf("player '%s' already exists in this game", name)
			}
			return nil, fmt.Errorf("failed to add player %s: %w", name, err)
		}
	}

	if err := uc.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(game), nil
}

// toResponse converts an entity into a DTO
func (uc *AddPlayersUseCase) toResponse(game *entities.Game) *dto.GameResponse {
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
