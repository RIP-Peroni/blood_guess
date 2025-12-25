package entities

import (
	"errors"
	"time"

	"RIP-Peroni/blood_guess/internal/domain/value_objects"

	"github.com/google/uuid"
)

type PredictionID string

type Prediction struct {
	id            PredictionID
	gameID        GameID
	userID        UserID
	playerSlotID  PlayerSlotID
	predictedRole value_objects.Role // Используем Role вместо string
	pointsAwarded *int
	createdAt     time.Time
}

func NewPrediction(gameID GameID, userID UserID, playerSlotID PlayerSlotID, predictedRoleStr string) (*Prediction, error) {
	role, err := value_objects.FromString(predictedRoleStr)
	if err != nil {
		return nil, err
	}

	return &Prediction{
		id:            PredictionID(uuid.New().String()),
		gameID:        gameID,
		userID:        userID,
		playerSlotID:  playerSlotID,
		predictedRole: role,
		pointsAwarded: nil,
		createdAt:     time.Now(),
	}, nil
}

func (p *Prediction) ID() PredictionID {
	return p.id
}

func (p *Prediction) GameID() GameID {
	return p.gameID
}

func (p *Prediction) UserID() UserID {
	return p.userID
}

func (p *Prediction) PlayerSlotID() PlayerSlotID {
	return p.playerSlotID
}

func (p *Prediction) PredictedRole() value_objects.Role {
	return p.predictedRole
}

// PointsAwarded returns the accrued points if they have already been calculated
func (p *Prediction) PointsAwarded() (int, bool) {
	if p.pointsAwarded == nil {
		return 0, false
	}
	return *p.pointsAwarded, true
}

func (p *Prediction) CreatedAt() time.Time {
	return p.createdAt
}

// AwardPoints awards points for this prediction
// Can only be called once
func (p *Prediction) AwardPoints(points int) error {
	if p.pointsAwarded != nil {
		return errors.New("points already awarded for this prediction")
	}
	p.pointsAwarded = &points
	return nil
}

func (p *Prediction) IsCorrect(realRoleStr string) bool {
	realRole, err := value_objects.FromString(realRoleStr)
	if err != nil {
		return false
	}
	return p.predictedRole == realRole
}

// HasPointsAwarded - checks whether points have already been awarded
func (p *Prediction) HasPointsAwarded() bool {
	return p.pointsAwarded != nil
}
