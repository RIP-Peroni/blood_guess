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

// CopyPlayersHandler processes the /copyplayers command
type CopyPlayersHandler struct {
	*BaseHandler
	copyPlayersInput ports.CopyPlayersInput
	gameFinder       usecases.GameFinderInterface
}

func NewCopyPlayersHandler(
	bot BotClient,
	copyPlayersInput ports.CopyPlayersInput,
	gameFinder usecases.GameFinderInterface,
) *CopyPlayersHandler {
	return &CopyPlayersHandler{
		BaseHandler:      NewBaseHandler(bot),
		copyPlayersInput: copyPlayersInput,
		gameFinder:       gameFinder,
	}
}

func (h *CopyPlayersHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "copyplayers")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args != "" {
		// Старый формат с ID игры - для обратной совместимости
		_ = h.SendHTML(update.Message.Chat.ID,
			`<i>Примечание: Теперь команда /copyplayers не требует ID игры. Она автоматически копирует игроков в последнюю созданную игру.</i>`)
	}

	// Находим последнюю игру в статусе CREATED
	targetGame, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusCreated)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для копирования игроков!</b>

Нет игр в статусе "создана". Возможные причины:
1. Игра еще не создана - используйте <code>/newgame</code>
2. Игра уже перешла в другой статус - используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	if targetGame.CreatorID() != update.Message.From.ID {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Только создатель игры может копировать игроков.")
	}

	command := dto.CopyPlayersCommand{
		TargetGameID: string(targetGame.ID()),
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

Теперь можно добавить ещё игроков или открыть прогнозы с помощью <code>/openpred</code>

<b>Примечание:</b> Роли игрокам не назначены. Вы можете назначить их командой <code>/setrealrole</code> после завершения игры.`,
		h.EscapeHTML(response.Name),
		len(playerNames),
		strings.Join(playerNames, ", "))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *CopyPlayersHandler) Command() string {
	return "copyplayers"
}

func (h *CopyPlayersHandler) Description() string {
	return "Скопировать игроков из последней завершенной игры"
}
