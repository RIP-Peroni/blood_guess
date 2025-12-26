package telegram

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot represents Telegram bot
type Bot struct {
	api       *tgbotapi.BotAPI
	config    *Config
	handlers  map[string]CommandHandler
	startTime time.Time
}

// handleCommand handles a command
func (b *Bot) handleCommand(update tgbotapi.Update) {
	command := update.Message.Command()
	handler, exists := b.handlers[command]
	if !exists {
		b.SendMessage(update.Message.Chat.ID,
			fmt.Sprintf("Неизвестная команда: /%s\n\nИспользуйте /help для списка команд", command))
		return
	}

	if err := handler.Handle(update); err != nil {
		log.Printf("Error handling command /%s: %v", command, err)
		b.SendHTML(update.Message.Chat.ID,
			fmt.Sprintf("❌ Произошла ошибка при обработке команды.\n\nОшибка: %v", err))
	}
}

// handleMessage processes regular messages
func (b *Bot) handleMessage(update tgbotapi.Update) {
	msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Я понимаю только команды. Используй /help для списка команд")
	msgConfig.ReplyToMessageID = update.Message.MessageID

	if _, err := b.api.Send(msgConfig); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

// CommandHandler - interface of command handler
type CommandHandler interface {
	Handle(update tgbotapi.Update) error
	Command() string
	Description() string
}

// NewBot creates new bot
func NewBot(config *Config) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot API: %w", err)
	}

	botAPI.Debug = config.Debug

	return &Bot{
		api:       botAPI,
		config:    config,
		handlers:  make(map[string]CommandHandler),
		startTime: time.Now(),
	}, nil
}

// GetAPI returns BotAPI for using in handlers
func (b *Bot) GetAPI() *tgbotapi.BotAPI {
	return b.api
}

// RegisterHandler registers command's handler
func (b *Bot) RegisterHandler(handler CommandHandler) {
	b.handlers[handler.Command()] = handler
}

// Start launches a bot
func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.api.Self.UserName)
	log.Printf("Bot started at %s", b.startTime.Format(time.RFC3339))
	log.Printf("Registered handlers: %d", len(b.handlers))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update)
		} else {
			b.handleMessage(update)
		}
	}
	return nil
}

// Stop stops a bot
func (b *Bot) Stop() {
	log.Println("Stopping bot...")
	b.api.StopReceivingUpdates()
	log.Println("Bot stopped")
}

// SendMessage sends message to chat
func (b *Bot) SendMessage(chatID int64, text string) {
	msgConfig := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msgConfig); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
