package handlers

import (
	"log"

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

// SendMarkdown sends a message with Markdown markup
func (s *MessageSender) SendMarkdown(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
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
	case "markdown":
		return h.sender.SendMarkdown(chatID, text)
	case "html":
		return h.sender.SendHTML(chatID, text)
	default:
		return h.sender.SendText(chatID, text)
	}
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
