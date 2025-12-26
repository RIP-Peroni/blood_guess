package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendText sends a plain text message (no HTML)
func (b *Bot) SendText(chatID int64, text string) {
	msgConfig := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msgConfig); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// SendHTML sends a message with HTML markup
func (b *Bot) SendHTML(chatID int64, text string) {
	msgConfig := tgbotapi.NewMessage(chatID, text)
	msgConfig.ParseMode = tgbotapi.ModeHTML
	if _, err := b.api.Send(msgConfig); err != nil {
		log.Printf("Error sending HTML message: %v", err)
	}
}
