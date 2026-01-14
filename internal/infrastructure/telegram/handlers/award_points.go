package handlers

import (
	"errors"
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/application/usecases"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AwardPointsHandler handles /awardpoints command
type AwardPointsHandler struct {
	*BaseHandler
	awardPointsInput ports.AwardPointsInput
	gameFinder       usecases.GameFinderInterface
}

func NewAwardPointsHandler(
	bot BotClient,
	awardPointsInput ports.AwardPointsInput,
	gameFinder usecases.GameFinderInterface,
) *AwardPointsHandler {
	return &AwardPointsHandler{
		BaseHandler:      NewBaseHandler(bot),
		awardPointsInput: awardPointsInput,
		gameFinder:       gameFinder,
	}
}

func (h *AwardPointsHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "awardpoints")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args != "" {
		// Старый формат с ID игры - для обратной совместимости
		_ = h.SendHTML(update.Message.Chat.ID,
			`<i>Примечание: Теперь команда /awardpoints не требует ID игры. Она автоматически находит последнюю игру в статусе FINISHED.</i>`)
	}

	// Находим последнюю игру в статусе FINISHED
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusFinished)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для начисления очков!</b>

Нет игр в статусе "завершена". Возможные причины:
1. Игра еще не завершена - используйте <code>/finish</code>
2. Используйте <code>/games</code> для просмотра игр`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может начислять очки.")
	}

	command := dto.AwardPointsCommand{
		GameID:  string(game.ID()),
		AdminID: update.Message.From.ID,
	}

	response, err := h.awardPointsInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось начислить очки: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	// Build results message
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`🎉 <b>Очки начислены!</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>💰 Валюта начислена:</b> %v
<b>⏰ Начислено:</b> %s

<b>🏆 Результаты по пользователям:</b>
`,
		h.EscapeHTML(response.Name),
		response.Status,
		response.CurrencyAwarded,
		response.EndedAt.Format("02.01.2006 15:04")))

	if len(response.UserScores) > 0 {
		for userID, score := range response.UserScores {
			sb.WriteString(fmt.Sprintf("  Пользователь %s: %+d очков\n",
				h.EscapeHTML(userID), score))
		}
	} else {
		sb.WriteString("  Нет прогнозов для подсчета\n")
	}

	sb.WriteString("\nОчки успешно начислены! 🎭")

	return h.SendHTML(update.Message.Chat.ID, sb.String())
}

func (h *AwardPointsHandler) Command() string {
	return "awardpoints"
}

func (h *AwardPointsHandler) Description() string {
	return "Начислить очки за прогнозы (после установки всех ролей)"
}
