package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/application/usecases"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewGameHandler handles /newgame command
type NewGameHandler struct {
	*BaseHandler
	createGameInput ports.CreateGameInput
	gameFinder      usecases.GameFinderInterface
}

func NewNewGameHandler(
	bot BotClient,
	createGameInput ports.CreateGameInput,
	gameFinder usecases.GameFinderInterface,
) *NewGameHandler {
	return &NewGameHandler{
		BaseHandler:     NewBaseHandler(bot),
		createGameInput: createGameInput,
		gameFinder:      gameFinder,
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

<b>Пример:</b> <code>/newgame BMR 15.01.2024</code>

<b>Примечание:</b>
• Можно создать новую игру только если нет других активных игр
• Активная игра - любая игра со статусом, отличным от FINISHED`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Проверяем, можно ли создать новую игру
	canCreate, err := h.gameFinder.CanCreateNewGame()
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при проверке активных игр: %v", err))
	}

	if !canCreate {
		// Пытаемся найти активную игру
		activeGame, err := h.gameFinder.FindActiveGame()
		if err == nil && activeGame != nil {
			return h.SendHTML(update.Message.Chat.ID,
				fmt.Sprintf(`❌ <b>Нельзя создать новую игру!</b>

Уже есть активная игра:
<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>👥 Игроков:</b> %d

Дождитесь завершения текущей игры или завершите её командой <code>/finish</code>`,
					h.EscapeHTML(activeGame.Name()),
					activeGame.Status(),
					len(activeGame.Players())))
		}

		return h.SendHTML(update.Message.Chat.ID,
			"❌ Уже есть активная игра. Дождитесь её завершения.")
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

	successMsg := fmt.Sprintf(`✅ <b>Игра создана!</b>

<b>📛 Название:</b> %s
<b>🆔 ID игры:</b> <code>%s</code>
<b>📊 Статус:</b> %s
<b>👑 Создатель:</b> вы
<b>👥 Игроков:</b> %d

Теперь добавьте игроков с помощью:
• <code>/addplayers Вася Петя Миша Коля</code> - добавить несколько игроков`,
		escapedName, escapedID, response.Status, len(response.Players))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *NewGameHandler) Command() string {
	return "newgame"
}

func (h *NewGameHandler) Description() string {
	return "Создать новую игру (если нет других активных игр)"
}
