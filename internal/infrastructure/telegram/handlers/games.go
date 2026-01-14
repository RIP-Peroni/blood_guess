package handlers

import (
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/usecases"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GamesHandler struct {
	*BaseHandler
	gameFinder *usecases.GameFinder
}

func NewGamesHandler(bot BotClient, gameFinder *usecases.GameFinder) *GamesHandler {
	return &GamesHandler{
		BaseHandler: NewBaseHandler(bot),
		gameFinder:  gameFinder,
	}
}

func (h *GamesHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "games")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	// Используем GameFinder для получения активных игр через его репозиторий
	activeGames, err := h.gameFinder.GameRepo.FindActiveGames()
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Ошибка при получении списка игр: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	// ... остальной код без изменений
	if len(activeGames) == 0 {
		message := `<b>🎮 Активных игр нет</b>

Создайте новую игру с помощью <code>/newgame</code>`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Группируем игры по статусу
	var createdGames, openGames, closedGames, inProgressGames []*entities.Game

	for _, game := range activeGames {
		switch game.Status() {
		case entities.GameStatusCreated:
			createdGames = append(createdGames, game)
		case entities.GameStatusPredictionsOpen:
			openGames = append(openGames, game)
		case entities.GameStatusPredictionsClosed:
			closedGames = append(closedGames, game)
		case entities.GameStatusInProgress:
			inProgressGames = append(inProgressGames, game)
		}
	}

	var sb strings.Builder
	sb.WriteString("<b>🎮 Активные игры</b>\n\n")

	if len(createdGames) > 0 {
		sb.WriteString("<b>🆕 Созданные (ожидают открытия прогнозов):</b>\n")
		for i, game := range createdGames {
			escapedName := h.EscapeHTML(game.Name())
			escapedID := h.EscapeHTML(string(game.ID()))

			sb.WriteString(fmt.Sprintf("%d. %s (ID: <code>%s</code>)\n   👥 %d игроков\n",
				i+1, escapedName, escapedID, len(game.Players())))
		}
		sb.WriteString("\n")
	}

	if len(openGames) > 0 {
		sb.WriteString("<b>🎯 Открыты для прогнозов:</b>\n")
		for i, game := range openGames {
			escapedName := h.EscapeHTML(game.Name())
			escapedID := h.EscapeHTML(string(game.ID()))

			sb.WriteString(fmt.Sprintf("%d. %s (ID: <code>%s</code>)\n   👥 %d игроков\n",
				i+1, escapedName, escapedID, len(game.Players())))
		}
		sb.WriteString("\n")
	}

	if len(closedGames) > 0 {
		sb.WriteString("<b>⏸️ Прогнозы закрыты:</b>\n")
		for i, game := range closedGames {
			escapedName := h.EscapeHTML(game.Name())
			escapedID := h.EscapeHTML(string(game.ID()))

			sb.WriteString(fmt.Sprintf("%d. %s (ID: <code>%s</code>)\n   👥 %d игроков\n",
				i+1, escapedName, escapedID, len(game.Players())))
		}
		sb.WriteString("\n")
	}

	if len(inProgressGames) > 0 {
		sb.WriteString("<b>🎲 В процессе игры:</b>\n")
		for i, game := range inProgressGames {
			escapedName := h.EscapeHTML(game.Name())
			escapedID := h.EscapeHTML(string(game.ID()))

			sb.WriteString(fmt.Sprintf("%d. %s (ID: <code>%s</code>)\n   👥 %d игроков\n",
				i+1, escapedName, escapedID, len(game.Players())))
		}
	}

	sb.WriteString("\n")
	sb.WriteString("<b>📝</b> Используйте <code>/gameinfo &lt;ID_игры&gt;</code> для подробной информации")

	return h.SendHTML(update.Message.Chat.ID, sb.String())
}

func (h *GamesHandler) Command() string {
	return "games"
}

func (h *GamesHandler) Description() string {
	return "Показать активные игры"
}
