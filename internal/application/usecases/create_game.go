package usecases

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"errors"
	"strings"
)

var (
	ErrEmptyGameName  = errors.New("game name cannot be empty")
	ErrInvalidCreator = errors.New("creator id must be positive")
)

// CreateGameUseCase implements use case "creating game"
type CreateGameUseCase struct {
	gameRepo ports.GameRepository
}

func NewCreateGameUseCase(gameRepo ports.GameRepository) *CreateGameUseCase {
	return &CreateGameUseCase{
		gameRepo: gameRepo,
	}
}

// Execute executes creation of game
func (uc *CreateGameUseCase) Execute(command dto.CreateGameCommand) (*dto.GameResponse, error) {
	if err := uc.validateCommand(command); err != nil {
		return nil, err
	}

	game := entities.NewGame(command.Name, command.CreatorID)
	if err := uc.gameRepo.Save(game); err != nil {
		return nil, err
	}

	return uc.toResponse(game), nil
}

func (uc *CreateGameUseCase) validateCommand(command dto.CreateGameCommand) error {
	if strings.TrimSpace(command.Name) == "" {
		return ErrEmptyGameName
	}
	if command.CreatorID <= 0 {
		return ErrInvalidCreator
	}
	return nil
}

func (uc *CreateGameUseCase) toResponse(game *entities.Game) *dto.GameResponse {
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
