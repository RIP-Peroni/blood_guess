package handlers

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct {
	*BaseHandler
	commands map[string]string //command -> description
}

func NewHelpHandler(bot BotClient, commands map[string]string) *HelpHandler {
	return &HelpHandler{
		NewBaseHandler(bot),
		commands,
	}
}

func (h *HelpHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "help")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, h.buildHelpMessage())
	msgConfig.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := h.bot.Send(msgConfig)
	return err
}

func (h *HelpHandler) Command() string {
	return "help"
}

func (h *HelpHandler) Description() string {
	return "Показать список команд"
}

func (h *HelpHandler) buildHelpMessage() string {
	message := "🩸 *Доступные команды:*\n\n"
	for cmd, desc := range h.commands {
		message += fmt.Sprintf("* /%s - %s\n", cmd, desc)
	}

	message += "\n\n_Играй в Кровь на Часовой Башне, делай прогнозы и получай школьные койны!_🪙"
	return message
}
