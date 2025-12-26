package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ==================== MessageSender ====================

// MessageSender предоставляет удобные методы для отправки сообщений
type MessageSender struct {
	bot *tgbotapi.BotAPI
}

func NewMessageSender(bot *tgbotapi.BotAPI) *MessageSender {
	return &MessageSender{bot: bot}
}

// SendText отправляет простое текстовое сообщение
func (s *MessageSender) SendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.bot.Send(msg)
	return err
}

// SendMarkdown отправляет сообщение с Markdown разметкой
func (s *MessageSender) SendMarkdown(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := s.bot.Send(msg)
	return err
}

// SendHTML отправляет сообщение с HTML разметкой
func (s *MessageSender) SendHTML(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := s.bot.Send(msg)
	return err
}

// Reply отвечает на сообщение
func (s *MessageSender) Reply(update tgbotapi.Update, text string, parseMode string) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if parseMode != "" {
		msg.ParseMode = parseMode
	}
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := s.bot.Send(msg)
	return err
}

// ==================== BaseHandler ====================

// BaseHandler предоставляет общую функциональность для всех обработчиков
type BaseHandler struct {
	bot    *tgbotapi.BotAPI
	sender *MessageSender
	logger *log.Logger
}

func NewBaseHandler(bot *tgbotapi.BotAPI) *BaseHandler {
	return &BaseHandler{
		bot:    bot,
		sender: NewMessageSender(bot),
		logger: log.Default(),
	}
}

// SendMessage отправляет сообщение с указанным форматом
func (h *BaseHandler) SendMessage(chatID int64, text string, format string) error {
	switch format {
	case "markdown":
		return h.sender.SendMarkdown(chatID, text)
	case "html":
		return h.sender.SendHTML(chatID, text)
	default:
		return h.sender.SendText(chatID, text)
	}
}

// LogCommand логирует использование команды
func (h *BaseHandler) LogCommand(update tgbotapi.Update, command string) {
	user := update.Message.From
	h.logger.Printf("Command /%s from user %d (%s %s @%s)",
		command,
		user.ID,
		user.FirstName,
		user.LastName,
		user.UserName,
	)
}

// RequireUser проверяет, что пользователь существует
func (h *BaseHandler) RequireUser(update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.From != nil
}

// GetUserID возвращает ID пользователя
func (h *BaseHandler) GetUserID(update tgbotapi.Update) int64 {
	if !h.RequireUser(update) {
		return 0
	}
	return update.Message.From.ID
}
