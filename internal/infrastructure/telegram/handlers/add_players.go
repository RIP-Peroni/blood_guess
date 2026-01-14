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

// AddPlayersHandler processes the /addplayers command
type AddPlayersHandler struct {
	*BaseHandler
	addPlayersInput ports.AddPlayersInput
	gameFinder      usecases.GameFinderInterface
}

func NewAddPlayersHandler(
	bot BotClient,
	addPlayersInput ports.AddPlayersInput,
	gameFinder usecases.GameFinderInterface,
) *AddPlayersHandler {
	return &AddPlayersHandler{
		BaseHandler:     NewBaseHandler(bot),
		addPlayersInput: addPlayersInput,
		gameFinder:      gameFinder,
	}
}

func (h *AddPlayersHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "addplayers")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/addplayers &lt;имя1&gt; &lt;имя2&gt; ...</code>

<b>Пример:</b> <code>/addplayers Вася Петя Коля Миша Алекс</code>

<b>Примечание:</b>
• Игроки будут добавлены в последнюю созданную игру (статус: CREATED)
• Имена игроков должны быть уникальными в рамках игры
• Можно использовать имена с пробелами, заключив их в кавычки: <code>/addplayers "Вася Пупкин" Петя</code>
• Роли игрокам не назначаются. Их можно назначить позже командой <code>/setrealrole</code>`

		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Находим последнюю игру в статусе CREATED
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusCreated)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для добавления игроков!</b>

Нет игр в статусе "создана". Возможные причины:
1. Игра еще не создана - используйте <code>/newgame</code>
2. Игра уже перешла в другой статус - используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может добавлять игроков.")
	}

	// Parsing arguments with quote support
	playerNames := parseArguments(args)

	command := dto.AddPlayersCommand{
		GameID:      string(game.ID()),
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

Теперь можно добавить ещё игроков или открыть прогнозы с помощью <code>/openpred</code>`,
		h.EscapeHTML(response.Name),
		len(playerNames),
		strings.Join(playerNames, ", "),
		len(response.Players))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *AddPlayersHandler) Command() string {
	return "addplayers"
}

func (h *AddPlayersHandler) Description() string {
	return "Добавить несколько игроков в последнюю созданную игру"
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
