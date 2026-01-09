package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AssignRolesHandler processes the /assignroles command
type AssignRolesHandler struct {
	*BaseHandler
	gameRepo ports.GameRepository
}

func NewAssignRolesHandler(
	bot BotClient,
	gameRepo ports.GameRepository,
) *AssignRolesHandler {
	return &AssignRolesHandler{
		BaseHandler: NewBaseHandler(bot),
		gameRepo:    gameRepo,
	}
}

func (h *AssignRolesHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "assignroles")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/assignroles &lt;ID_игры&gt;</code>

<b>Пример:</b> <code>/assignroles abc123</code>

<b>Что делает:</b>
• Показывает список игроков без назначенных ролей
• Позволяет назначить роли интерактивно

<b>Доступные роли:</b>
• <code>townsfolk</code> - горожанин
• <code>outsider</code> - изгой
• <code>minion</code> - приспешник
• <code>demon</code> - демон

<b>Примечание:</b> Только создатель игры может назначать роли.`

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
			"❌ Только создатель игры может назначать роли.")
	}

	var playersWithoutRoles []entities.PlayerSlot
	for _, player := range game.Players() {
		if player.AssignedRole == "" {
			playersWithoutRoles = append(playersWithoutRoles, player)
		}
	}

	if len(playersWithoutRoles) == 0 {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf(`✅ <b>Все игроки уже имеют назначенные роли!</b>

<b>🎮 Игра:</b> %s
<b>👥 Игроков:</b> %d

Все роли назначены. Можете открывать прогнозы командой <code>/openpred %s</code>`,
				h.EscapeHTML(game.Name()),
				len(game.Players()),
				h.EscapeHTML(gameID)))
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<b>🎭 Назначение ролей для игры: %s</b>

<b>👥 Игроки без ролей (%d):</b>
`,
		h.EscapeHTML(game.Name()),
		len(playersWithoutRoles)))

	for i, player := range playersWithoutRoles {
		sb.WriteString(fmt.Sprintf("\n%d. <b>%s</b> (ID: <code>%s</code>)", i+1, h.EscapeHTML(player.Name), player.ID))
	}

	sb.WriteString(fmt.Sprintf(`

<b>📝 Для назначения роли используйте команду:</b>
<code>/setrole %s &lt;ID_игрока&gt; &lt;роль&gt;</code>

<b>Пример:</b> <code>/setrole %s %s demon</code>

<b>Доступные роли:</b>
• townsfolk - горожанин
• outsider - изгой
• minion - приспешник
• demon - демон`,
		h.EscapeHTML(gameID),
		h.EscapeHTML(gameID),
		h.EscapeHTML(string(playersWithoutRoles[0].ID))))

	return h.SendHTML(update.Message.Chat.ID, sb.String())
}

func (h *AssignRolesHandler) Command() string {
	return "assignroles"
}

func (h *AssignRolesHandler) Description() string {
	return "Показать игроков без ролей и инструкцию по назначению ролей"
}
