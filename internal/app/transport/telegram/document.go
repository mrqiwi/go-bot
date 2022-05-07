package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (ctrl *TelegramController) HandleDocument(msg *tgbotapi.Message) error {
	switch msg.Document.MimeType {
	case "application/x-bittorrent":
		return ctrl.DownloadFile(msg.Document)
	default:
		return nil
	}
}

func (ctrl *TelegramController) DownloadFile(filename *tgbotapi.Document) error {
	file, err := ctrl.bot.GetFile(tgbotapi.FileConfig{FileID: filename.FileID})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", ctrl.downloadPath, filename.FileName)

	err = ctrl.httpClient.DownloadFile(path, file.Link(ctrl.token))
	if err != nil {
		return err
	}

	return nil
}
