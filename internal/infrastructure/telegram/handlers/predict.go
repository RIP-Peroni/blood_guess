package handlers

import (
	"RIP-Peroni/blood_guess/internal/domain/constants"
	"RIP-Peroni/blood_guess/internal/infrastructure/telegram/formatting"
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
		message := fmt.Sprintf(`<b>Использование:</b> <code>/predict &lt;ID_игры&gt; &lt;ID_игрока&gt; &lt;роль&gt;</code>

<b>Пример:</b> <code>/predict abc123 def456 demon</code>

<b>Как получить ID игрока:</b>
Используйте команду <code>/games</code> для просмотра списка игр.
Затем используйте <code>/gameinfo &lt;ID_игры&gt;</code> для просмотра ID игроков.

<b>Доступные роли для прогноза:</b>
• <code>demon</code> - демон (злая роль)
• <code>minion</code> - приспешник (злая роль)

<b>❗ Внимание:</b> Прогнозировать можно только злые роли! Горожане и изгои не прогнозируются.

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

	if !dto.IsPredictableRole(role) {
		return h.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf(`❌ Недопустимая роль для прогноза: <code>%s</code> %s

<b>Можно прогнозировать только злые роли:</b>
• <code>demon</code> - демон %s
• <code>minion</code> - приспешник %s

Горожане (townsfolk) и изгои (outsider) <b>не прогнозируются</b>.`,
				h.EscapeHTML(role),
				formatting.RoleEmoji(role),
				formatting.RoleEmoji("demon"),
				formatting.RoleEmoji("minion")))
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

	roleEmoji := formatting.RoleEmoji(role)

	successMsg := fmt.Sprintf(`%s <b>Прогноз на злую роль сохранен!</b>

<b>🎮 Игра:</b> %s
<b>👤 Игрок:</b> %s
<b>🎭 Ваш прогноз:</b> %s %s (%s)

<b>📊 Система очков:</b>
• Если %s окажется демоном: <b>+%d очков</b> %s
• Если %s окажется приспешником: <b>+%d очков</b> %s
• Если ошибётесь: <b>штраф %d очков</b>

<b>📝 Ваши прогнозы в этой игре:</b>
ID прогноза: <code>%s</code>
Создан: %s

Теперь можно сделать прогнозы для других игроков!`,
		roleEmoji,
		h.EscapeHTML(game.Name()),
		h.EscapeHTML(playerName),
		roleEmoji,
		h.EscapeHTML(role),
		formatting.RoleDisplayName(role),
		h.EscapeHTML(playerName),
		constants.PointsForDemon,
		formatting.RoleEmoji("demon"),
		h.EscapeHTML(playerName),
		constants.PointsForMinion,
		formatting.RoleEmoji("minion"),
		constants.PenaltyForRole(role),
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
