package usecases

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"errors"
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

	// Find the last game in FINISHED status
	finishedGames, err := uc.gameRepo.FindByStatus(entities.GameStatusFinished)
	if err != nil {
		return nil, fmt.Errorf("failed to find finished games: %w", err)
	}

	if len(finishedGames) == 0 {
		return nil, fmt.Errorf("no finished games found")
	}

	// Sort by creation date (latest first)
	sort.Slice(finishedGames, func(i, j int) bool {
		return finishedGames[i].CreatedAt().After(finishedGames[j].CreatedAt())
	})

	sourceGame := finishedGames[0]

	if len(sourceGame.Players()) == 0 {
		return nil, fmt.Errorf("source game has no players")
	}

	for _, player := range sourceGame.Players() {
		if player.Name == "" {
			continue
		}

		playerExists := false
		for _, existingPlayer := range targetGame.Players() {
			if existingPlayer.Name == player.Name {
				playerExists = true
				break
			}
		}

		if playerExists {
			continue
		}

		if err := targetGame.AddPlayer(player.Name); err != nil {
			if errors.Is(err, entities.ErrPlayerAlreadyAdded) {
				continue
			}
			return nil, fmt.Errorf("failed to copy player %s: %w", player.Name, err)
		}
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
