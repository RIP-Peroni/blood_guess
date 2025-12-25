package services

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/value_objects"
)

type BasicScoringRules struct {
	pointsForDemon   int
	pointsForMinion  int
	pointsForGood    int
	penaltyForDemon  int
	penaltyForMinion int
	penaltyForGood   int
}

func NewBasicScoringRules() *BasicScoringRules {
	return &BasicScoringRules{
		pointsForDemon:   10,
		pointsForMinion:  5,
		pointsForGood:    2,
		penaltyForDemon:  3,
		penaltyForMinion: 2,
		penaltyForGood:   1,
	}
}

// CalculatePointsForUser calculates points for all user predictions
// Returns the total score
func (s *BasicScoringRules) CalculatePointsForUser(predictions []*entities.Prediction, realRoles map[string]string) int {
	totalScore := 0

	for _, prediction := range predictions {
		realRole, exists := realRoles[string(prediction.PlayerSlotID())]
		if !exists {
			continue
		}

		if prediction.IsCorrect(realRole) {
			totalScore += s.getPointsForRole(realRole)
		} else {
			totalScore -= s.getPenaltyForRole(prediction.PredictedRole())
		}
	}

	return totalScore
}

// CalculatePointsForEachPrediction calculates points for each prediction separately
// Useful for recording points in each prediction
func (s *BasicScoringRules) CalculatePointsForEachPrediction(predictions []entities.Prediction, realRoles map[string]string) map[entities.PredictionID]int {
	result := make(map[entities.PredictionID]int)

	for _, prediction := range predictions {
		realRole, exists := realRoles[string(prediction.PlayerSlotID())]
		if !exists {
			result[prediction.ID()] = 0
			continue
		}

		if prediction.IsCorrect(realRole) {
			result[prediction.ID()] = s.getPointsForRole(realRole)
		} else {
			result[prediction.ID()] = -s.getPenaltyForRole(prediction.PredictedRole())
		}
	}

	return result
}

func (s *BasicScoringRules) getPointsForRole(role string) int {
	switch role {
	case "demon":
		return s.pointsForDemon
	case "minion":
		return s.pointsForMinion
	default: // townsfolk, outsider, and other "good" roles
		return s.pointsForGood
	}
}

func (s *BasicScoringRules) getPenaltyForRole(role value_objects.Role) int {
	switch role {
	case "demon":
		return s.penaltyForDemon
	case "minion":
		return s.penaltyForMinion
	default: // townsfolk, outsider
		return s.penaltyForGood
	}
}
