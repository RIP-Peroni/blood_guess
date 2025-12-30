package dto

import (
	"errors"
	"strings"
)

var (
	ErrEmptyGameName      = errors.New("game name cannot be empty")
	ErrInvalidCreatorID   = errors.New("creator ID must be positive")
	ErrEmptyGameID        = errors.New("game ID cannot be empty")
	ErrEmptyUserID        = errors.New("user ID cannot be empty")
	ErrEmptyPlayerSlotID  = errors.New("player slot ID cannot be empty")
	ErrEmptyPredictedRole = errors.New("predicted role cannot be empty")
	ErrInvalidRole        = errors.New("invalid role")
	ErrEmptyAdminID       = errors.New("admin ID cannot be empty")
	ErrEmptyPlayerName    = errors.New("player name cannot be empty")
)

// List of allowed roles for the game "Blood on the Clocktower"
var validRoles = map[string]bool{
	"demon":     true,
	"minion":    true,
	"townsfolk": true,
	"outsider":  true,
}

// CreateGameCommand - command for creating a game
type CreateGameCommand struct {
	Name      string
	CreatorID int64
}

// Validate checks the correctness of the game creation command
func (c *CreateGameCommand) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrEmptyGameName
	}
	if c.CreatorID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// TrimmedName returns the game name without extra spaces
func (c *CreateGameCommand) TrimmedName() string {
	return strings.TrimSpace(c.Name)
}

// AddPlayerCommand - команда для добавления игрока
type AddPlayerCommand struct {
	GameID       string
	PlayerName   string
	AssignedRole string
	AdminID      int64
}

// Validate checks the correctness of the player add command
func (c *AddPlayerCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if strings.TrimSpace(c.PlayerName) == "" {
		return ErrEmptyPlayerName
	}
	if c.AssignedRole == "" {
		return ErrEmptyPredictedRole
	}
	if !IsValidRole(c.AssignedRole) {
		return ErrInvalidRole
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// OpenPredictionsCommand - command to open predictions
type OpenPredictionsCommand struct {
	GameID  string
	AdminID int64
}

// Validate checks the correctness of the forecast opening command
func (c *OpenPredictionsCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// SubmitPredictionCommand - command to send a prediction
type SubmitPredictionCommand struct {
	GameID        string
	UserID        string
	PlayerSlotID  string
	PredictedRole string
}

// Validate checks the correctness of the prediction sending command
func (c *SubmitPredictionCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.UserID == "" {
		return ErrEmptyUserID
	}
	if c.PlayerSlotID == "" {
		return ErrEmptyPlayerSlotID
	}
	if c.PredictedRole == "" {
		return ErrEmptyPredictedRole
	}
	if !c.IsValidRole() {
		return ErrInvalidRole
	}
	return nil
}

// IsValidRole checks if a role is valid
func (c *SubmitPredictionCommand) IsValidRole() bool {
	return IsValidRole(c.PredictedRole)
}

// ClosePredictionsCommand - command to close predictions
type ClosePredictionsCommand struct {
	GameID  string
	AdminID int64
}

// Validate checks the correctness of the prediction closing command
func (c *ClosePredictionsCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// FinishGameCommand - command to end the game
type FinishGameCommand struct {
	GameID    string
	AdminID   int64
	RealRoles map[string]string // PlayerSlotID -> RealRole
}

// Validate checks the correctness of the game finishing command
func (c *FinishGameCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// SetPlayerRoleCommand - command to set the player's real role
type SetPlayerRoleCommand struct {
	GameID       string
	PlayerSlotID string
	RealRole     string
	AdminID      int64
}

// Validate checks the correctness of the role setting command
func (c *SetPlayerRoleCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.PlayerSlotID == "" {
		return ErrEmptyPlayerSlotID
	}
	if c.RealRole == "" {
		return ErrEmptyPredictedRole
	}
	if !IsValidRole(c.RealRole) {
		return ErrInvalidRole
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	return validRoles[role]
}

// GetValidRoles returns a list of valid roles
func GetValidRoles() []string {
	roles := make([]string, 0, len(validRoles))
	for role := range validRoles {
		roles = append(roles, role)
	}
	return roles
}

type StartGameCommand struct {
	GameID  string
	AdminID int64
}

func (c *StartGameCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}
