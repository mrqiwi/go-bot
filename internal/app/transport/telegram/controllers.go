package telegram

import (
	"go.uber.org/zap"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Commander interface {
	PcReboot() error
	PcOn() error
	PcOff() error
	PcStatus() error
}

type TelegramController struct {
	settings  Settings
	logger    *zap.SugaredLogger
	bot       *tgbotapi.BotAPI
	commander Commander
}

type Settings struct {
	Token   string
	ChatIDs []int64
}

func NewTelegramController(
	settings Settings,
	logger *zap.SugaredLogger,
	bot *tgbotapi.BotAPI,
	commander Commander,
) TelegramController {
	return TelegramController{
		settings:  settings,
		logger:    logger,
		bot:       bot,
		commander: commander,
	}
}

func (ctrl TelegramController) EventLoop() error {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить меню"},
		{Command: "pc_on", Description: "Включить ПК"},
		{Command: "pc_off", Description: "Выключить ПК"},
		{Command: "pc_reboot", Description: "Перезагрузить ПК"},
		{Command: "pc_status", Description: "Узнать статус ПК"},
	}
	_, err := ctrl.bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		ctrl.logger.Errorf("Cannot requst commands: %s", err)
	}

	menu := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/pc_on"),
			tgbotapi.NewKeyboardButton("/pc_off"),
			tgbotapi.NewKeyboardButton("/pc_reboot"),
			tgbotapi.NewKeyboardButton("/pc_status"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/help"),
		),
	)

	updates := ctrl.bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 60,
	})

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

		ctrl.handleCommand(update.Message, menu)
	}

	return nil
}

func (ctrl TelegramController) handleCommand(msg *tgbotapi.Message, menu tgbotapi.ReplyKeyboardMarkup) {
	switch msg.Command() {
	case "start":
		ctrl.sendMessageWithMenu(menu, msg.Chat.ID, "Привет! Выбери команду с клавиатуры ниже:")

	case "pc_on":
		err := ctrl.commander.PcOn()
		if err != nil {
			ctrl.sendMessage(msg.Chat.ID, "ПК невозможно включить")
		} else {
			ctrl.sendMessage(msg.Chat.ID, "ПК включается…")
		}

	case "pc_off":
		err := ctrl.commander.PcOff()
		if err != nil {
			ctrl.sendMessage(msg.Chat.ID, "ПК невозможно выключить")
		} else {
			ctrl.sendMessage(msg.Chat.ID, "ПК выключается…")
		}

	case "pc_reboot":
		err := ctrl.commander.PcReboot()
		if err != nil {
			ctrl.sendMessage(msg.Chat.ID, "ПК невозможнжо перезагрузить")
		} else {
			ctrl.sendMessage(msg.Chat.ID, "ПК перезагружается…")
		}

	case "pc_status":
		err := ctrl.commander.PcStatus()
		if err != nil {
			ctrl.sendMessage(msg.Chat.ID, "ПК выключен")
		} else {
			ctrl.sendMessage(msg.Chat.ID, "ПК включен")
		}

	case "help":
		ctrl.sendMessage(msg.Chat.ID, "Используй команды:\n/pc_on\n/pc_off\n/pc_reboot\n/pc_status")

	default:
		ctrl.sendMessage(msg.Chat.ID, "Неизвестная команда. Нажми /start, чтобы открыть меню.")
	}
}

func (ctrl TelegramController) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := ctrl.bot.Send(msg)
	if err != nil {
		ctrl.logger.Errorf("Cannot send message: %s", err)
	}
}

func (ctrl TelegramController) sendMessageWithMenu(menu tgbotapi.ReplyKeyboardMarkup, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	msg.ReplyMarkup = menu
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
