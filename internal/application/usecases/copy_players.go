package usecases

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"fmt"
	"sort"
)

// CopyPlayersUseCase реализует копирование игроков из последней игры
type CopyPlayersUseCase struct {
	gameRepo ports.GameRepository
}

func NewCopyPlayersUseCase(gameRepo ports.GameRepository) *CopyPlayersUseCase {
	return &CopyPlayersUseCase{
		gameRepo: gameRepo,
	}
}

// findLastGameByCreator Finds the creator's latest game (except the one specified)
func (uc *CopyPlayersUseCase) findLastGameByCreator(creatorID int64, excludeGameID string) (*entities.Game, error) {
	games, err := uc.gameRepo.FindActiveGames()
	if err != nil {
		return nil, fmt.Errorf("failed to get active games: %w", err)
	}

	// Filter games by creator and exclude the current game
	var creatorGames []*entities.Game
	for _, game := range games {
		if game == nil {
			continue
		}

		// We check that the game has the correct ID and creator.
		if game.CreatorID() == creatorID && string(game.ID()) != excludeGameID {
			creatorGames = append(creatorGames, game)
		}
	}

	if len(creatorGames) == 0 {
		return nil, nil
	}

	// Sort by creation date (latest first)
	sort.Slice(creatorGames, func(i, j int) bool {
		return creatorGames[i].CreatedAt().After(creatorGames[j].CreatedAt())
	})

	return creatorGames[0], nil
}

// Execute copies players from the last game
func (uc *CopyPlayersUseCase) Execute(command dto.CopyPlayersCommand) (*dto.GameResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	targetGame, err := uc.gameRepo.FindByID(entities.GameID(command.TargetGameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	if targetGame.CreatorID() != command.AdminID {
		return nil, ErrNotGameCreator
	}

	if targetGame.Status() != entities.GameStatusCreated {
		return nil, fmt.Errorf("%w: current status is %s", ErrInvalidGameState, targetGame.Status())
	}

	sourceGame, err := uc.findLastGameByCreator(command.AdminID, command.TargetGameID)
	if err != nil {
		return nil, fmt.Errorf("failed to find previous game: %w", err)
	}

	if sourceGame == nil {
		return nil, fmt.Errorf("no previous games found for this creator")
	}

	if len(sourceGame.Players()) == 0 {
		return nil, fmt.Errorf("source game has no players")
	}

	if err := targetGame.CopyPlayersFrom(sourceGame); err != nil {
		return nil, fmt.Errorf("failed to copy players: %w", err)
	}

	if err := uc.gameRepo.Update(targetGame); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(targetGame), nil
}

// toResponse converts an entity into a DTO
func (uc *CopyPlayersUseCase) toResponse(game *entities.Game) *dto.GameResponse {
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
