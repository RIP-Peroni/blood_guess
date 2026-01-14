package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/application/usecases"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// MyPredictHandler обрабатывает команду /mypredict
type MyPredictHandler struct {
	*BaseHandler
	gameRepo       ports.GameRepository
	userRepo       ports.UserRepository
	predictionRepo ports.PredictionRepository
	gameFinder     usecases.GameFinderInterface
}

func NewMyPredictHandler(
	bot BotClient,
	gameRepo ports.GameRepository,
	userRepo ports.UserRepository,
	predictionRepo ports.PredictionRepository,
	gameFinder usecases.GameFinderInterface,
) *MyPredictHandler {
	return &MyPredictHandler{
		BaseHandler:    NewBaseHandler(bot),
		gameRepo:       gameRepo,
		userRepo:       userRepo,
		predictionRepo: predictionRepo,
		gameFinder:     gameFinder,
	}
}

func (h *MyPredictHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "mypredict")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	telegramID := update.Message.From.ID
	user, err := h.userRepo.FindByTelegramID(telegramID)
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Вы еще не делали прогнозов. Сначала создайте прогноз с помощью <code>/predict</code>")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	var game *entities.Game

	if args == "" {
		// Пытаемся найти последнюю игру с открытыми или закрытыми прогнозами
		game, err = h.gameFinder.FindLastGameByStatus(entities.GameStatusPredictionsOpen)
		if err != nil {
			game, err = h.gameFinder.FindLastGameByStatus(entities.GameStatusPredictionsClosed)
			if err != nil {
				game, err = h.gameFinder.FindLatestGame()
				if err != nil {
					return h.SendHTML(update.Message.Chat.ID,
						"❌ Не найдено ни одной игры. Сначала создайте игру с помощью <code>/newgame</code>")
				}
			}
		}
	} else {
		// Ищем игру по ID
		game, err = h.gameRepo.FindByID(entities.GameID(args))
		if err != nil {
			return h.SendHTML(update.Message.Chat.ID,
				fmt.Sprintf("❌ Игра с ID <code>%s</code> не найдена.", h.EscapeHTML(args)))
		}
	}

	// Получаем прогнозы пользователя для этой игры
	predictions, err := h.predictionRepo.FindByGameAndUser(game.ID(), user.ID())
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Ошибка при получении прогнозов. Попробуйте еще раз.")
	}

	if len(predictions) == 0 {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf(`❌ У вас нет прогнозов в игре "%s".

Сделайте прогноз с помощью команды: <code>/predict</code>`,
				h.EscapeHTML(game.Name())))
	}

	// Собираем информацию об игроках для отображения
	playerMap := make(map[string]string) // playerSlotID -> playerName
	for _, player := range game.Players() {
		playerMap[string(player.ID)] = player.Name
	}

	// Формируем сообщение
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<b>🎯 Ваши прогнозы</b>

<b>🎮 Игра:</b> %s
<b>📊 Статус игры:</b> %s
<b>👤 Ваш прогноз:</b>
`,
		h.EscapeHTML(game.Name()),
		game.Status()))

	for _, prediction := range predictions {
		playerName := playerMap[string(prediction.PlayerSlotID())]
		if playerName == "" {
			playerName = "Неизвестный игрок"
		}

		roleEmoji := "❓"
		switch prediction.PredictedRole().String() {
		case "demon":
			roleEmoji = "👹"
		case "minion":
			roleEmoji = "😈"
		}

		pointsInfo := ""
		if points, awarded := prediction.PointsAwarded(); awarded {
			sign := "+"
			if points < 0 {
				sign = ""
			}
			pointsInfo = fmt.Sprintf(" (%s%d очков)", sign, points)
		}

		sb.WriteString(fmt.Sprintf("• <b>%s</b> → %s %s%s\n",
			h.EscapeHTML(playerName),
			roleEmoji,
			prediction.PredictedRole().String(),
			pointsInfo))
	}

	sb.WriteString(fmt.Sprintf("\n<b>📈 Всего прогнозов:</b> %d", len(predictions)))

	if game.Status() == entities.GameStatusFinished && game.AllRealRolesSet() {
		sb.WriteString("\n\n<i>Очки уже начислены. Используйте <code>/profile</code> для просмотра баланса.</i>")
	} else if game.Status() == entities.GameStatusPredictionsOpen {
		sb.WriteString("\n\n<i>Вы можете изменить прогноз до закрытия прогнозов.</i>")
	}

	return h.SendHTML(update.Message.Chat.ID, sb.String())
}

func (h *MyPredictHandler) Command() string {
	return "mypredict"
}

func (h *MyPredictHandler) Description() string {
	return "Показать мои прогнозы в текущей или указанной игре"
}
