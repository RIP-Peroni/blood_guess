package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// FinishGameHandler handles /finish command
type FinishGameHandler struct {
	*BaseHandler
	finishGameInput ports.FinishGameInput
	gameRepo        ports.GameRepository
}

func NewFinishGameHandler(
	bot BotClient,
	finishGameInput ports.FinishGameInput,
	gameRepo ports.GameRepository,
) *FinishGameHandler {
	return &FinishGameHandler{
		BaseHandler:     NewBaseHandler(bot),
		finishGameInput: finishGameInput,
		gameRepo:        gameRepo,
	}
}

func (h *FinishGameHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "finish")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/finish &lt;ID_игры&gt;</code>

<b>Пример:</b> <code>/finish abc123</code>

<b>Требования:</b>
• Игра должна быть в процессе (реальная игра начата)
• Для всех игроков должны быть установлены реальные роли (командой <code>/setrole</code>)
• Только создатель игры может завершать игру

<b>Результат:</b>
• Очки будут подсчитаны по прогнозам
• Валюта будет начислена пользователям
• Игра будет завершена`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	gameID := args

	// First, check if game exists and user is creator
	game, err := h.gameRepo.FindByID(entities.GameID(gameID))
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Игра с ID <code>%s</code> не найдена.", h.EscapeHTML(gameID)))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может завершать игру.")
	}

	// Check current status for better error messages
	if game.Status() != entities.GameStatusInProgress {
		statusMessages := map[entities.GameStatus]string{
			entities.GameStatusCreated:           "❌ Игра еще не начата. Сначала откройте прогнозы командой <code>/openpred</code>.",
			entities.GameStatusPredictionsOpen:   "❌ Прогнозы еще открыты. Сначала закройте их командой <code>/closepred</code>.",
			entities.GameStatusPredictionsClosed: "❌ Игра еще не начата. Начните реальную игру командой <code>/startgame</code>.",
			entities.GameStatusFinished:          "❌ Игра уже завершена!",
		}

		if msg, ok := statusMessages[game.Status()]; ok {
			return h.SendHTML(update.Message.Chat.ID, msg)
		}
	}

	// Check if all real roles are set
	missingRoles := []string{}
	for _, player := range game.Players() {
		if !player.IsRealRoleSet {
			missingRoles = append(missingRoles, player.Name)
		}
	}

	if len(missingRoles) > 0 {
		message := fmt.Sprintf(`❌ <b>Не все реальные роли установлены!</b>

Следующие игроки не имеют реальной роли:
%s

Используйте команду <code>/setrole %s &lt;ID_игрока&gt; &lt;роль&gt;</code> для установки ролей.

<b>Пример:</b> <code>/setrole %s %s demon</code>`,
			strings.Join(missingRoles, "\n• "),
			h.EscapeHTML(gameID),
			h.EscapeHTML(gameID),
			h.EscapeHTML(string(game.Players()[0].ID)))

		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Ask for confirmation
	confirmationMsg := fmt.Sprintf(`⚠️ <b>Подтверждение завершения игры</b>

<b>🎮 Игра:</b> %s
<b>👥 Игроков:</b> %d
<b>📊 Прогнозов сделано:</b> ?

Вы уверены, что хотите завершить игру и подсчитать результаты?

После завершения:
• Очки будут подсчитаны и начислены
• Прогнозы больше нельзя будет изменить
• Игра будет завершена

Для подтверждения отправьте: <code>/finish %s confirm</code>`,
		h.EscapeHTML(game.Name()),
		len(game.Players()),
		h.EscapeHTML(gameID))

	// If confirmation provided
	if strings.HasSuffix(args, " confirm") || strings.Contains(args, " confirm ") {
		// Remove "confirm" from args
		gameID = strings.TrimSpace(strings.Replace(args, "confirm", "", 1))

		command := dto.FinishGameCommand{
			GameID:  gameID,
			AdminID: update.Message.From.ID,
		}

		response, err := h.finishGameInput.Execute(command)
		if err != nil {
			errorMsg := fmt.Sprintf("❌ Не удалось завершить игру: %v", err)
			return h.SendText(update.Message.Chat.ID, errorMsg)
		}

		// Build results message
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf(`🎉 <b>Игра завершена!</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>👥 Игроков:</b> %d
<b>💰 Валюта начислена:</b> %v
<b>⏰ Завершена:</b> %s

<b>📈 Результаты по игрокам:</b>
`,
			h.EscapeHTML(response.Name),
			response.Status,
			len(response.PlayerResults),
			response.CurrencyAwarded,
			response.EndedAt.Format("02.01.2006 15:04")))

		for _, player := range response.PlayerResults {
			sb.WriteString(fmt.Sprintf("\n<b>👤 %s</b>\n", h.EscapeHTML(player.PlayerName)))
			sb.WriteString(fmt.Sprintf("  Назначенная роль: %s\n", player.AssignedRole))
			sb.WriteString(fmt.Sprintf("  Реальная роль: <b>%s</b>\n", player.RealRole))

			if len(player.Predictions) > 0 {
				sb.WriteString("  Прогнозы:\n")
				for _, pred := range player.Predictions {
					icon := "❌"
					if pred.IsCorrect {
						icon = "✅"
					}
					sb.WriteString(fmt.Sprintf("    %s @%s: %s (%+d очков)\n",
						icon,
						h.EscapeHTML(pred.Username),
						pred.PredictedRole,
						pred.Points))
				}
			} else {
				sb.WriteString("  Прогнозов нет\n")
			}
		}

		sb.WriteString("\n<b>🏆 Общие результаты:</b>\n")
		if len(response.UserScores) > 0 {
			for userID, score := range response.UserScores {
				// Try to get username
				// Note: We don't have userRepo here, so we show userID
				sb.WriteString(fmt.Sprintf("  Пользователь %s: %+d очков\n",
					h.EscapeHTML(userID), score))
			}
		} else {
			sb.WriteString("  Нет прогнозов для подсчета\n")
		}

		sb.WriteString("\nИгра завершена! Спасибо за участие! 🎭")

		return h.SendHTML(update.Message.Chat.ID, sb.String())
	}

	// If no confirmation, ask for it
	return h.SendHTML(update.Message.Chat.ID, confirmationMsg)
}

func (h *FinishGameHandler) Command() string {
	return "finish"
}

func (h *FinishGameHandler) Description() string {
	return "Завершить игру и подсчитать результаты"
}
