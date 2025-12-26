package handlers

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct {
	*BaseHandler
	commands map[string]string // command -> description
}

func NewHelpHandler(bot BotClient, commands map[string]string) *HelpHandler {
	return &HelpHandler{
		BaseHandler: NewBaseHandler(bot),
		commands:    commands,
	}
}

func (h *HelpHandler) Handle(update tgbotapi.Update) error {
	h.LogCommand(update, "help")

	if !h.RequireUser(update) {
		return h.SendMessage(update.Message.Chat.ID, "Ошибка: не удалось определить пользователя", "")
	}

	message := h.buildHelpMessage()
	return h.SendHTML(update.Message.Chat.ID, message)
}

func (h *HelpHandler) buildHelpMessage() string {
	var sb strings.Builder
	sb.WriteString("<b>🩸 Доступные команды:</b>\n\n")

	for cmd, desc := range h.commands {
		// Экранируем описание команды
		escapedDesc := h.EscapeHTML(desc)
		sb.WriteString(fmt.Sprintf("• <code>/%s</code> - %s\n", cmd, escapedDesc))
	}

	sb.WriteString("\n\n")
	sb.WriteString("<i>Играй в Кровь на Часовой Башне, делай прогнозы и получай школьные койны!</i> 🪙")

	return sb.String()
}

func (h *HelpHandler) Command() string {
	return "help"
}

func (h *HelpHandler) Description() string {
	return "Показать список команд"
}
