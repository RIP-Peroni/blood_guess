package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AddPlayerHandler handles /addplayer command
type AddPlayerHandler struct {
	*BaseHandler
	addPlayerInput ports.AddPlayerInput
	gameRepo       ports.GameRepository
}

func NewAddPlayerHandler(bot BotClient, addPlayerInput ports.AddPlayerInput, gameRepo ports.GameRepository) *AddPlayerHandler {
	return &AddPlayerHandler{
		BaseHandler:    NewBaseHandler(bot),
		addPlayerInput: addPlayerInput,
		gameRepo:       gameRepo,
	}
}

func (h *AddPlayerHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "addplayer")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/addplayer &lt;ID_игры&gt; &lt;имя_игрока&gt; &lt;роль&gt;</code>

<b>Пример:</b> <code>/addplayer abc123 Иван townsfolk</code>

<b>Доступные роли:</b>
• <code>townsfolk</code> - горожанин
• <code>outsider</code> - изгой
• <code>minion</code> - приспешник
• <code>demon</code> - демон

<b>Примечание:</b> Имя игрока должно быть уникальным в рамках игры.`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	parts := strings.Fields(args)
	if len(parts) < 3 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Недостаточно аргументов. Используйте: <code>/addplayer &lt;ID_игры&gt; &lt;имя_игрока&gt; &lt;роль&gt;</code>")
	}

	gameID := parts[0]
	playerName := parts[1]
	role := parts[2]

	game, err := h.gameRepo.FindByID(entities.GameID(gameID))
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Игра с ID <code>%s</code> не найдена.", h.EscapeHTML(gameID)))
	}

	// Проверяем, является ли пользователь создателем игры
	if game.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может добавлять игроков.")
	}

	// Проверяем валидность роли
	validRoles := []string{"townsfolk", "outsider", "minion", "demon"}
	isValidRole := false
	for _, validRole := range validRoles {
		if role == validRole {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Недопустимая роль: <code>%s</code>. Допустимые роли: townsfolk, outsider, minion, demon.", h.EscapeHTML(role)))
	}

	command := dto.AddPlayerCommand{
		GameID:       gameID,
		PlayerName:   playerName,
		AssignedRole: role,
		AdminID:      update.Message.From.ID,
	}

	response, err := h.addPlayerInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось добавить игрока: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	successMsg := fmt.Sprintf(`✅ <b>Игрок добавлен!</b>

<b>🎮 Игра:</b> %s
<b>👤 Игрок:</b> %s
<b>🎭 Роль:</b> %s
<b>👥 Всего игроков в игре:</b> %d

Теперь можно добавить ещё игроков или открыть прогнозы с помощью <code>/openpred %s</code>`,
		h.EscapeHTML(response.Name),
		h.EscapeHTML(playerName),
		h.EscapeHTML(role),
		len(response.Players),
		h.EscapeHTML(gameID))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *AddPlayerHandler) Command() string {
	return "addplayer"
}

func (h *AddPlayerHandler) Description() string {
	return "Добавить игрока в игру"
}
