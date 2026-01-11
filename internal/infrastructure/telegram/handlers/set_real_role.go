package handlers

import (
	"RIP-Peroni/blood_guess/internal/domain/constants"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/formatting"
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SetRealRoleHandler processes the /setrealrole command
type SetRealRoleHandler struct {
	*BaseHandler
	setRealRoleInput ports.SetRealRoleInput
}

func NewSetRealRoleHandler(
	bot BotClient,
	setRealRoleInput ports.SetRealRoleInput,
) *SetRealRoleHandler {
	return &SetRealRoleHandler{
		BaseHandler:      NewBaseHandler(bot),
		setRealRoleInput: setRealRoleInput,
	}
}

func (h *SetRealRoleHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "setrealrole")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := fmt.Sprintf(`<b>Использование:</b> <code>/setrealrole &lt;ID_игры&gt; &lt;ID_игрока&gt; &lt;реальная_роль&gt;</code>

<b>Пример:</b> <code>/setrealrole abc123 def456 demon</code>

<b>❗ Внимание:</b> Эта команда устанавливает РЕАЛЬНУЮ роль игрока после окончания игры.
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

	parts := strings.Fields(args)
	if len(parts) < 3 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Недостаточно аргументов. Используйте: <code>/setrealrole &lt;ID_игры&gt; &lt;ID_игрока&gt; &lt;реальная_роль&gt;</code>")
	}

	gameID := parts[0]
	playerSlotID := parts[1]
	role := parts[2]

	if !dto.IsValidGameRole(role) {
		validRoles := dto.GetValidGameRoles()
		var rolesList []string
		for _, r := range validRoles {
			rolesList = append(rolesList, fmt.Sprintf("<code>%s</code> %s", r, formatting.RoleEmoji(r)))
		}

		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf(`❌ Недопустимая роль: <code>%s</code>

<b>Допустимые реальные роли:</b>
%s`,
				h.EscapeHTML(role),
				strings.Join(rolesList, "\n")))
	}

	command := dto.SetRealRoleCommand{
		GameID:       gameID,
		PlayerSlotID: playerSlotID,
		RealRole:     role,
		AdminID:      update.Message.From.ID,
	}

	if err := h.setRealRoleInput.Execute(command); err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось установить реальную роль: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	roleEmoji := formatting.RoleEmoji(role)
	roleDisplayName := formatting.RoleDisplayName(role)

	isEvilRole := formatting.IsEvilRole(role)
	pointsInfo := ""

	if isEvilRole {
		pointsInfo = fmt.Sprintf(`

<b>📊 Влияние на очки пользователей:</b>
• Если кто-то предсказал для этого игрока правильную злую роль: <b>+%d очков</b> %s
• Если кто-то ошибся с злой ролью: <b>-%d очков</b>`,
			constants.PointsForRole(role),
			formatting.RoleEmoji(role),
			constants.PenaltyForRole(role))
	} else {
		pointsInfo = `

<b>📊 Влияние на очки пользователей:</b>
• Если кто-то предсказал для этого игрока злую роль: <b>штраф</b>`
	}

	successMsg := fmt.Sprintf(`✅ <b>Реальная роль установлена!</b>

<b>🎮 Игра:</b> <code>%s</code>
<b>👤 Игрок ID:</b> <code>%s</code>
<b>🎭 Реальная роль:</b> %s <b>%s</b> (%s)%s

Продолжайте устанавливать реальные роли для других игроков.`,
		h.EscapeHTML(gameID),
		h.EscapeHTML(playerSlotID),
		roleEmoji,
		h.EscapeHTML(role),
		roleDisplayName,
		pointsInfo)

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *SetRealRoleHandler) Command() string {
	return "setrealrole"
}

func (h *SetRealRoleHandler) Description() string {
	return "Установить реальную роль игрока после игры"
}
