package telegram

import (
	"log"
	"strings"

	"go-bot/internal/app/transmission"
	HTTP "go-bot/internal/app/transport/http"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Commander interface {
	Reboot() error
	Poweroff() error
}

type TelegramController struct {
	settings    Settings
	bot         *tgbotapi.BotAPI
	commander   Commander
	httpClient  HTTP.HTTPClient
	transClient transmission.Transmission
}

type Settings struct {
	Token        string
	DownloadPath string
	ChatIDs      []int64
}

func NewTelegramController(
	settings Settings,
	bot *tgbotapi.BotAPI,
	commander Commander,
	httpClient HTTP.HTTPClient,
	transClient transmission.Transmission,
) TelegramController {
	return TelegramController{
		settings:    settings,
		bot:         bot,
		commander:   commander,
		httpClient:  httpClient,
		transClient: transClient,
	}
}

func (ctrl TelegramController) EventLoop() error {
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

		if !ctrl.verifyChatID(update.Message.Chat.ID) {
			log.Printf("Unknown user: %d", update.Message.Chat.ID)
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

func (ctrl TelegramController) verifyChatID(chatID int64) bool {
	for _, id := range ctrl.settings.ChatIDs {
		if id == chatID {
			return true
		}
	}

	return false
}
