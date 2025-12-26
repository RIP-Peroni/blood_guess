package handlers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ===================== MessageSender ======================

// MessageSender provides convenient methods for sending messages
type MessageSender struct {
	bot BotClient
}

func NewMessageSender(bot BotClient) *MessageSender {
	return &MessageSender{bot: bot}
}

// SendText sends a simple text message
func (s *MessageSender) SendText(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := s.bot.Send(msg)
	return err
}

// SendHTML sends a message with HTML markup
func (s *MessageSender) SendHTML(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := s.bot.Send(msg)
	return err
}

// Reply responds to the message
func (s *MessageSender) Reply(update tgbotapi.Update, text string, parseMode string) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if parseMode != "" {
		msg.ParseMode = parseMode
	}
	msg.ReplyToMessageID = update.Message.MessageID
	_, err := s.bot.Send(msg)
	return err
}

// ====================== BaseHandler ======================

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	bot    BotClient
	sender *MessageSender
	logger *log.Logger
}

func NewBaseHandler(bot BotClient) *BaseHandler {
	return &BaseHandler{
		bot:    bot,
		sender: NewMessageSender(bot),
		logger: log.Default(),
	}
}

// SendMessage sends a message with the specified format
func (h *BaseHandler) SendMessage(chatID int64, text string, format string) error {
	switch format {
	case "html":
		return h.sender.SendHTML(chatID, text)
	default:
		return h.sender.SendText(chatID, text)
	}
}

// SendHTML is a shortcut for sending HTML messages
func (h *BaseHandler) SendHTML(chatID int64, text string) error {
	return h.sender.SendHTML(chatID, text)
}

// SendText is a shortcut for sending plain text messages
func (h *BaseHandler) SendText(chatID int64, text string) error {
	return h.sender.SendText(chatID, text)
}

// EscapeHTML escapes HTML special characters
func (h *BaseHandler) EscapeHTML(text string) string {
	// Экранируем базовые HTML символы
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	text = strings.ReplaceAll(text, "\"", "&quot;")
	text = strings.ReplaceAll(text, "'", "&#39;")
	return text
}

// LogCommand logs command usage
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

// RequireUser checks that the user exists
func (h *BaseHandler) RequireUser(update tgbotapi.Update) bool {
	return update.Message != nil && update.Message.From != nil
}

// GetUserID returns the user ID
func (h *BaseHandler) GetUserID(update tgbotapi.Update) int64 {
	if !h.RequireUser(update) {
		return 0
	}
	return update.Message.From.ID
}
