package usecases

import (
	"fmt"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

// SetRealRoleUseCase implements the installation of the player's real role after the game
type SetRealRoleUseCase struct {
	gameRepo ports.GameRepository
}

func NewSetRealRoleUseCase(gameRepo ports.GameRepository) *SetRealRoleUseCase {
	return &SetRealRoleUseCase{
		gameRepo: gameRepo,
	}
}

// Execute sets the player's actual role after the game ends
func (uc *SetRealRoleUseCase) Execute(command dto.SetRealRoleCommand) error {
	if err := command.Validate(); err != nil {
		return fmt.Errorf("invalid command: %w", err)
	}

	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return ErrGameNotFound
	}

	if game.CreatorID() != command.AdminID {
		return ErrNotGameCreator
	}

	if game.Status() != entities.GameStatusInProgress && game.Status() != entities.GameStatusFinished {
		return fmt.Errorf("%w: can only set real roles during or after the game (current: %s)",
			ErrInvalidGameState, game.Status())
	}

	// НОВОЕ: Разрешаем только злые роли
	if command.RealRole != "demon" && command.RealRole != "minion" {
		return fmt.Errorf("можно устанавливать только злые роли (demon, minion). Мирные роли устанавливаются автоматически")
	}

	if err := game.SetPlayerRealRole(entities.PlayerSlotID(command.PlayerSlotID), command.RealRole); err != nil {
		return fmt.Errorf("failed to set real role: %w", err)
	}

	if err := uc.gameRepo.Update(game); err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	return nil
}
