package usecases

import (
	"fmt"
	"time"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"
)

// AwardPointsUseCase implements the "award points" use case
type AwardPointsUseCase struct {
	gameRepo       ports.GameRepository
	userRepo       ports.UserRepository
	predictionRepo ports.PredictionRepository
	scoringService services.ScoringRules
}

func NewAwardPointsUseCase(
	gameRepo ports.GameRepository,
	userRepo ports.UserRepository,
	predictionRepo ports.PredictionRepository,
	scoringService services.ScoringRules,
) *AwardPointsUseCase {
	return &AwardPointsUseCase{
		gameRepo:       gameRepo,
		userRepo:       userRepo,
		predictionRepo: predictionRepo,
		scoringService: scoringService,
	}
}

// Execute awards points to users based on their predictions
func (uc *AwardPointsUseCase) Execute(command dto.AwardPointsCommand) (*dto.FinishGameResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	// Get the game
	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	if game.CreatorID() != command.AdminID {
		return nil, ErrNotGameCreator
	}

	if game.Status() != entities.GameStatusFinished {
		return nil, fmt.Errorf("%w: current status is %s, expected FINISHED", ErrInvalidGameState, game.Status())
	}

	// Check if all real roles are set
	if !game.AllRealRolesSet() {
		var missingRoles []string
		for _, player := range game.Players() {
			if !player.IsRealRoleSet {
				missingRoles = append(missingRoles, player.Name)
			}
		}
		return nil, fmt.Errorf("not all real roles are set. Missing roles for players: %v", missingRoles)
	}

	// Get all predictions for this game
	predictions, err := uc.predictionRepo.FindByGame(entities.GameID(command.GameID))
	if err != nil {
		return nil, fmt.Errorf("failed to find predictions: %w", err)
	}

	// Build real roles map
	realRoles := make(map[string]string)
	for _, player := range game.Players() {
		realRoles[string(player.ID)] = player.RealRole
	}

	// Calculate points for each prediction
	pointsPerPrediction := uc.scoringService.CalculatePointsForEachPrediction(predictions, realRoles)

	// Group predictions by user and calculate total points per user
	userScores := make(map[entities.UserID]int)
	userPredictions := make(map[entities.UserID][]*entities.Prediction)

	for _, prediction := range predictions {
		userID := prediction.UserID()
		userPredictions[userID] = append(userPredictions[userID], prediction)
		points := pointsPerPrediction[prediction.ID()]
		userScores[userID] += points
	}

	// Award points to users and update predictions
	userScoreMap := make(map[string]int)
	for userID, totalScore := range userScores {
		// Get user
		user, err := uc.userRepo.FindById(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to find user %s: %w", userID, err)
		}

		// Change user balance
		_, err = user.ChangeBalance(totalScore)
		if err != nil {
			return nil, fmt.Errorf("failed to change balance for user %s: %w", userID, err)
		}

		// Update user
		if err := uc.userRepo.Update(user); err != nil {
			return nil, fmt.Errorf("failed to update user %s: %w", userID, err)
		}

		userScoreMap[string(userID)] = totalScore

		// Award points to each prediction
		for _, prediction := range userPredictions[userID] {
			points := pointsPerPrediction[prediction.ID()]
			if err := prediction.AwardPoints(points); err != nil {
				return nil, fmt.Errorf("failed to award points to prediction %s: %w", prediction.ID(), err)
			}

			if err := uc.predictionRepo.Update(prediction); err != nil {
				return nil, fmt.Errorf("failed to update prediction %s: %w", prediction.ID(), err)
			}
		}
	}

	return uc.toResponse(game, userScoreMap), nil
}

// toResponse converts domain entities to response DTO
func (uc *AwardPointsUseCase) toResponse(game *entities.Game, userScores map[string]int) *dto.FinishGameResponse {
	endedAt := time.Now()
	if game.EndedAt() != nil {
		endedAt = *game.EndedAt()
	}

	return &dto.FinishGameResponse{
		GameID:          string(game.ID()),
		Name:            game.Name(),
		Status:          game.Status(),
		PlayerResults:   []dto.PlayerResultResponse{},
		UserScores:      userScores,
		CurrencyAwarded: true,
		EndedAt:         endedAt,
	}
}
