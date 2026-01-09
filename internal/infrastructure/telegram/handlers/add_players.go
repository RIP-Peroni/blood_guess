package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AddPlayersHandler processes the /addplayers command
type AddPlayersHandler struct {
	*BaseHandler
	addPlayersInput ports.AddPlayersInput
	gameRepo        ports.GameRepository
}

func NewAddPlayersHandler(
	bot BotClient,
	addPlayersInput ports.AddPlayersInput,
	gameRepo ports.GameRepository,
) *AddPlayersHandler {
	return &AddPlayersHandler{
		BaseHandler:     NewBaseHandler(bot),
		addPlayersInput: addPlayersInput,
		gameRepo:        gameRepo,
	}
}

func (h *AddPlayersHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "addplayers")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/addplayers &lt;ID_игры&gt; &lt;имя1&gt; &lt;имя2&gt; ...</code>

<b>Пример:</b> <code>/addplayers abc123 Вася Петя Коля Миша Алекс</code>

<b>Примечание:</b>
• Имена игроков должны быть уникальными в рамках игры
• Можно использовать имена с пробелами, заключив их в кавычки: <code>/addplayers abc123 "Вася Пупкин" Петя</code>
• Роли игрокам не назначаются. Их можно назначить позже командой <code>/assignroles</code>`

		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Parsing arguments with quote support
	parts := parseArguments(args)
	if len(parts) < 2 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Недостаточно аргументов. Используйте: <code>/addplayers &lt;ID_игры&gt; &lt;имя1&gt; &lt;имя2&gt; ...</code>")
	}

	gameID := parts[0]
	playerNames := parts[1:]

	game, err := h.gameRepo.FindByID(entities.GameID(gameID))
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Игра с ID <code>%s</code> не найдена.", h.EscapeHTML(gameID)))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может добавлять игроков.")
	}

	command := dto.AddPlayersCommand{
		GameID:      gameID,
		PlayerNames: playerNames,
		AdminID:     update.Message.From.ID,
	}

	response, err := h.addPlayersInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось добавить игроков: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	successMsg := fmt.Sprintf(`✅ <b>Игроки добавлены!</b>

<b>🎮 Игра:</b> %s
<b>👥 Добавлено игроков:</b> %d
<b>📋 Список:</b> %s
<b>👤 Всего игроков в игре:</b> %d

Теперь можно добавить ещё игроков или открыть прогнозы с помощью <code>/openpred %s</code>

<b>Примечание:</b> Роли игрокам не назначены. Вы можете назначить их командой <code>/assignroles %s</code>`,
		h.EscapeHTML(response.Name),
		len(playerNames),
		strings.Join(playerNames, ", "),
		len(response.Players),
		h.EscapeHTML(gameID),
		h.EscapeHTML(gameID))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *AddPlayersHandler) Command() string {
	return "addplayers"
}

func (h *AddPlayersHandler) Description() string {
	return "Добавить нескольких игроков в игру (без указания ролей)"
}

// parseArguments парсит аргументы с поддержкой кавычек
func parseArguments(input string) []string {
	var result []string
	var current strings.Builder
	inQuotes := false
	escapeNext := false

	for _, r := range input {
		if escapeNext {
			current.WriteRune(r)
			escapeNext = false
			continue
		}

		switch r {
		case '\\':
			escapeNext = true
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				current.WriteRune(r)
			} else {
				if current.Len() > 0 {
					result = append(result, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}
