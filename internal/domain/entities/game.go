package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoPlayers          = errors.New("game must have at least one player")
	ErrPlayerAlreadyAdded = errors.New("player already added to this game")
	ErrEmptyPlayerName    = errors.New("player name cannot be empty")
	ErrInvalidGameState   = errors.New("invalid game state for this operation")
	ErrPlayerNotFound     = errors.New("player not found")
	ErrEmptyRole          = errors.New("role cannot be empty")
)

type GameID string
type PlayerSlotID string

type GameStatus string

const (
	GameStatusCreated           GameStatus = "created"
	GameStatusPredictionsOpen   GameStatus = "predictions_open"
	GameStatusPredictionsClosed GameStatus = "predictions_closed"
	GameStatusInProgress        GameStatus = "in_progress"
	GameStatusFinished          GameStatus = "finished"
)

type PlayerSlot struct {
	ID            PlayerSlotID
	UserID        int64
	Name          string
	AssignedRole  string
	RealRole      string
	IsRealRoleSet bool
}

type Game struct {
	id        GameID
	name      string
	creatorID int64
	status    GameStatus
	players   []PlayerSlot
	createdAt time.Time
	startedAt *time.Time
	endedAt   *time.Time
}

func NewGame(name string, creatorID int64) *Game {
	return &Game{
		id:        GameID(uuid.New().String()),
		name:      name,
		creatorID: creatorID,
		status:    GameStatusCreated,
		players:   []PlayerSlot{},
		createdAt: time.Now(),
	}
}

func (g *Game) ID() GameID {
	return g.id
}

func (g *Game) Name() string {
	return g.name
}

func (g *Game) CreatorID() int64 {
	return g.creatorID
}

func (g *Game) Status() GameStatus {
	return g.status
}

func (g *Game) Players() []PlayerSlot {
	return g.players
}

func (g *Game) CreatedAt() time.Time {
	return g.createdAt
}

func (g *Game) AddPlayer(userID int64, name string, assignedRole string) error {
	if name == "" {
		return ErrEmptyPlayerName
	}

	for _, player := range g.players {
		if player.UserID == userID {
			return ErrPlayerAlreadyAdded
		}
	}

	player := PlayerSlot{
		ID:            PlayerSlotID(uuid.New().String()),
		UserID:        userID,
		Name:          name,
		AssignedRole:  assignedRole,
		RealRole:      "",
		IsRealRoleSet: false,
	}

	g.players = append(g.players, player)
	return nil
}

func (g *Game) OpenPredictions() error {
	if g.status != GameStatusCreated {
		return ErrInvalidGameState
	}

	if len(g.players) == 0 {
		return ErrNoPlayers
	}

	g.status = GameStatusPredictionsOpen
	return nil
}

func (g *Game) ClosePredictions() error {
	if g.status != GameStatusPredictionsOpen {
		return ErrInvalidGameState
	}

	g.status = GameStatusPredictionsClosed
	return nil
}

func (g *Game) Start() error {
	if g.status != GameStatusPredictionsClosed {
		return ErrInvalidGameState
	}

	now := time.Now()
	g.startedAt = &now
	g.status = GameStatusInProgress
	return nil
}

func (g *Game) Finish() error {
	if g.status != GameStatusInProgress {
		return ErrInvalidGameState
	}

	now := time.Now()
	g.endedAt = &now
	g.status = GameStatusFinished
	return nil
}

func (g *Game) CanAcceptPredictions() bool {
	return g.status == GameStatusPredictionsOpen
}

func (g *Game) SetPlayerRealRole(playerID PlayerSlotID, realRole string) error {
	if realRole == "" {
		return ErrEmptyRole
	}

	for i, player := range g.players {
		if player.ID == playerID {
			g.players[i].RealRole = realRole
			g.players[i].IsRealRoleSet = true
			return nil
		}
	}

	return ErrPlayerNotFound
}

// Вспомогательный метод для поиска игрока по ID
func (g *Game) FindPlayerByID(playerID PlayerSlotID) (*PlayerSlot, error) {
	for _, player := range g.players {
		if player.ID == playerID {
			return &player, nil
		}
	}
	return nil, ErrPlayerNotFound
}
