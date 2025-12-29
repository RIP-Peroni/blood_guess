package usecases

import (
	"fmt"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
)

// SubmitPredictionUseCase implements use case "submitting prediction"
type SubmitPredictionUseCase struct {
	gameRepo       ports.GameRepository
	userRepo       ports.UserRepository
	predictionRepo ports.PredictionRepository
}

func NewSubmitPredictionUseCase(
	gameRepo ports.GameRepository,
	userRepo ports.UserRepository,
	predictionRepo ports.PredictionRepository,
) *SubmitPredictionUseCase {
	return &SubmitPredictionUseCase{
		gameRepo:       gameRepo,
		userRepo:       userRepo,
		predictionRepo: predictionRepo,
	}
}

// Execute executes submission of prediction
func (uc *SubmitPredictionUseCase) Execute(command dto.SubmitPredictionCommand) (*dto.PredictionResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command format: %w", err)
	}

	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	if !game.CanAcceptPredictions() {
		return nil, ErrGameNotAcceptingPredictions
	}

	_, err = uc.userRepo.FindById(entities.UserID(command.UserID))
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !uc.playerSlotExists(game, command.PlayerSlotID) {
		return nil, ErrPlayerSlotNotFound
	}

	existingPredictions, err := uc.predictionRepo.FindByGameAndUser(
		entities.GameID(command.GameID),
		entities.UserID(command.UserID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing predictions: %w", err)
	}

	for _, existing := range existingPredictions {
		if string(existing.PlayerSlotID()) == command.PlayerSlotID {
			return nil, fmt.Errorf("%w: user %s already predicted for slot %s",
				ErrAlreadyPredicted, command.UserID, command.PlayerSlotID)
		}
	}

	prediction, err := entities.NewPrediction(
		entities.GameID(command.GameID),
		entities.UserID(command.UserID),
		entities.PlayerSlotID(command.PlayerSlotID),
		command.PredictedRole,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create prediction: %w", err)
	}

	if err := uc.predictionRepo.Save(prediction); err != nil {
		return nil, fmt.Errorf("failed to save prediction: %w", err)
	}

	return uc.toResponse(prediction), nil
}

// playerSlotExists checks if the player slot exists in the game
func (uc *SubmitPredictionUseCase) playerSlotExists(game *entities.Game, playerSlotID string) bool {
	for _, player := range game.Players() {
		if string(player.ID) == playerSlotID {
			return true
		}
	}
	return false
}

// toResponse converts a domain entity into a response DTO
func (uc *SubmitPredictionUseCase) toResponse(prediction *entities.Prediction) *dto.PredictionResponse {
	points, awarded := prediction.PointsAwarded()

	return &dto.PredictionResponse{
		ID:            string(prediction.ID()),
		GameID:        string(prediction.GameID()),
		UserID:        string(prediction.UserID()),
		PlayerSlotID:  string(prediction.PlayerSlotID()),
		PredictedRole: string(prediction.PredictedRole()),
		Points:        points,
		PointsAwarded: awarded,
		CreatedAt:     prediction.CreatedAt(),
	}
}
