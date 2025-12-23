package services

import "RIP-Peroni/blood_guess/internal/domain/entities"

type ScoringRules interface {
	// CalculatePointsForUser calculates the total number of points for the user
	CalculatePointsForUser(predictions []entities.Prediction, realRoles map[string]string) int

	// CalculatePointsForEachPrediction calculates points for each prediction separately
	CalculatePointsForEachPrediction(predictions []entities.Prediction, realRoles map[string]string) map[entities.PredictionID]int
}
