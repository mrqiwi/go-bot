package telegram

import (
	"go.uber.org/zap"
	"strings"
	"time"

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
	logger      *zap.SugaredLogger
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
	logger *zap.SugaredLogger,
	bot *tgbotapi.BotAPI,
	commander Commander,
	httpClient HTTP.HTTPClient,
	transClient transmission.Transmission,
) TelegramController {
	return TelegramController{
		settings:    settings,
		logger:      logger,
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
		ctrl.logger.Errorf("Cannot get update channel: %s", err)
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !ctrl.verifyChatID(update.Message.Chat.ID) {
			ctrl.logger.Infof("Unknown user: %d", update.Message.Chat.ID)
			continue
		}

		if oldMessage(update.Message.Date) {
			continue
		}

		ctrl.HandleMessage(update.Message)
	}

	return nil
}

func (ctrl TelegramController) HandleMessage(msg *tgbotapi.Message) {
	if msg.Document != nil {
		err := ctrl.HandleDocument(msg)
		if err != nil {
			ctrl.logger.Errorf("Handle document: %s", err)
			return
		}
	}

	if msg.Text != "" {
		ctrl.HandleText(strings.ToLower(msg.Text), msg.Chat.ID)
		return
	}
}

func (ctrl *TelegramController) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := ctrl.bot.Send(msg)
	if err != nil {
		ctrl.logger.Errorf("Cannot send message: %s", err)
	}
}

func (ctrl TelegramController) verifyChatID(chatID int64) bool {
	for _, id := range ctrl.settings.ChatIDs {
		if id == chatID {
			return true
		}
	}

	return false
}

func oldMessage(timestamp int) bool {
	const waitTime = 2

	msgTime := time.Unix(int64(timestamp), 0)
	nowTime := time.Now()

	if nowTime.Sub(msgTime).Minutes() > waitTime {
		return true
	}

	return false
}
