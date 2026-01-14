package handlers

import (
	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/application/usecases"
	"RIP-Peroni/blood_guess/internal/domain/constants"
	"RIP-Peroni/blood_guess/internal/domain/entities"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/formatting"
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PredictHandler handles /predict command
type PredictHandler struct {
	*BaseHandler
	submitPredictionInput ports.SubmitPredictionInput
	gameFinder            usecases.GameFinderInterface
	userRepo              ports.UserRepository
}

func NewPredictHandler(
	bot BotClient,
	submitPredictionInput ports.SubmitPredictionInput,
	gameFinder usecases.GameFinderInterface,
	userRepo ports.UserRepository,
) *PredictHandler {
	return &PredictHandler{
		BaseHandler:           NewBaseHandler(bot),
		submitPredictionInput: submitPredictionInput,
		gameFinder:            gameFinder,
		userRepo:              userRepo,
	}
}

func (h *PredictHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "predict")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	args := strings.TrimSpace(update.Message.CommandArguments())
	if args == "" {
		message := fmt.Sprintf(`<b>Использование:</b> <code>/predict &lt;имя игрока&gt; &lt;роль&gt; [&lt;имя игрока&gt; &lt;роль&gt; ...]</code>

<b>Пример:</b> <code>/predict Вася demon Петя minion Коля minion</code>

<b>Примечание:</b>
• Можно сделать несколько прогнозов за один раз
• Имена игроков должны быть такими же, как при добавлении (регистр важен)
• Доступные роли для прогноза: <code>demon</code> и <code>minion</code>

<b>Система начисления очков:</b>
✅ Угадал демона: %d очков
✅ Угадал приспешника: %d очков
❌ Ошибся с демоном: %d очка
❌ Ошибся с приспешником: %d очка
🙅 Прогнозы на мирные роли не учитываются

Вы можете сделать только один прогноз на каждого игрока.`,
			constants.PointsForDemon,
			constants.PointsForMinion,
			constants.PenaltyForDemon,
			constants.PenaltyForMinion)
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Находим последнюю игру в статусе PREDICTIONS_OPEN
	game, err := h.gameFinder.FindLastGameByStatus(entities.GameStatusPredictionsOpen)
	if err != nil {
		if errors.Is(err, usecases.ErrNoGamesWithStatus) {
			return h.SendHTML(update.Message.Chat.ID,
				`❌ <b>Не найдена игра для прогнозов!</b>

Нет игр в статусе "прогнозы открыты". Возможные причины:
1. Игра еще не создана - используйте <code>/newgame</code>
2. Прогнозы еще не открыты - используйте <code>/openpred</code>
3. Игра уже перешла в другой статус - используйте <code>/games</code> для просмотра`)
		}
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Ошибка при поиске игры: %v", err))
	}

	// Получаем или создаем пользователя
	telegramID := update.Message.From.ID
	user, err := h.userRepo.FindByTelegramID(telegramID)
	if err != nil {
		username := update.Message.From.UserName
		if username == "" {
			username = update.Message.From.FirstName
		}
		user = entities.NewUser(telegramID, username)
		if err := h.userRepo.Save(user); err != nil {
			return h.SendHTML(update.Message.Chat.ID,
				"❌ Ошибка при создании пользователя. Попробуйте еще раз.")
		}
	}

	// Парсим аргументы: чередование имени игрока и роли
	parts := parseArguments(args)
	if len(parts)%2 != 0 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Нечетное количество аргументов. Ожидается формат: <code>/predict имя1 роль1 имя2 роль2 ...</code>")
	}

	var successfulPredictions []string
	var errorMessages []string

	// Обрабатываем пары (имя, роль)
	for i := 0; i < len(parts); i += 2 {
		playerName := parts[i]
		role := parts[i+1]

		// Находим игрока по имени в игре
		var playerSlotID string
		for _, player := range game.Players() {
			if player.Name == playerName {
				playerSlotID = string(player.ID)
				break
			}
		}

		if playerSlotID == "" {
			errorMessages = append(errorMessages, fmt.Sprintf("❌ Игрок '%s' не найден в игре", playerName))
			continue
		}

		if !dto.IsPredictableRole(role) {
			errorMessages = append(errorMessages, fmt.Sprintf("❌ Недопустимая роль для игрока '%s': %s", playerName, role))
			continue
		}

		command := dto.SubmitPredictionCommand{
			GameID:        string(game.ID()),
			UserID:        string(user.ID()),
			PlayerSlotID:  playerSlotID,
			PredictedRole: role,
		}

		_, err := h.submitPredictionInput.Execute(command)
		if err != nil {
			errorMessages = append(errorMessages, fmt.Sprintf("❌ Не удалось сохранить прогноз для '%s': %v", playerName, err))
		} else {
			successfulPredictions = append(successfulPredictions,
				fmt.Sprintf("✅ %s → %s", playerName, role))
		}
	}

	// Формируем итоговое сообщение
	var sb strings.Builder

	if len(successfulPredictions) > 0 {
		sb.WriteString(fmt.Sprintf(`%s <b>Прогнозы сохранены!</b>

<b>🎮 Игра:</b> %s
<b>👤 Ваши прогнозы:</b>
%s
`,
			formatting.RoleEmoji("demon"),
			h.EscapeHTML(game.Name()),
			strings.Join(successfulPredictions, "\n")))
	}

	if len(errorMessages) > 0 {
		sb.WriteString("\n<b>⚠️ Ошибки:</b>\n")
		sb.WriteString(strings.Join(errorMessages, "\n"))
	}

	if len(successfulPredictions) == 0 && len(errorMessages) == 0 {
		sb.WriteString("❌ Не удалось обработать ни одного прогноза.")
	}

	return h.SendHTML(update.Message.Chat.ID, sb.String())
}

func (h *PredictHandler) Command() string {
	return "predict"
}

func (h *PredictHandler) Description() string {
	return "Сделать прогноз на роли игроков в последней игре с открытыми прогнозами"
}
