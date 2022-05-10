package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
)

func TelegramApiInit(token string, logger *zap.SugaredLogger) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	//bot.Debug = true

	logger.Infof("Authorized on account %s", bot.Self.UserName)

	return bot, err
}
