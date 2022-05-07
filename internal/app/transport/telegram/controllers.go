package telegram

import (
	"log"
	"strings"

	HTTP "go-bot/internal/app/transport/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Commander interface {
	Reboot() error
	Poweroff() error
}

type TelegramController struct {
	bot          *tgbotapi.BotAPI
	commander    Commander
	httpClient   HTTP.HTTPClient
	token        string
	downloadPath string
}

func NewTelegramController(
	bot *tgbotapi.BotAPI,
	commander Commander,
	httpClient HTTP.HTTPClient,
	token string,
	downloadPath string,
) TelegramController {
	return TelegramController{
		bot:          bot,
		commander:    commander,
		httpClient:   httpClient,
		token:        token,
		downloadPath: downloadPath,
	}
}

func (ctrl *TelegramController) EventLoop() error {
	updates, err := ctrl.bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 60,
	})
	if err != nil {
		log.Printf("Cannot get update channel: %s", err)
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Document != nil {
			err = ctrl.HandleDocument(update.Message)
			if err != nil {
				log.Printf("Handle document: %s", err)
				continue
			}
		} else {
			ctrl.HandleMessage(strings.ToLower(update.Message.Text), update.Message.Chat.ID)
		}
	}

	return nil
}
