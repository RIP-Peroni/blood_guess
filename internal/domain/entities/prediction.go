package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PredictionID string

// Prediction - user prediction of a specific player's role
// PointsAwarded set AFTER the game is completed
type Prediction struct {
	id            PredictionID
	gameID        GameID
	userID        UserID
	playerSlotID  PlayerSlotID
	predictedRole string
	pointsAwarded *int // Nil - points not yet awarded, indicator to distinguish 0 points from "not yet counted"
	createdAt     time.Time
}

func NewPrediction(gameID GameID, userID UserID, playerSlotID PlayerSlotID, predictedRole string) *Prediction {
	return &Prediction{
		id:            PredictionID(uuid.New().String()),
		gameID:        gameID,
		userID:        userID,
		playerSlotID:  playerSlotID,
		predictedRole: predictedRole,
		pointsAwarded: nil, // The points will be set later.
		createdAt:     time.Now(),
	}
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

func (p *Prediction) PredictedRole() string {
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

func (p *Prediction) IsCorrect(realRole string) bool {
	return p.predictedRole == realRole
}

// HasPointsAwarded - checks whether points have already been awarded
func (p *Prediction) HasPointsAwarded() bool {
	return p.pointsAwarded != nil
}
