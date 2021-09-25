package telegram

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Commander interface {
	Reboot() error
	Poweroff() error
}

type TelegramController struct {
	bot       *tgbotapi.BotAPI
	commander Commander
}

func NewTelegramController(bot *tgbotapi.BotAPI, commander Commander) TelegramController {
	return TelegramController{
		bot:       bot,
		commander: commander,
	}
}

func (ctrl *TelegramController) EventLoop() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := ctrl.bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Printf("Cannot get update channel: %s", err)
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		ctrl.handleMessage(strings.ToLower(update.Message.Text), update.Message.Chat.ID)
	}

	return nil
}

func (ctrl *TelegramController) Reboot(chatID int64) {
	err := ctrl.commander.Reboot()
	if err != nil {
		log.Printf("Cannot run reboot: %s", err)
		return
	}

	ctrl.sendMessage(chatID, "reboot...")
}

func (ctrl *TelegramController) Poweroff(chatID int64) {
	err := ctrl.commander.Poweroff()
	if err != nil {
		log.Printf("Cannot run poweroff: %s", err)
		return
	}

	ctrl.sendMessage(chatID, "shutdowing...")
}

func (ctrl *TelegramController) handleMessage(text string, chatID int64) {
	switch text {
	case "reset", "reboot", "restart", "перезагрузка", "рестарт":
		ctrl.Reboot(chatID)
	case "poweroff", "off", "выключение", "выкл":
		ctrl.Poweroff(chatID)
	default:
		ctrl.sendMessage(chatID, "unknown command")
	}
}

func (ctrl *TelegramController) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := ctrl.bot.Send(msg)
	if err != nil {
		log.Printf("Cannot send message: %s", err)
	}
}
