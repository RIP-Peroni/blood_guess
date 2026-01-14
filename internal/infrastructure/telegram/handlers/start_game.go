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

// StartGameHandler handles /startgame command
type StartGameHandler struct {
	*BaseHandler
	startGameInput ports.StartGameInput
	gameFinder     *usecases.GameFinder
}

func NewStartGameHandler(
	bot BotClient,
	startGameInput ports.StartGameInput,
	gameFinder *usecases.GameFinder,
) *StartGameHandler {
	return &StartGameHandler{
		BaseHandler:    NewBaseHandler(bot),
		startGameInput: startGameInput,
		gameFinder:     gameFinder,
	}
}

func (h *StartGameHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "startgame")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args != "" {
		// Старый формат с ID игры - для обратной совместимости
		h.SendHTML(update.Message.Chat.ID,
			`<i>Примечание: Теперь команда /startgame не требует ID игры. Она автоматически находит последнюю игру с закрытыми прогнозами.</i>`)
	}

	// Находим последнюю игру в статусе PREDICTIONS_CLOSED
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusPredictionsClosed)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для начала!</b>

Нет игр в статусе "прогнозы закрыты". Возможные причины:
1. Игра еще не создана - используйте <code>/newgame</code>
2. Прогнозы еще не закрыты - используйте <code>/closepred</code>
3. Игра уже перешла в другой статус - используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может начинать реальную игру.")
	}

	command := dto.StartGameCommand{
		GameID:  string(game.ID()),
		AdminID: update.Message.From.ID,
	}

	response, err := h.startGameInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось начать игру: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	successMsg := fmt.Sprintf(`🎲 <b>Реальная игра началась!</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>👥 Игроков:</b> %d

Теперь можно играть в реальной игре! После окончания игры установите реальные роли игроков с помощью <code>/setrealrole</code>

<b>Пример:</b> <code>/setrealrole Вася demon Коля minion</code>

Когда все роли установлены, завершите игру командой <code>/finish</code> для подсчета очков.`,
		h.EscapeHTML(response.Name),
		response.Status,
		len(response.Players))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *StartGameHandler) Command() string {
	return "startgame"
}

func (h *StartGameHandler) Description() string {
	return "Начать реальную игру (после закрытия прогнозов)"
}
