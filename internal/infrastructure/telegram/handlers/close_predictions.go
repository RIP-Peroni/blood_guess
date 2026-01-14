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

// ClosePredictionsHandler handles /closepred command
type ClosePredictionsHandler struct {
	*BaseHandler
	closePredictionsInput ports.ClosePredictionsInput
	gameFinder            *usecases.GameFinder
}

func NewClosePredictionsHandler(
	bot BotClient,
	closePredictionsInput ports.ClosePredictionsInput,
	gameFinder *usecases.GameFinder,
) *ClosePredictionsHandler {
	return &ClosePredictionsHandler{
		BaseHandler:           NewBaseHandler(bot),
		closePredictionsInput: closePredictionsInput,
		gameFinder:            gameFinder,
	}
}

func (h *ClosePredictionsHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "closepred")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args != "" {
		// Старый формат с ID игры - для обратной совместимости
		h.SendHTML(update.Message.Chat.ID,
			`<i>Примечание: Теперь команда /closepred не требует ID игры. Она автоматически находит последнюю игру с открытыми прогнозами.</i>`)
	}

	// Находим последнюю игру в статусе PREDICTIONS_OPEN
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusPredictionsOpen)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для закрытия прогнозов!</b>

Нет игр в статусе "прогнозы открыты". Возможные причины:
1. Игра еще не создана - используйте <code>/newgame</code>
2. Прогнозы еще не открыты - используйте <code>/openpred</code>
3. Игра уже перешла в другой статус - используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может закрывать прогнозы.")
	}

	command := dto.ClosePredictionsCommand{
		GameID:  string(game.ID()),
		AdminID: update.Message.From.ID,
	}

	response, err := h.closePredictionsInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось закрыть прогнозы: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	successMsg := fmt.Sprintf(`🔒 <b>Прогнозы закрыты!</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>👥 Игроков:</b> %d

Прогнозы больше не принимаются. Можно начинать реальную игру!

Для начала игры используйте: <code>/startgame</code>`,
		h.EscapeHTML(response.Name),
		response.Status,
		len(response.Players))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *ClosePredictionsHandler) Command() string {
	return "closepred"
}

func (h *ClosePredictionsHandler) Description() string {
	return "Закрыть прогнозы для последней игры с открытыми прогнозами"
}
