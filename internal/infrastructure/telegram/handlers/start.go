package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// StartHandler handles /start command
type StartHandler struct {
	*BaseHandler
}

func NewStartHandler(bot BotClient) *StartHandler {
	return &StartHandler{BaseHandler: NewBaseHandler(bot)}
}

func (h *StartHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "start")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	user := update.Message.From

	var username string
	if user.UserName != "" {
		username = "@" + user.UserName
	} else {
		username = user.FirstName
	}

	// Экранируем username для HTML
	escapedUsername := h.EscapeHTML(username)

	message := fmt.Sprintf(`<b>👋 Привет, %s!</b>

🎭 <b>Добро пожаловать в бота для угадывания ролей в игре "Кровь на Часовой Башне"!</b>

Я помогу тебе:
• Создавать игры и добавлять игроков
• Делать прогнозы на роли игроков
• Получать очки за прогнозы и менять их на школьные койнсы!

<b>📚</b> Используй <code>/help</code> для списка команд
<b>🎮</b> Используй <code>/newgame</code> для создания новой игры

<i>Удачных прогнозов!</i> 🎯`, escapedUsername)

	return h.SendHTML(update.Message.Chat.ID, message)
}

func (h *StartHandler) Command() string {
	return "start"
}

func (h *StartHandler) Description() string {
	return "Начать работу с ботом"
}
