package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartGameHandler handles /startgame command
type StartGameHandler struct {
	*BaseHandler
	startGameInput ports.StartGameInput
	gameRepo       ports.GameRepository
}

func NewStartGameHandler(
	bot BotClient,
	startGameInput ports.StartGameInput,
	gameRepo ports.GameRepository,
) *StartGameHandler {
	return &StartGameHandler{
		BaseHandler:    NewBaseHandler(bot),
		startGameInput: startGameInput,
		gameRepo:       gameRepo,
	}
}

func (h *StartGameHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "startgame")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/startgame &lt;ID_игры&gt;</code>

<b>Пример:</b> <code>/startgame abc123</code>

<b>Как получить ID игры:</b>
Используйте команду <code>/games</code> для просмотра списка игр.

<b>Примечание:</b> Только создатель игры может начинать реальную игру. Прогнозы должны быть закрыты.`
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
			"❌ Только создатель игры может начинать реальную игру.")
	}

	// Check current status for better error messages
	if game.Status() != entities.GameStatusPredictionsClosed {
		statusMessages := map[entities.GameStatus]string{
			entities.GameStatusCreated:         "❌ Нельзя начать игру. Игра только создана. Сначала откройте прогнозы командой <code>/openpred</code>.",
			entities.GameStatusPredictionsOpen: "❌ Нельзя начать игру. Прогнозы еще открыты. Сначала закройте прогнозы командой <code>/closepred</code>.",
			entities.GameStatusInProgress:      "❌ Игра уже начата!",
			entities.GameStatusFinished:        "❌ Игра уже завершена!",
		}

		if msg, ok := statusMessages[game.Status()]; ok {
			return h.SendHTML(update.Message.Chat.ID, msg)
		}
	}

	command := dto.StartGameCommand{
		GameID:  gameID,
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
<b>🆔 ID игры:</b> <code>%s</code>

Теперь можно играть в реальной игре! После окончания игры установите реальные роли игроков с помощью <code>/setrole %s &lt;ID_игрока&gt; &lt;реальная_роль&gt;</code>

<b>Пример:</b> <code>/setrole %s player-123 demon</code>

Когда все роли установлены, завершите игру командой <code>/finish %s</code> для подсчета очков.`,
		h.EscapeHTML(response.Name),
		response.Status,
		len(response.Players),
		h.EscapeHTML(gameID),
		h.EscapeHTML(gameID),
		h.EscapeHTML(gameID),
		h.EscapeHTML(gameID))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *StartGameHandler) Command() string {
	return "startgame"
}

func (h *StartGameHandler) Description() string {
	return "Начать реальную игру (после закрытия прогнозов)"
}
