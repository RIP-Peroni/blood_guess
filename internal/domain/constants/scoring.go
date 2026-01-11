package constants

const (
	PointsForDemon   = 10
	PointsForMinion  = 5
	PenaltyForDemon  = 3
	PenaltyForMinion = 2
)

// PointsForRole returns points for correctly guessing the role
func PointsForRole(role string) int {
	switch role {
	case "demon":
		return PointsForDemon
	case "minion":
		return PointsForMinion
	default:
		return 0
	}
}

// PenaltyForRole returns the penalty for guessing the role incorrectly
func PenaltyForRole(role string) int {
	switch role {
	case "demon":
		return PenaltyForDemon
	case "minion":
		return PenaltyForMinion
	default:
		return 0
	}
}
