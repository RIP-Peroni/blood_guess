package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GameInfoHandler handles /gameinfo command
type GameInfoHandler struct {
	*BaseHandler
	gameRepo       ports.GameRepository
	userRepo       ports.UserRepository
	predictionRepo ports.PredictionRepository
}

func NewGameInfoHandler(
	bot BotClient,
	gameRepo ports.GameRepository,
	userRepo ports.UserRepository,
	predictionRepo ports.PredictionRepository,
) *GameInfoHandler {
	return &GameInfoHandler{
		BaseHandler:    NewBaseHandler(bot),
		gameRepo:       gameRepo,
		userRepo:       userRepo,
		predictionRepo: predictionRepo,
	}
}

func (h *GameInfoHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "gameinfo")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := `<b>Использование:</b> <code>/gameinfo &lt;ID_игры&gt;</code>

<b>Пример:</b> <code>/gameinfo abc123</code>

Используйте <code>/games</code> для просмотра списка игр и их ID.`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	gameID := args

	game, err := h.gameRepo.FindByID(entities.GameID(gameID))
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Игра с ID <code>%s</code> не найдена.", h.EscapeHTML(gameID)))
	}

	telegramID := update.Message.From.ID
	user, err := h.userRepo.FindByTelegramID(telegramID)
	var userID entities.UserID
	if err == nil {
		userID = user.ID()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>🎮 Информация об игре</b>\n\n"))
	sb.WriteString(fmt.Sprintf("<b>📛 Название:</b> %s\n", h.EscapeHTML(game.Name())))
	sb.WriteString(fmt.Sprintf("<b>🆔 ID игры:</b> <code>%s</code>\n", h.EscapeHTML(gameID)))
	sb.WriteString(fmt.Sprintf("<b>📊 Статус:</b> %s\n", game.Status()))
	sb.WriteString(fmt.Sprintf("<b>👑 Создатель:</b> ID %d\n", game.CreatorID()))
	sb.WriteString(fmt.Sprintf("<b>📅 Создана:</b> %s\n", game.CreatedAt().Format("02.01.2006 15:04")))

	sb.WriteString(fmt.Sprintf("\n<b>👥 Игроки (%d):</b>\n", len(game.Players())))
	if len(game.Players()) == 0 {
		sb.WriteString("   (нет игроков)\n")
	} else {
		for i, player := range game.Players() {
			sb.WriteString(fmt.Sprintf("%d. <b>%s</b>\n", i+1, h.EscapeHTML(player.Name)))
			sb.WriteString(fmt.Sprintf("   🆔: <code>%s</code>\n", player.ID))

			if player.IsRealRoleSet && (game.Status() == entities.GameStatusFinished) {
				sb.WriteString(fmt.Sprintf("   ✅ Реальная роль: <b>%s</b>\n", player.RealRole))
			}

			if userID != "" {
				predictions, err := h.predictionRepo.FindByGameAndUser(game.ID(), userID)
				if err == nil {
					for _, pred := range predictions {
						if string(pred.PlayerSlotID()) == string(player.ID) {
							sb.WriteString(fmt.Sprintf("   🎯 Ваш прогноз: %s\n", pred.PredictedRole()))
							if points, awarded := pred.PointsAwarded(); awarded {
								sb.WriteString(fmt.Sprintf("   📊 Начислено очков: %d\n", points))
							}
						}
					}
				}
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n<b>📝 Доступные действия:</b>\n")
	switch game.Status() {
	case entities.GameStatusCreated:
		sb.WriteString("• Используйте <code>/addplayer</code> чтобы добавить игроков\n")
		sb.WriteString("• Используйте <code>/openpred</code> чтобы открыть прогнозы\n")
	case entities.GameStatusPredictionsOpen:
		sb.WriteString("• Используйте <code>/predict</code> чтобы сделать прогноз\n")
		sb.WriteString("• Используйте <code>/closepred</code> чтобы закрыть прогнозы\n")
	case entities.GameStatusPredictionsClosed:
		sb.WriteString("• Прогнозы закрыты. Ожидание начала реальной игры.\n")
	case entities.GameStatusInProgress:
		sb.WriteString("• Игра в процессе. Устанавливайте реальные роли с помощью <code>/setrole</code>\n")
	case entities.GameStatusFinished:
		sb.WriteString("• Игра завершена. Результаты подсчитаны.\n")
	}

	return h.SendHTML(update.Message.Chat.ID, sb.String())
}

func (h *GameInfoHandler) Command() string {
	return "gameinfo"
}

func (h *GameInfoHandler) Description() string {
	return "Показать подробную информацию об игре"
}
