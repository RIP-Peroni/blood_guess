package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewGameHandler handles /newgame command
type NewGameHandler struct {
	*BaseHandler
	createGameInput ports.CreateGameInput
}

func NewNewGameHandler(bot BotClient, createGameInput ports.CreateGameInput) *NewGameHandler {
	return &NewGameHandler{
		BaseHandler:     NewBaseHandler(bot),
		createGameInput: createGameInput,
	}
}

func (h *NewGameHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "newgame")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/newgame &lt;название игры&gt;</code>

<b>Пример:</b> <code>/newgame Игра от Ивана 15.02.2024</code>`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	command := dto.CreateGameCommand{
		Name:      args,
		CreatorID: update.Message.From.ID,
	}

	response, err := h.createGameInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось создать игру: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	escapedName := h.EscapeHTML(response.Name)
	escapedID := h.EscapeHTML(response.ID)

	successMsg := fmt.Sprintf(`<b>🎮 Игра создана!</b>

<b>📛 Название:</b> %s
<b>🆔 ID игры:</b> <code>%s</code>
<b>📊 Статус:</b> %s
<b>👑 Создатель:</b> вы
<b>👥 Игроков:</b> %d

Теперь добавьте нескольких игроков сразу при помощи <code>/addplayers</code>
или по одному игроку за раз при помощи  <code>/addplayer</code>`,
		escapedName, escapedID, response.Status, len(response.Players))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *NewGameHandler) Command() string {
	return "newgame"
}

func (h *NewGameHandler) Description() string {
	return "Создать новую игру"
}
