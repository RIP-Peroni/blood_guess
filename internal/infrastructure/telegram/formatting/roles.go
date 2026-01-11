package formatting

import "RIP-Peroni/blood_guess/internal/domain/value_objects"

// RoleEmoji returns emoji for the role
func RoleEmoji(role string) string {
	switch role {
	case "demon":
		return "👹"
	case "minion":
		return "😈"
	case "townsfolk":
		return "👨‍🌾"
	case "outsider":
		return "🚶"
	default:
		return "❓"
	}
}

// RoleDisplayName returns the display name of the role in Russian
func RoleDisplayName(role string) string {
	switch role {
	case "demon":
		return "демон"
	case "minion":
		return "приспешник"
	case "townsfolk":
		return "горожанин"
	case "outsider":
		return "изгой"
	default:
		return role
	}
}

// IsEvilRole проверяет, является ли роль злой
func IsEvilRole(role string) bool {
	r, err := value_objects.FromString(role)
	if err != nil {
		return false
	}
	return r.IsEvil()
}

// IsPeacefulRole проверяет, является ли роль мирной
func IsPeacefulRole(role string) bool {
	r, err := value_objects.FromString(role)
	if err != nil {
		return false
	}
	return r.IsPeaceful()
}
