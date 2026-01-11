package value_objects

// PredictableRoles contains roles that can be predicted
var PredictableRoles = map[Role]bool{
	RoleDemon:  true,
	RoleMinion: true,
}

// IsPredictable checks whether this role can be predicted
func IsPredictable(role Role) bool {
	return role.IsEvil()
}

// GetPredictableRoles returns a list of roles that can be predicted
func GetPredictableRoles() []Role {
	roles := make([]Role, 0, len(PredictableRoles))
	for role := range PredictableRoles {
		roles = append(roles, role)
	}
	return roles
}
