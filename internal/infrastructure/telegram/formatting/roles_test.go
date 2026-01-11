package formatting

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoleEmoji(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected string
	}{
		{"demon role", "demon", "👹"},
		{"minion role", "minion", "😈"},
		{"townsfolk role", "townsfolk", "👨‍🌾"},
		{"outsider role", "outsider", "🚶"},
		{"unknown role", "unknown", "❓"},
		{"empty role", "", "❓"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoleEmoji(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRoleDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected string
	}{
		{"demon role", "demon", "демон"},
		{"minion role", "minion", "приспешник"},
		{"townsfolk role", "townsfolk", "горожанин"},
		{"outsider role", "outsider", "изгой"},
		{"unknown role", "unknown", "unknown"},
		{"empty role", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoleDisplayName(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsEvilRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"demon is evil", "demon", true},
		{"minion is evil", "minion", true},
		{"townsfolk is not evil", "townsfolk", false},
		{"outsider is not evil", "outsider", false},
		{"unknown is not evil", "unknown", false},
		{"empty is not evil", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEvilRole(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsPeacefulRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"demon is not peaceful", "demon", false},
		{"minion is not peaceful", "minion", false},
		{"townsfolk is peaceful", "townsfolk", true},
		{"outsider is peaceful", "outsider", true},
		{"unknown is not peaceful", "unknown", false},
		{"empty is not peaceful", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPeacefulRole(tt.role)
			assert.Equal(t, tt.expected, result)
		})
	}
}
