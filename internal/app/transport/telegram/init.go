package telegram

import (
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func TelegramApiInit(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("Cannot init telegram bot api: %s", err)
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot, err
}
