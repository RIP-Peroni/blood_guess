package usecases

import (
	"RIP-Peroni/blood_guess/internal/domain/constants"
	"fmt"
	"time"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/services"
)

// FinishGameUseCase implements the "finish game" use case
type FinishGameUseCase struct {
	gameRepo       ports.GameRepository
	userRepo       ports.UserRepository
	predictionRepo ports.PredictionRepository
	scoringService services.ScoringRules
}

func NewFinishGameUseCase(
	gameRepo ports.GameRepository,
	userRepo ports.UserRepository,
	predictionRepo ports.PredictionRepository,
	scoringService services.ScoringRules,
) *FinishGameUseCase {
	return &FinishGameUseCase{
		gameRepo:       gameRepo,
		userRepo:       userRepo,
		predictionRepo: predictionRepo,
		scoringService: scoringService,
	}
}

// Execute finishes the game and calculates results
func (uc *FinishGameUseCase) Execute(command dto.FinishGameCommand) (*dto.FinishGameResponse, error) {
	if err := command.Validate(); err != nil {
		return nil, fmt.Errorf("invalid command: %w", err)
	}

	// Get the game
	game, err := uc.gameRepo.FindByID(entities.GameID(command.GameID))
	if err != nil {
		return nil, ErrGameNotFound
	}

	// Check if user is the game creator
	if game.CreatorID() != command.AdminID {
		return nil, ErrNotGameCreator
	}

	// Check if game is in correct state (in progress)
	if game.Status() != entities.GameStatusInProgress {
		return nil, fmt.Errorf("%w: current status is %s", ErrInvalidGameState, game.Status())
	}

	// Check if all real roles are set
	if !uc.allRealRolesSet(game) {
		return nil, fmt.Errorf("cannot finish game: not all players have real roles set")
	}

	// Get all predictions for this game
	predictions, err := uc.predictionRepo.FindByGame(game.ID())
	if err != nil {
		return nil, fmt.Errorf("failed to get predictions: %w", err)
	}

	// Prepare real roles map
	realRoles := make(map[string]string)
	for _, player := range game.Players() {
		realRoles[string(player.ID)] = player.RealRole
	}

	// Group predictions by user
	predictionsByUser := uc.groupPredictionsByUser(predictions)

	// Calculate points for each user
	userScores := make(map[string]int)
	userPredictions := make(map[string][]*entities.Prediction)

	for userID, userPreds := range predictionsByUser {
		// Convert to slice of entities.Prediction (not pointers)
		totalScore := uc.scoringService.CalculatePointsForUser(userPreds, realRoles)
		userScores[userID] = totalScore
		userPredictions[userID] = userPreds
	}

	// Award points to users and update predictions
	for userID, score := range userScores {
		// Find user
		user, err := uc.userRepo.FindById(entities.UserID(userID))
		if err != nil {
			// If user not found, skip (should not happen)
			continue
		}

		// Award currency to user (balance won't go negative)
		if score != 0 {
			_, err := user.ChangeBalance(score)
			if err != nil {
				return nil, fmt.Errorf("failed to update user balance: %w", err)
			}

			// Update user
			if err := uc.userRepo.Update(user); err != nil {
				return nil, fmt.Errorf("failed to save user: %w", err)
			}
		}

		// Award points to each prediction (for record keeping, can be negative)
		userPreds := userPredictions[userID]
		for _, pred := range userPreds {
			// Calculate points for this specific prediction
			realRole := realRoles[string(pred.PlayerSlotID())]
			var points int
			if pred.IsCorrect(realRole) {
				points = constants.PointsForRole(realRole)
			} else {
				points = -constants.PenaltyForRole(pred.PredictedRole().String())
			}

			// Award points to prediction (can be negative)
			if err := pred.AwardPoints(points); err != nil {
				return nil, fmt.Errorf("failed to award points to prediction: %w", err)
			}

			// Update prediction
			if err := uc.predictionRepo.Update(pred); err != nil {
				return nil, fmt.Errorf("failed to update prediction: %w", err)
			}
		}
	}

	// Finish the game
	if err := game.Finish(); err != nil {
		return nil, fmt.Errorf("failed to finish game: %w", err)
	}

	// Save the updated game
	if err := uc.gameRepo.Update(game); err != nil {
		return nil, fmt.Errorf("failed to update game: %w", err)
	}

	return uc.toResponse(game, userScores, predictionsByUser), nil
}

// allRealRolesSet checks if all players have real roles set
func (uc *FinishGameUseCase) allRealRolesSet(game *entities.Game) bool {
	for _, player := range game.Players() {
		if !player.IsRealRoleSet {
			return false
		}
	}
	return true
}

// groupPredictionsByUser groups predictions by user ID
func (uc *FinishGameUseCase) groupPredictionsByUser(predictions []*entities.Prediction) map[string][]*entities.Prediction {
	result := make(map[string][]*entities.Prediction)

	for _, prediction := range predictions {
		userID := string(prediction.UserID())
		result[userID] = append(result[userID], prediction)
	}

	return result
}

// toResponse converts domain entities to response DTO
func (uc *FinishGameUseCase) toResponse(
	game *entities.Game,
	userScores map[string]int,
	predictionsByUser map[string][]*entities.Prediction,
) *dto.FinishGameResponse {
	// Build player results
	playerResults := make([]dto.PlayerResultResponse, 0, len(game.Players()))

	for _, player := range game.Players() {
		// Find predictions for this player
		playerPredictions := make([]dto.PredictionResultResponse, 0)

		for userID, preds := range predictionsByUser {
			for _, pred := range preds {
				if string(pred.PlayerSlotID()) == string(player.ID) {
					// Get user for username
					user, err := uc.userRepo.FindById(entities.UserID(userID))
					username := ""
					if err == nil {
						username = user.Username()
					}

					// Check if prediction is correct
					isCorrect := pred.IsCorrect(player.RealRole)

					// Get points for this prediction
					points := 0
					if pointsAwarded, ok := pred.PointsAwarded(); ok {
						points = pointsAwarded
					}

					playerPredictions = append(playerPredictions, dto.PredictionResultResponse{
						UserID:        userID,
						Username:      username,
						PredictedRole: pred.PredictedRole().String(),
						Points:        points,
						IsCorrect:     isCorrect,
					})
				}
			}
		}

		playerResults = append(playerResults, dto.PlayerResultResponse{
			PlayerID:    string(player.ID),
			PlayerName:  player.Name,
			RealRole:    player.RealRole,
			Predictions: playerPredictions,
		})
	}

	// Get endedAt time
	endedAt := time.Now()
	if game.EndedAt() != nil {
		endedAt = *game.EndedAt()
	}

	return &dto.FinishGameResponse{
		GameID:          string(game.ID()),
		Name:            game.Name(),
		Status:          game.Status(),
		PlayerResults:   playerResults,
		UserScores:      userScores,
		CurrencyAwarded: len(userScores) > 0,
		EndedAt:         endedAt,
	}
}
