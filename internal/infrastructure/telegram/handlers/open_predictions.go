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

// OpenPredictionsHandler handles /openpred command
type OpenPredictionsHandler struct {
	*BaseHandler
	openPredictionsInput ports.OpenPredictionsInput
	gameFinder           *usecases.GameFinder
}

func NewOpenPredictionsHandler(
	bot BotClient,
	openPredictionsInput ports.OpenPredictionsInput,
	gameFinder *usecases.GameFinder,
) *OpenPredictionsHandler {
	return &OpenPredictionsHandler{
		BaseHandler:          NewBaseHandler(bot),
		openPredictionsInput: openPredictionsInput,
		gameFinder:           gameFinder,
	}
}

func (h *OpenPredictionsHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "openpred")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args != "" {
		// Старый формат с ID игры - для обратной совместимости
		_ = h.SendHTML(update.Message.Chat.ID,
			`<i>Примечание: Теперь команда /openpred не требует ID игры. Она автоматически находит последнюю созданную игру.</i>`)
	}

	// Находим последнюю игру в статусе CREATED
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusCreated)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для открытия прогнозов!</b>

Нет игр в статусе "создана". Возможные причины:
1. Игра еще не создана - используйте <code>/newgame</code>
2. Игроки еще не добавлены - используйте <code>/addplayers</code>
3. Игра уже перешла в другой статус - используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может открывать прогнозы.")
	}

	if len(game.Players()) == 0 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Нельзя открыть прогнозы для игры без игроков. Сначала добавьте игроков с помощью <code>/addplayers</code>.")
	}

	command := dto.OpenPredictionsCommand{
		GameID:  string(game.ID()),
		AdminID: update.Message.From.ID,
	}

	response, err := h.openPredictionsInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось открыть прогнозы: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	successMsg := fmt.Sprintf(`✅ <b>Прогнозы открыты!</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>👥 Игроков:</b> %d

Теперь участники могут делать прогнозы с помощью команды <code>/predict</code>.

Для закрытия прогнозов используйте: <code>/closepred</code>`,
		h.EscapeHTML(response.Name),
		response.Status,
		len(response.Players))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *OpenPredictionsHandler) Command() string {
	return "openpred"
}

func (h *OpenPredictionsHandler) Description() string {
	return "Открыть прогнозы для последней созданной игры"
}
