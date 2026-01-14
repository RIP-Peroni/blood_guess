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

// FinishGameHandler handles /finish command
type FinishGameHandler struct {
	*BaseHandler
	finishGameInput ports.FinishGameInput
	gameFinder      usecases.GameFinderInterface
}

func NewFinishGameHandler(
	bot BotClient,
	finishGameInput ports.FinishGameInput,
	gameFinder usecases.GameFinderInterface,
) *FinishGameHandler {
	return &FinishGameHandler{
		BaseHandler:     NewBaseHandler(bot),
		finishGameInput: finishGameInput,
		gameFinder:      gameFinder,
	}
}

func (h *FinishGameHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "finish")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args != "" {
		// Старый формат с ID игры - для обратной совместимости
		_ = h.SendHTML(update.Message.Chat.ID,
			`<i>Примечание: Теперь команда /finish не требует ID игры. Она автоматически находит последнюю игру в процессе.</i>`)
	}

	// Находим последнюю игру в статусе IN_PROGRESS
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusInProgress)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для завершения!</b>

Нет игр в статусе "в процессе". Возможные причины:
1. Игра еще не начата - используйте <code>/startgame</code>
2. Игра уже завершена - используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может завершать игру.")
	}

	// Immediately finish the game without confirmation
	command := dto.FinishGameCommand{
		GameID:  string(game.ID()),
		AdminID: update.Message.From.ID,
	}

	response, err := h.finishGameInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось завершить игру: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	successMsg := fmt.Sprintf(`✅ <b>Игра завершена!</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>⏰ Завершена:</b> %s

Теперь установите реальные <b>злые роли</b> игроков с помощью <code>/setrealrole</code>

<b>Пример:</b> <code>/setrealrole Вася demon Коля minion</code>

<b>Примечание:</b> Устанавливайте только злые роли (демон и приспешников). Все остальные игроки автоматически считаются мирными.`,
		h.EscapeHTML(response.Name),
		response.Status,
		response.EndedAt.Format("02.01.2006 15:04"))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *FinishGameHandler) Command() string {
	return "finish"
}

func (h *FinishGameHandler) Description() string {
	return "Завершить последнюю игру в процессе"
}
