package usecases

import (
	"errors"
	"fmt"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

var (
	ErrGameNotInCreatedState = errors.New("game is not in 'created' state")
	ErrDuplicatePlayerName   = errors.New("player with this name already exists in the game")
)

// AddPlayerUseCase implements use case "adding player to game"
type AddPlayerUseCase struct {
	gameRepo ports.GameRepository
}

func NewAddPlayerUseCase(gameRepo ports.GameRepository) *AddPlayerUseCase {
	return &AddPlayerUseCase{
		gameRepo: gameRepo,
	}
}

// Execute executes adding a player to the game
func (uc *AddPlayerUseCase) Execute(command dto.AddPlayerCommand) (*dto.GameResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	// Проверяем, что пользователь - создатель игры
	if game.CreatorID() != command.AdminID {
		return nil, ErrNotGameCreator
	}

	// Проверяем, что игра находится в состоянии created (можно добавлять игроков)
	if game.Status() != entities.GameStatusCreated {
		return nil, fmt.Errorf("%w: current status is %s", ErrGameNotInCreatedState, game.Status())
	}

	// Добавляем игрока
	if err := game.AddPlayer(command.PlayerName, command.AssignedRole); err != nil {
		if err.Error() == "player with this name already exists in this game" {
			return nil, ErrDuplicatePlayerName
		}
		return nil, fmt.Errorf("failed to add player: %w", err)
	}

	// Сохраняем обновленную игру
	if err := uc.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(game), nil
}

// toResponse converts the domain entity into a response DTO
func (uc *AddPlayerUseCase) toResponse(game *entities.Game) *dto.GameResponse {
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
