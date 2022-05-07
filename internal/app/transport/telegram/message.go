package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (ctrl *TelegramController) Reboot(chatID int64) {
	err := ctrl.commander.Reboot()
	if err != nil {
		log.Printf("Cannot run reboot: %s", err)
		return
	}

	ctrl.SendMessage(chatID, "reboot...")
}

func (ctrl *TelegramController) Poweroff(chatID int64) {
	err := ctrl.commander.Poweroff()
	if err != nil {
		log.Printf("Cannot run poweroff: %s", err)
		return
	}

	ctrl.SendMessage(chatID, "shutdowing...")
}

func (ctrl *TelegramController) HandleMessage(text string, chatID int64) {
	switch text {
	case "reset", "reboot", "restart", "перезагрузка", "рестарт":
		ctrl.Reboot(chatID)
	case "poweroff", "off", "выключение", "выкл":
		ctrl.Poweroff(chatID)
	default:
		ctrl.SendMessage(chatID, "unknown command")
	}
}

func (ctrl *TelegramController) SendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := ctrl.bot.Send(msg)
	if err != nil {
		log.Printf("Cannot send message: %s", err)
	}
}
