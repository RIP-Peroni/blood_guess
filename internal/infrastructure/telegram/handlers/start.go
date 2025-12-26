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
		username = fmt.Sprintf("@%s", user.UserName)
	} else {
		username = user.FirstName
	}

	message := fmt.Sprintf(`👋 Привет, %s!

🎭 *Добро пожаловать в бота для угадывания ролей в игре "Кровь на Часовой Башне"!*

Я помогу тебе:
• Создавать игры и добавлять игроков
• Делать прогнозы на роли игроков
• Получать очки за прогнозы и менять их на школьные койнсы!

📚 Используй /help для списка команд
🎮 Используй /newgame для создания новой игры

_Удачных прогнозов!_ 🎯`, username)
	msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	msgConfig.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := h.bot.Send(msgConfig)
	return err
}

func (h *StartHandler) Command() string {
	return "start"
}

func (h *StartHandler) Description() string {
	return "Начать работу с ботом"
}
