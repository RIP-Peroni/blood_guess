package dto

import (
	"testing"
	"time"

	"RIP-Peroni/blood_guess/internal/domain/entities"

	"github.com/stretchr/testify/assert"
)

// ==================== TESTS FOR COMMANDS ====================

func TestCreateGameCommand(t *testing.T) {
	t.Run("creating a command with correct data", func(t *testing.T) {
		cmd := CreateGameCommand{
			Name:      "test game",
			CreatorID: 12345,
		}

		assert.Equal(t, "test game", cmd.Name)
		assert.Equal(t, int64(12345), cmd.CreatorID)
	})

	t.Run("validation of command", func(t *testing.T) {
		tests := []struct {
			name      string
			command   CreateGameCommand
			wantError bool
		}{
			{
				name: "valid command",
				command: CreateGameCommand{
					Name:      "test",
					CreatorID: 1,
				},
				wantError: false,
			},
			{
				name: "empty name",
				command: CreateGameCommand{
					Name:      "",
					CreatorID: 1,
				},
				wantError: true,
			},
			{
				name: "null CreatorID",
				command: CreateGameCommand{
					Name:      "test",
					CreatorID: 0,
				},
				wantError: true,
			},
			{
				name: "negative CreatorID",
				command: CreateGameCommand{
					Name:      "test",
					CreatorID: -1,
				},
				wantError: true,
			},
			{
				name: "only white spaces in name",
				command: CreateGameCommand{
					Name:      "   ",
					CreatorID: 1,
				},
				wantError: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.command.Validate()
				if tt.wantError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestSubmitPredictionCommand(t *testing.T) {
	t.Run("creation of command prediction", func(t *testing.T) {
		cmd := SubmitPredictionCommand{
			GameID:        "game-123",
			UserID:        "user-456",
			PlayerSlotID:  "slot-789",
			PredictedRole: "demon",
		}

		assert.Equal(t, "game-123", cmd.GameID)
		assert.Equal(t, "user-456", cmd.UserID)
		assert.Equal(t, "slot-789", cmd.PlayerSlotID)
		assert.Equal(t, "demon", cmd.PredictedRole)
	})

	t.Run("validation of command prediction", func(t *testing.T) {
		tests := []struct {
			name      string
			command   SubmitPredictionCommand
			wantError bool
		}{
			{
				name: "valid command",
				command: SubmitPredictionCommand{
					GameID:        "game-123",
					UserID:        "user-456",
					PlayerSlotID:  "slot-789",
					PredictedRole: "demon",
				},
				wantError: false,
			},
			{
				name: "empty GameID",
				command: SubmitPredictionCommand{
					GameID:        "",
					UserID:        "user-456",
					PlayerSlotID:  "slot-789",
					PredictedRole: "demon",
				},
				wantError: true,
			},
			{
				name: "empty UserID",
				command: SubmitPredictionCommand{
					GameID:        "game-123",
					UserID:        "",
					PlayerSlotID:  "slot-789",
					PredictedRole: "demon",
				},
				wantError: true,
			},
			{
				name: "empty PlayerSlotID",
				command: SubmitPredictionCommand{
					GameID:        "game-123",
					UserID:        "user-456",
					PlayerSlotID:  "",
					PredictedRole: "demon",
				},
				wantError: true,
			},
			{
				name: "empty PredictedRole",
				command: SubmitPredictionCommand{
					GameID:        "game-123",
					UserID:        "user-456",
					PlayerSlotID:  "slot-789",
					PredictedRole: "",
				},
				wantError: true,
			},
			{
				name: "incorrect role",
				command: SubmitPredictionCommand{
					GameID:        "game-123",
					UserID:        "user-456",
					PlayerSlotID:  "slot-789",
					PredictedRole: "invalid_role",
				},
				wantError: true, // If we have role validation
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.command.Validate()
				if tt.wantError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

func TestOpenPredictionsCommand(t *testing.T) {
	cmd := OpenPredictionsCommand{
		GameID:  "game-123",
		AdminID: 12345,
	}

	assert.Equal(t, "game-123", cmd.GameID)
	assert.Equal(t, int64(12345), cmd.AdminID)
}

func TestClosePredictionsCommand(t *testing.T) {
	cmd := ClosePredictionsCommand{
		GameID:  "game-123",
		AdminID: 12345,
	}

	assert.Equal(t, "game-123", cmd.GameID)
	assert.Equal(t, int64(12345), cmd.AdminID)
}

func TestFinishGameCommand(t *testing.T) {
	cmd := FinishGameCommand{
		GameID:  "game-123",
		AdminID: 12345,
		RealRoles: map[string]string{
			"slot-1": "demon",
			"slot-2": "minion",
		},
	}

	assert.Equal(t, "game-123", cmd.GameID)
	assert.Equal(t, int64(12345), cmd.AdminID)
	assert.Len(t, cmd.RealRoles, 2)
	assert.Equal(t, "demon", cmd.RealRoles["slot-1"])
	assert.Equal(t, "minion", cmd.RealRoles["slot-2"])
}

// ==================== TESTS FOR RESPONSES ====================

func TestGameResponse(t *testing.T) {
	t.Run("creating response with game", func(t *testing.T) {
		players := []PlayerResponse{
			{
				ID:           "player-1",
				Name:         "Player 1",
				AssignedRole: "townsfolk",
				RealRole:     "demon",
			},
		}

		response := GameResponse{
			ID:        "game-123",
			Name:      "test game",
			Status:    entities.GameStatusFinished,
			CreatorID: 12345,
			Players:   players,
			CreatedAt: time.Now(),
		}

		assert.Equal(t, "game-123", response.ID)
		assert.Equal(t, "test game", response.Name)
		assert.Equal(t, entities.GameStatusFinished, response.Status)
		assert.Equal(t, int64(12345), response.CreatorID)
		assert.Len(t, response.Players, 1)
		assert.Equal(t, "player-1", response.Players[0].ID)
		assert.NotZero(t, response.CreatedAt)
	})

	t.Run("game without players", func(t *testing.T) {
		response := GameResponse{
			ID:        "game-123",
			Name:      "Empty game",
			Status:    entities.GameStatusCreated,
			CreatorID: 12345,
			Players:   []PlayerResponse{},
		}

		assert.Empty(t, response.Players)
	})
}

func TestPlayerResponse(t *testing.T) {
	player := PlayerResponse{
		ID:           "player-1",
		Name:         "Василий",
		AssignedRole: "townsfolk",
		RealRole:     "demon",
	}

	assert.Equal(t, "player-1", player.ID)
	assert.Equal(t, "Василий", player.Name)
	assert.Equal(t, "townsfolk", player.AssignedRole)
	assert.Equal(t, "demon", player.RealRole)

	t.Run("player without real role", func(t *testing.T) {
		player := PlayerResponse{
			ID:           "player-2",
			Name:         "Мария",
			AssignedRole: "minion",
			RealRole:     "",
		}

		assert.Empty(t, player.RealRole)
	})
}

func TestPredictionResponse(t *testing.T) {
	t.Run("prediction without awarding points", func(t *testing.T) {
		response := PredictionResponse{
			ID:            "prediction-123",
			GameID:        "game-123",
			UserID:        "user-456",
			PlayerSlotID:  "slot-789",
			PredictedRole: "demon",
			Points:        0,
			PointsAwarded: false,
			CreatedAt:     time.Now(),
		}

		assert.Equal(t, "prediction-123", response.ID)
		assert.Equal(t, "demon", response.PredictedRole)
		assert.Equal(t, 0, response.Points)
		assert.False(t, response.PointsAwarded)
		assert.NotZero(t, response.CreatedAt)
	})

	t.Run("prediction with awarded points", func(t *testing.T) {
		response := PredictionResponse{
			ID:            "prediction-123",
			GameID:        "game-123",
			UserID:        "user-456",
			PlayerSlotID:  "slot-789",
			PredictedRole: "demon",
			Points:        10,
			PointsAwarded: true,
			CreatedAt:     time.Now(),
		}

		assert.Equal(t, 10, response.Points)
		assert.True(t, response.PointsAwarded)
	})
}

func TestUserResponse(t *testing.T) {
	response := UserResponse{
		ID:         "user-123",
		TelegramID: 12345,
		Username:   "test_user",
		Balance:    100,
		CreatedAt:  time.Now(),
	}

	assert.Equal(t, "user-123", response.ID)
	assert.Equal(t, int64(12345), response.TelegramID)
	assert.Equal(t, "test_user", response.Username)
	assert.Equal(t, 100, response.Balance)
	assert.NotZero(t, response.CreatedAt)
}

func TestCalculateResultsResponse(t *testing.T) {
	response := CalculateResultsResponse{
		GameID: "game-123",
		Scores: map[string]int{
			"user-1": 15,
			"user-2": 10,
			"user-3": -5,
		},
	}

	assert.Equal(t, "game-123", response.GameID)
	assert.Len(t, response.Scores, 3)
	assert.Equal(t, 15, response.Scores["user-1"])
	assert.Equal(t, -5, response.Scores["user-3"])
}

// ==================== TESTS OF METHODS DTO ====================

func TestCreateGameCommand_TrimName(t *testing.T) {
	cmd := CreateGameCommand{
		Name:      "  test game  ",
		CreatorID: 12345,
	}

	trimmed := cmd.TrimmedName()
	assert.Equal(t, "test game", trimmed)
}

func TestSubmitPredictionCommand_IsValidRole(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"demon", true},
		{"minion", true},
		{"townsfolk", true},
		{"outsider", true},
		{"traveler", false},
		{"invalid", false},
		{"", false},
		{"DEMON", false}, // case matters
	}

	cmd := SubmitPredictionCommand{}
	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			cmd.PredictedRole = tt.role
			isValid := cmd.IsValidRole()
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

func TestGameResponse_HasRealRolesSet(t *testing.T) {
	t.Run("all real roles are set", func(t *testing.T) {
		response := GameResponse{
			Players: []PlayerResponse{
				{RealRole: "demon"},
				{RealRole: "minion"},
				{RealRole: "townsfolk"},
			},
		}

		assert.True(t, response.HasRealRolesSet())
	})

	t.Run("not all real roles are set", func(t *testing.T) {
		response := GameResponse{
			Players: []PlayerResponse{
				{RealRole: "demon"},
				{RealRole: ""}, // empty role
				{RealRole: "townsfolk"},
			},
		}

		assert.False(t, response.HasRealRolesSet())
	})

	t.Run("no players", func(t *testing.T) {
		response := GameResponse{
			Players: []PlayerResponse{},
		}

		assert.True(t, response.HasRealRolesSet())
	})
}

func TestPredictionResponse_IsCorrect(t *testing.T) {
	t.Run("prediction is correct", func(t *testing.T) {
		response := PredictionResponse{
			PredictedRole: "demon",
		}

		assert.True(t, response.IsCorrect("demon"))
	})

	t.Run("prediction is not correct", func(t *testing.T) {
		response := PredictionResponse{
			PredictedRole: "demon",
		}

		assert.False(t, response.IsCorrect("minion"))
	})

	t.Run("empty real role", func(t *testing.T) {
		response := PredictionResponse{
			PredictedRole: "demon",
		}

		assert.False(t, response.IsCorrect(""))
	})
}

// ==================== TESTS FOR CONVERSIONS ====================

func TestFromDomainGame(t *testing.T) {
	t.Run("conversion domain game into DTO", func(t *testing.T) {
		// create domain game
		game := entities.NewGame("test game", 12345)

		// add players
		err := game.AddPlayer(111, "Player 1", "townsfolk")
		assert.NoError(t, err)
		err = game.AddPlayer(222, "Player 2", "outsider")
		assert.NoError(t, err)

		// convert
		response := FromDomainGame(game)

		assert.Equal(t, string(game.ID()), response.ID)
		assert.Equal(t, game.Name(), response.Name)
		assert.Equal(t, game.Status(), response.Status)
		assert.Equal(t, game.CreatorID(), response.CreatorID)
		assert.Len(t, response.Players, 2)
		assert.Equal(t, "Player 1", response.Players[0].Name)
		assert.Equal(t, "townsfolk", response.Players[0].AssignedRole)
	})
}

func TestFromDomainPrediction(t *testing.T) {
	t.Run("converting domain prediction into DTO", func(t *testing.T) {
		prediction := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)

		response := FromDomainPrediction(prediction)

		assert.Equal(t, string(prediction.ID()), response.ID)
		assert.Equal(t, string(prediction.GameID()), response.GameID)
		assert.Equal(t, string(prediction.UserID()), response.UserID)
		assert.Equal(t, string(prediction.PlayerSlotID()), response.PlayerSlotID)
		assert.Equal(t, prediction.PredictedRole(), response.PredictedRole)

		// Check that points have not been awarded
		points, awarded := prediction.PointsAwarded()
		assert.Equal(t, points, response.Points)
		assert.Equal(t, awarded, response.PointsAwarded)
	})

	t.Run("prediction with awarded points", func(t *testing.T) {
		prediction := entities.NewPrediction(
			"game-123",
			"user-456",
			"slot-789",
			"demon",
		)

		err := prediction.AwardPoints(10)
		assert.NoError(t, err)

		response := FromDomainPrediction(prediction)

		assert.Equal(t, 10, response.Points)
		assert.True(t, response.PointsAwarded)
	})
}

// ==================== ROLE VALIDATION TESTS ====================

func TestValidRoles(t *testing.T) {
	validRoles := GetValidRoles()

	assert.Contains(t, validRoles, "demon")
	assert.Contains(t, validRoles, "minion")
	assert.Contains(t, validRoles, "townsfolk")
	assert.Contains(t, validRoles, "outsider")

	assert.NotContains(t, validRoles, "traveler")
	assert.NotContains(t, validRoles, "invalid")
	assert.NotContains(t, validRoles, "")
}

func TestIsValidRole(t *testing.T) {
	assert.True(t, IsValidRole("demon"))
	assert.True(t, IsValidRole("minion"))
	assert.True(t, IsValidRole("townsfolk"))
	assert.False(t, IsValidRole("invalid"))
	assert.False(t, IsValidRole(""))
	assert.False(t, IsValidRole("DEMON")) // case matters
}
