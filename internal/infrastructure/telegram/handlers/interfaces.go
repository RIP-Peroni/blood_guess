package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type BotClient interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	GetMe() (tgbotapi.User, error)
}
