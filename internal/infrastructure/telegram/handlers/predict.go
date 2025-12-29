package handlers

import (
	"fmt"
	"strings"

	"RIP-Peroni/blood_guess/internal/application/dto"
	"RIP-Peroni/blood_guess/internal/application/ports"
	"RIP-Peroni/blood_guess/internal/domain/entities"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PredictHandler handles /predict command
type PredictHandler struct {
	*BaseHandler
	submitPredictionInput ports.SubmitPredictionInput
	gameRepo              ports.GameRepository
	userRepo              ports.UserRepository
}

func NewPredictHandler(
	bot BotClient,
	submitPredictionInput ports.SubmitPredictionInput,
	gameRepo ports.GameRepository,
	userRepo ports.UserRepository,
) *PredictHandler {
	return &PredictHandler{
		BaseHandler:           NewBaseHandler(bot),
		submitPredictionInput: submitPredictionInput,
		gameRepo:              gameRepo,
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
		message := `<b>Использование:</b> <code>/predict &lt;ID_игры&gt; &lt;ID_игрока&gt; &lt;роль&gt;</code>

<b>Пример:</b> <code>/predict abc123 def456 demon</code>

<b>Как получить ID игрока:</b>
Используйте команду <code>/games</code> для просмотра списка игр.
Затем используйте <code>/gameinfo &lt;ID_игры&gt;</code> для просмотра ID игроков.

<b>Доступные роли:</b>
• <code>townsfolk</code> - горожанин
• <code>outsider</code> - изгой
• <code>minion</code> - приспешник
• <code>demon</code> - демон

<b>Примечание:</b> Вы можете сделать только один прогноз на каждого игрока.`
		return h.SendHTML(update.Message.Chat.ID, message)
	}

	// Parse arguments: /predict <game_id> <player_id> <role>
	parts := strings.Fields(args)
	if len(parts) < 3 {
		return h.SendHTML(update.Message.Chat.ID,
			"❌ Недостаточно аргументов. Используйте: <code>/predict &lt;ID_игры&gt; &lt;ID_игрока&gt; &lt;роль&gt;</code>")
	}

	gameID := parts[0]
	playerSlotID := parts[1]
	role := parts[2]

	game, err := h.gameRepo.FindByID(entities.GameID(gameID))
	if err != nil {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Игра с ID <code>%s</code> не найдена.", h.EscapeHTML(gameID)))
	}

	if !game.CanAcceptPredictions() {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Прогнозы для этой игры не принимаются. Текущий статус: <b>%s</b>.", game.Status()))
	}

	playerFound := false
	var playerName string
	for _, player := range game.Players() {
		if string(player.ID) == playerSlotID {
			playerFound = true
			playerName = player.Name
			break
		}
	}

	if !playerFound {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Игрок с ID <code>%s</code> не найден в этой игре.", h.EscapeHTML(playerSlotID)))
	}

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

	command := dto.SubmitPredictionCommand{
		GameID:        gameID,
		UserID:        string(user.ID()),
		PlayerSlotID:  playerSlotID,
		PredictedRole: role,
	}

	response, err := h.submitPredictionInput.Execute(command)
	if err != nil {
		errorMsg := fmt.Sprintf("❌ Не удалось сохранить прогноз: %v", err)
		return h.SendText(update.Message.Chat.ID, errorMsg)
	}

	successMsg := fmt.Sprintf(`✅ <b>Прогноз сохранен!</b>

<b>🎮 Игра:</b> %s
<b>👤 Игрок:</b> %s
<b>🎭 Ваш прогноз:</b> %s

<b>📝 Ваши прогнозы в этой игре:</b>
ID прогноза: <code>%s</code>
Создан: %s

Теперь можно сделать прогнозы для других игроков или дождаться начала игры!`,
		h.EscapeHTML(game.Name()),
		h.EscapeHTML(playerName),
		h.EscapeHTML(role),
		h.EscapeHTML(response.ID),
		response.CreatedAt.Format("02.01.2006 15:04"))

	return h.SendHTML(update.Message.Chat.ID, successMsg)
}

func (h *PredictHandler) Command() string {
	return "predict"
}

func (h *PredictHandler) Description() string {
	return "Сделать прогноз на роль игрока в игре"
}
