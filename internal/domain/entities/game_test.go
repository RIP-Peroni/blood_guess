package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGame_NewGame(t *testing.T) {
	t.Run("creating new game", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		assert.NotEmpty(t, game.ID())
		assert.Equal(t, "Test game", game.Name())
		assert.Equal(t, int64(12345), game.CreatorID())
		assert.Equal(t, GameStatusCreated, game.Status())
		assert.Empty(t, game.Players())
		assert.WithinDuration(t, time.Now(), game.CreatedAt(), time.Second)
	})
}

func TestGame_AddPlayer(t *testing.T) {
	t.Run("adding player to game", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)
		assert.Len(t, game.Players(), 1)

		player := game.Players()[0]
		assert.Equal(t, int64(11111), player.UserID)
		assert.Equal(t, "player 1", player.Name)
		assert.Equal(t, "townsfolk", player.AssignedRole)
		assert.Empty(t, player.RealRole)
		assert.False(t, player.IsRealRoleSet)
	})
	t.Run("can't add a player with empty name", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		err := game.AddPlayer(11111, "", "townsfolk")
		assert.Error(t, err)
		assert.Len(t, game.Players(), 0)
	})
	t.Run("can't add same player twice", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)

		err = game.AddPlayer(11111, "player 1", "outsider")
		assert.Error(t, err)
		assert.Len(t, game.Players(), 1)
	})
}

func TestGame_OpenPredictions(t *testing.T) {
	t.Run("successfully open predictions", func(t *testing.T) {
		game := NewGame("Test game", 12345)
		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)

		err = game.OpenPredictions()
		assert.NoError(t, err)
		assert.Equal(t, GameStatusPredictionsOpen, game.Status())
	})
	t.Run("can't open predictions without players", func(t *testing.T) {
		game := NewGame("Test game", 12345)

		err := game.OpenPredictions()
		assert.Error(t, err)
		assert.Equal(t, GameStatusCreated, game.Status())
	})
	t.Run("can't open predictions for finished game", func(t *testing.T) {
		game := NewGame("Test game", 12345)
		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)

		err = game.OpenPredictions()
		assert.NoError(t, err)

		err = game.ClosePredictions()
		assert.NoError(t, err)

		err = game.Start()
		assert.NoError(t, err)

		err = game.Finish()
		assert.NoError(t, err)

		err = game.OpenPredictions()
		assert.Error(t, err)
		assert.Equal(t, GameStatusFinished, game.Status())
	})
}

func TestGame_ClosePredictions(t *testing.T) {
	t.Run("successfully close predictions", func(t *testing.T) {
		game := NewGame("Test game", 12345)
		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)

		err = game.OpenPredictions()
		assert.NoError(t, err)

		err = game.ClosePredictions()
		assert.NoError(t, err)
		assert.Equal(t, GameStatusPredictionsClosed, game.Status())
	})
	t.Run("can't close predictions if they are not opened", func(t *testing.T) {
		game := NewGame("Test game", 12345)
		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)

		err = game.ClosePredictions()
		assert.Error(t, err)
		assert.Equal(t, GameStatusCreated, game.Status())
	})
}

func TestGame_CanAcceptPredictions(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Game)
		expected bool
	}{
		{
			name:     "The game is created but predictions are not opened - it is impossible to accept a prediction",
			setup:    func(g *Game) {},
			expected: false,
		},
		{
			name: "Predictions are opened - game can accept a prediction",
			setup: func(g *Game) {
				err := g.AddPlayer(11111, "player 1", "townsfolk")
				assert.NoError(t, err)

				err = g.OpenPredictions()
				assert.NoError(t, err)

			},
			expected: true,
		},
		{
			name: "Predictions are closed - game can not accept a prediction",
			setup: func(g *Game) {
				err := g.AddPlayer(11111, "player 1", "townsfolk")
				assert.NoError(t, err)

				err = g.OpenPredictions()
				assert.NoError(t, err)

				err = g.ClosePredictions()
				assert.NoError(t, err)

			},
			expected: false,
		},
		{
			name: "game finished - can't accept a prediction",
			setup: func(g *Game) {
				err := g.AddPlayer(11111, "player 1", "townsfolk")
				assert.NoError(t, err)

				err = g.OpenPredictions()
				assert.NoError(t, err)

				err = g.ClosePredictions()
				assert.NoError(t, err)

				err = g.Start()
				assert.NoError(t, err)

				err = g.Finish()
				assert.NoError(t, err)

			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame("Test Game", 12345)
			tt.setup(game)

			result := game.CanAcceptPredictions()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGame_SetPlayerRealRole(t *testing.T) {
	t.Run("set real role to a player", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)

		playerID := game.Players()[0].ID
		err = game.SetPlayerRealRole(playerID, "demon")

		assert.NoError(t, err)
		assert.Equal(t, "demon", game.Players()[0].RealRole)
		assert.True(t, game.Players()[0].IsRealRoleSet)
	})
	t.Run("can't set a role to unexisting player", func(t *testing.T) {
		game := NewGame("Test Game", 12345)

		err := game.SetPlayerRealRole("non-existing-id", "demon")
		assert.Error(t, err)
	})
	t.Run("can't set an empty role", func(t *testing.T) {
		game := NewGame("Test Game", 12345)
		err := game.AddPlayer(11111, "player 1", "townsfolk")
		assert.NoError(t, err)

		playerID := game.Players()[0].ID
		err = game.SetPlayerRealRole(playerID, "")

		assert.Error(t, err)
		assert.False(t, game.Players()[0].IsRealRoleSet)
	})
}
