package entities

import (
	"errors"
	"fmt"
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
	Name          string
	RealRole      string // Actual role after the game (can be any of 4)
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

func (g *Game) StartedAt() *time.Time {
	return g.startedAt
}

func (g *Game) EndedAt() *time.Time {
	return g.endedAt
}

func (g *Game) AddPlayer(name string) error {
	if name == "" {
		return ErrEmptyPlayerName
	}

	// Проверяем уникальность имени в игре
	for _, player := range g.players {
		if player.Name == name {
			return ErrPlayerAlreadyAdded
		}
	}

	player := PlayerSlot{
		ID:            PlayerSlotID(uuid.New().String()),
		Name:          name,
		RealRole:      "",
		IsRealRoleSet: false,
	}

	g.players = append(g.players, player)
	return nil
}

func (g *Game) AddPlayers(names []string) error {
	for _, name := range names {
		if err := g.AddPlayer(name); err != nil {
			return fmt.Errorf("failed to add player %s: %w", name, err)
		}
	}
	return nil
}

// CopyPlayersFrom copies players from another game (without roles)
func (g *Game) CopyPlayersFrom(sourceGame *Game) error {
	if sourceGame == nil {
		return errors.New("source game cannot be nil")
	}

	if len(sourceGame.Players()) == 0 {
		return fmt.Errorf("source game has no players")
	}

	for _, player := range sourceGame.Players() {
		if player.Name == "" {
			continue
		}

		playerExists := false
		for _, existingPlayer := range g.players {
			if existingPlayer.Name == player.Name {
				playerExists = true
				break
			}
		}

		if playerExists {
			continue
		}

		if err := g.AddPlayer(player.Name); err != nil {
			return fmt.Errorf("failed to copy player %s: %w", player.Name, err)
		}
	}

	return nil
}

// GetPlayerNames returns a list of player names
func (g *Game) GetPlayerNames() []string {
	names := make([]string, len(g.players))
	for i, player := range g.players {
		names[i] = player.Name
	}
	return names
}

// HasPlayersWithRealRoles checks if players have real roles set
func (g *Game) HasPlayersWithRealRoles() bool {
	for _, player := range g.players {
		if player.IsRealRoleSet {
			return true
		}
	}
	return false
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

// SetPlayerRealRole sets the player's actual role (which was assigned in the game)
func (g *Game) SetPlayerRealRole(playerID PlayerSlotID, realRole string) error {
	if realRole == "" {
		return ErrEmptyRole
	}

	if g.status != GameStatusInProgress && g.status != GameStatusFinished {
		return fmt.Errorf("cannot set real role: game must be in progress or finished")
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

func (g *Game) FindPlayerByID(playerID PlayerSlotID) (*PlayerSlot, error) {
	for _, player := range g.players {
		if player.ID == playerID {
			return &player, nil
		}
	}
	return nil, ErrPlayerNotFound
}

// AllRealRolesSet checks if real roles are set for all players
func (g *Game) AllRealRolesSet() bool {
	for _, player := range g.players {
		if !player.IsRealRoleSet {
			return false
		}
	}
	return true
}
