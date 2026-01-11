package services

import (
	"RIP-Peroni/blood_guess/internal/domain/constants"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/domain/value_objects"
)

type BasicScoringRules struct {
	pointsForDemon   int
	pointsForMinion  int
	penaltyForDemon  int
	penaltyForMinion int
}

func NewBasicScoringRules() *BasicScoringRules {
	return &BasicScoringRules{
		pointsForDemon:   10,
		pointsForMinion:  5,
		penaltyForDemon:  3,
		penaltyForMinion: 2,
	}
}

// CalculatePointsForUser calculates points for all user predictions
func (s *BasicScoringRules) CalculatePointsForUser(predictions []*entities.Prediction, realRoles map[string]string) int {
	totalScore := 0

	for _, prediction := range predictions {
		realRole, exists := realRoles[string(prediction.PlayerSlotID())]
		if !exists {
			continue
		}

		predictedRole := prediction.PredictedRole()
		if predictedRole.String() == realRole {
			// Правильный прогноз
			totalScore += constants.PointsForRole(realRole)
		} else {
			// Неправильный прогноз - штраф
			totalScore -= constants.PenaltyForRole(predictedRole.String())
		}
	}

	return totalScore
}

// CalculatePointsForEachPrediction calculates points for each prediction separately
func (s *BasicScoringRules) CalculatePointsForEachPrediction(predictions []*entities.Prediction, realRoles map[string]string) map[entities.PredictionID]int {
	result := make(map[entities.PredictionID]int)

	for _, prediction := range predictions {
		realRole, exists := realRoles[string(prediction.PlayerSlotID())]
		if !exists {
			result[prediction.ID()] = 0
			continue
		}

		predictedRole := prediction.PredictedRole()

		if predictedRole.String() == realRole {
			result[prediction.ID()] = constants.PointsForRole(realRole)
		} else {
			result[prediction.ID()] = -constants.PenaltyForRole(predictedRole.String())
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
	default:
		return 0
	}
}

func (s *BasicScoringRules) getPenaltyForRole(role value_objects.Role) int {
	switch role.String() {
	case "demon":
		return s.penaltyForDemon
	case "minion":
		return s.penaltyForMinion
	default:
		return 0
	}
}

// isEvilRole checks if the role is evil
func (s *BasicScoringRules) isEvilRole(role value_objects.Role) bool {
	return role.String() == "demon" || role.String() == "minion"
}
