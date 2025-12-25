package value_objects

import "errors"

var (
	ErrInvalidRole = errors.New("invalid role")
)

type Role string

const (
	RoleDemon     Role = "demon"
	RoleMinion    Role = "minion"
	RoleTownsfolk Role = "townsfolk"
	RoleOutsider  Role = "outsider"
)

// IsValid checks if a role is valid
func (r Role) IsValid() bool {
	switch r {
	case RoleDemon, RoleMinion, RoleTownsfolk, RoleOutsider:
		return true
	default:
		return false
	}
}

// String returns a string representation of the role
func (r Role) String() string {
	return string(r)
}

// FromString creates a Role from a string with validation
func FromString(roleStr string) (Role, error) {
	role := Role(roleStr)
	if !role.IsValid() {
		return "", ErrInvalidRole
	}
	return role, nil
}
