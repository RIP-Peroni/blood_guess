package dto

import (
	"RIP-Peroni/blood_guess/internal/domain/value_objects"
	"errors"
	"fmt"
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
	ErrEmptyPlayerName    = errors.New("player name cannot be empty")
	ErrEmptyPlayerNames   = errors.New("player names cannot be empty")
)

// List of allowed roles for the game "Blood on the Clocktower"
var (
	allGameRoles = map[string]bool{
		"demon":     true,
		"minion":    true,
		"townsfolk": true,
		"outsider":  true,
	}

	validRoles = map[string]bool{
		"demon":     true,
		"minion":    true,
		"townsfolk": true,
		"outsider":  true,
	}
)

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
	GameID     string
	PlayerName string
	AdminID    int64
}

// Validate checks the correctness of the player add command
func (c *AddPlayerCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if strings.TrimSpace(c.PlayerName) == "" {
		return ErrEmptyPlayerName
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
	if !IsPredictableRole(c.PredictedRole) {
		return fmt.Errorf("%w: only demon and minion roles can be predicted", ErrInvalidRole)
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

// IsValidRole checks if a role is valid
func IsValidRole(role string) bool {
	return validRoles[role]
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

// AddPlayersCommand - команда для добавления нескольких игроков
type AddPlayersCommand struct {
	GameID      string
	PlayerNames []string
	AdminID     int64
}

// Validate checks the command to add multiple players
func (c *AddPlayersCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if len(c.PlayerNames) == 0 {
		return ErrEmptyPlayerNames
	}
	for _, name := range c.PlayerNames {
		if strings.TrimSpace(name) == "" {
			return ErrEmptyPlayerName
		}
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// CopyPlayersCommand - command to copy players from a previous game
type CopyPlayersCommand struct {
	TargetGameID string
	AdminID      int64
}

// Validate checks the player copy command
func (c *CopyPlayersCommand) Validate() error {
	if c.TargetGameID == "" {
		return ErrEmptyGameID
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// IsValidGameRole checks if a role is valid for the game
func IsValidGameRole(role string) bool {
	return allGameRoles[role]
}

// GetValidGameRoles returns all valid game roles
func GetValidGameRoles() []string {
	roles := make([]string, 0, len(allGameRoles))
	for role := range allGameRoles {
		roles = append(roles, role)
	}
	return roles
}

// SubmitPredictionCommand - command to send a prediction
type SubmitPredictionCommand struct {
	GameID        string
	UserID        string
	PlayerSlotID  string
	PredictedRole string
}

// SetRealRoleCommand - command to set the player's real role after the game
type SetRealRoleCommand struct {
	GameID       string
	PlayerSlotID string
	RealRole     string
	AdminID      int64
}

// Validate checks the real role setting command
func (c *SetRealRoleCommand) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.PlayerSlotID == "" {
		return ErrEmptyPlayerSlotID
	}
	if c.RealRole == "" {
		return ErrEmptyPredictedRole
	}
	if !IsValidGameRole(c.RealRole) {
		return fmt.Errorf("%w: valid roles are demon, minion, townsfolk, outsider", ErrInvalidRole)
	}
	if c.AdminID <= 0 {
		return ErrInvalidCreatorID
	}
	return nil
}

// IsPredictableRole checks whether this role can be predicted
func IsPredictableRole(role string) bool {
	r, err := value_objects.FromString(role)
	if err != nil {
		return false
	}
	return value_objects.IsPredictable(r)
}
