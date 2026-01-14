package handlers

import (
	"errors"
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/application/usecases"
	"RIP-Peroni/blood_guess/internal/domain/constants"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/formatting"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SetRealRoleHandler processes the /setrealrole command
type SetRealRoleHandler struct {
	*BaseHandler
	setRealRoleInput ports.SetRealRoleInput
	gameFinder       usecases.GameFinderInterface
}

func NewSetRealRoleHandler(
	bot BotClient,
	setRealRoleInput ports.SetRealRoleInput,
	gameFinder usecases.GameFinderInterface,
) *SetRealRoleHandler {
	return &SetRealRoleHandler{
		BaseHandler:      NewBaseHandler(bot),
		setRealRoleInput: setRealRoleInput,
		gameFinder:       gameFinder,
	}
}

func (h *SetRealRoleHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "setrealrole")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := fmt.Sprintf(`<b>Использование:</b> <code>/setrealrole &lt;имя игрока&gt; &lt;роль&gt; [&lt;имя игрока&gt; &lt;роль&gt; ...]</code>

<b>Пример:</b> <code>/setrealrole Вася demon Коля minion</code>

<b>❗ Внимание:</b> Эта команда устанавливает РЕАЛЬНЫЕ роли игроков после окончания игры.
Используйте только после того, как реальная игра завершена.

<b>Доступные реальные роли:</b>
• <code>demon</code> - демон%s
• <code>minion</code> - приспешник%s
• <code>townsfolk</code> - горожанин%s
• <code>outsider</code> - изгой%s

<b>Примечание:</b> Только создатель игры может устанавливать реальные роли.`,
			formatting.RoleEmoji("demon"),
			formatting.RoleEmoji("minion"),
			formatting.RoleEmoji("townsfolk"),
			formatting.RoleEmoji("outsider"))
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Находим последнюю игру в статусе FINISHED
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusFinished)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для установки реальных ролей!</b>

Нет игр в статусе "завершена". Возможные причины:
1. Игра еще не завершена - используйте <code>/finish</code>
2. Игра еще не начата - используйте <code>/startgame</code>
3. Используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может устанавливать реальные роли.")
	}

	// Парсим аргументы: чередование имени игрока и роли
	parts := parseArguments(args)
	if len(parts)%2 != 0 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Нечетное количество аргументов. Ожидается формат: <code>/setrealrole имя1 роль1 имя2 роль2 ...</code>")
	}

	var successMessages []string
	var errorMessages []string

	// Обрабатываем пары (имя, роль)
	for i := 0; i < len(parts); i += 2 {
		playerName := parts[i]
		role := parts[i+1]

		if !dto.IsValidGameRole(role) {
			errorMessages = append(errorMessages,
				fmt.Sprintf("❌ Недопустимая роль для игрока '%s': %s", playerName, role))
			continue
		}

		// Находим игрока по имени
		var playerSlotID string
		for _, player := range game.Players() {
			if player.Name == playerName {
				playerSlotID = string(player.ID)
				break
			}
		}

		if playerSlotID == "" {
			errorMessages = append(errorMessages,
				fmt.Sprintf("❌ Игрок '%s' не найден в игре", playerName))
			continue
		}

		command := dto.SetRealRoleCommand{
			GameID:       string(game.ID()),
			PlayerSlotID: playerSlotID,
			RealRole:     role,
			AdminID:      update.Message.From.ID,
		}

		if err := h.setRealRoleInput.Execute(command); err != nil {
			errorMessages = append(errorMessages,
				fmt.Sprintf("❌ Не удалось установить роль для '%s': %v", playerName, err))
		} else {
			roleEmoji := formatting.RoleEmoji(role)
			roleDisplayName := formatting.RoleDisplayName(role)

			pointsInfo := ""
			if formatting.IsEvilRole(role) {
				pointsInfo = fmt.Sprintf(" (%+d очков за правильный прогноз)",
					constants.PointsForRole(role))
			}

			successMessages = append(successMessages,
				fmt.Sprintf("✅ %s → %s %s%s",
					playerName, roleEmoji, roleDisplayName, pointsInfo))
		}
	}

	// Формируем итоговое сообщение
	var sb strings.Builder

	if len(successMessages) > 0 {
		sb.WriteString(fmt.Sprintf(`✅ <b>Реальные роли установлены!</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус:</b> %s
<b>✅ Установленные роли:</b>
%s
`,
			h.EscapeHTML(game.Name()),
			game.Status(),
			strings.Join(successMessages, "\n")))
	}

	if len(errorMessages) > 0 {
		sb.WriteString("\n<b>⚠️ Ошибки:</b>\n")
		sb.WriteString(strings.Join(errorMessages, "\n"))
	}

	if len(successMessages) == 0 && len(errorMessages) == 0 {
		sb.WriteString("❌ Не удалось установить ни одной роли.")
	}

	return h.SendHTML(update.Message.Chat.ID, sb.String())
}

func (h *SetRealRoleHandler) Command() string {
	return "setrealrole"
}

func (h *SetRealRoleHandler) Description() string {
	return "Установить реальные роли игрокам после игры"
}
