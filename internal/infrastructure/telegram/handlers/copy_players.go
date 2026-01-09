package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// CopyPlayersHandler processes the /copyplayers command
type CopyPlayersHandler struct {
	*BaseHandler
	copyPlayersInput ports.CopyPlayersInput
	gameRepo         ports.GameRepository
}

func NewCopyPlayersHandler(
	bot BotClient,
	copyPlayersInput ports.CopyPlayersInput,
	gameRepo ports.GameRepository,
) *CopyPlayersHandler {
	return &CopyPlayersHandler{
		BaseHandler:      NewBaseHandler(bot),
		copyPlayersInput: copyPlayersInput,
		gameRepo:         gameRepo,
	}
}

func (h *CopyPlayersHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "copyplayers")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/copyplayers &lt;ID_новой_игры&gt;</code>

<b>Пример:</b> <code>/copyplayers abc123</code>

<b>Что делает:</b>
• Копирует имена игроков из вашей последней игры
• Не копирует роли игроков
• Игроки добавляются без назначенных ролей

<b>Примечание:</b> Только создатель игры может копировать игроков.`

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
			"❌ Только создатель игры может копировать игроков.")
	}

	command := dto.CopyPlayersCommand{
		TargetGameID: gameID,
		AdminID:      update.Message.From.ID,
	}

	response, err := h.copyPlayersInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось скопировать игроков: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	// Получаем имена скопированных игроков
	playerNames := make([]string, len(response.Players))
	for i, player := range response.Players {
		playerNames[i] = player.Name
	}

	successMsg := fmt.Sprintf(`✅ <b>Игроки скопированы из последней игры!</b>

<b>🎮 Игра:</b> %s
<b>👥 Скопировано игроков:</b> %d
<b>📋 Список:</b> %s

Теперь можно добавить ещё игроков или открыть прогнозы с помощью <code>/openpred %s</code>

<b>Примечание:</b> Роли игрокам не назначены. Вы можете назначить их командой <code>/assignroles %s</code>`,
		h.EscapeHTML(response.Name),
		len(playerNames),
		strings.Join(playerNames, ", "),
		h.EscapeHTML(gameID),
		h.EscapeHTML(gameID))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *CopyPlayersHandler) Command() string {
	return "copyplayers"
}

func (h *CopyPlayersHandler) Description() string {
	return "Скопировать игроков из последней игры (без ролей)"
}
