package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// OpenPredictionsHandler handles /openpred command
type OpenPredictionsHandler struct {
	*BaseHandler
	openPredictionsInput ports.OpenPredictionsInput
	gameRepo             ports.GameRepository
}

func NewOpenPredictionsHandler(
	bot BotClient,
	openPredictionsInput ports.OpenPredictionsInput,
	gameRepo ports.GameRepository,
) *OpenPredictionsHandler {
	return &OpenPredictionsHandler{
		BaseHandler:          NewBaseHandler(bot),
		openPredictionsInput: openPredictionsInput,
		gameRepo:             gameRepo,
	}
}

func (h *OpenPredictionsHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "openpred")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/openpred &lt;ID_игры&gt;</code>

<b>Пример:</b> <code>/openpred abc123</code>

<b>Примечание:</b> Только создатель игры может открывать прогнозы.`
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
			"❌ Только создатель игры может открывать прогнозы.")
	}

	if game.Status() != entities.GameStatusCreated {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Нельзя открыть прогнозы. Текущий статус игры: <b>%s</b>.", game.Status()))
	}

	if len(game.Players()) == 0 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Нельзя открыть прогнозы для игры без игроков. Сначала добавьте игроков с помощью <code>/addplayer</code>.")
	}

	command := dto.OpenPredictionsCommand{
		GameID:  gameID,
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

Для закрытия прогнозов используйте: <code>/closepred %s</code>`,
		h.EscapeHTML(response.Name),
		response.Status,
		len(response.Players),
		h.EscapeHTML(gameID))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *OpenPredictionsHandler) Command() string {
	return "openpred"
}

func (h *OpenPredictionsHandler) Description() string {
	return "Открыть прогнозы для игры"
}
