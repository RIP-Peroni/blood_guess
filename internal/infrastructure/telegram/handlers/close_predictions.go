package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ClosePredictionsHandler handles /closepred command
type ClosePredictionsHandler struct {
	*BaseHandler
	closePredictionsInput ports.ClosePredictionsInput
	gameRepo              ports.GameRepository
}

func NewClosePredictionsHandler(
	bot BotClient,
	closePredictionsInput ports.ClosePredictionsInput,
	gameRepo ports.GameRepository,
) *ClosePredictionsHandler {
	return &ClosePredictionsHandler{
		BaseHandler:           NewBaseHandler(bot),
		closePredictionsInput: closePredictionsInput,
		gameRepo:              gameRepo,
	}
}

func (h *ClosePredictionsHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "closepred")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/closepred &lt;ID_игры&gt;</code>

<b>Пример:</b> <code>/closepred abc123</code>

<b>Примечание:</b> Только создатель игры может закрывать прогнозы.`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	gameID := args

	game, err := h.gameRepo.FindByID(entities.GameID(gameID))
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Игра с ID <code>%s</code> не найдена.", h.EscapeHTML(gameID)))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может закрывать прогнозы.")
	}

	if game.Status() != entities.GameStatusPredictionsOpen {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Нельзя закрыть прогнозы. Текущий статус игры: <b>%s</b>. Прогнозы должны быть открыты.", game.Status()))
	}

	command := dto.ClosePredictionsCommand{
		GameID:  gameID,
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

Для начала игры используйте: <code>/startgame %s</code>`,
		h.EscapeHTML(response.Name),
		response.Status,
		len(response.Players),
		h.EscapeHTML(gameID))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *ClosePredictionsHandler) Command() string {
	return "closepred"
}

func (h *ClosePredictionsHandler) Description() string {
	return "Закрыть прогнозы для игры"
}
