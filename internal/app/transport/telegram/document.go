package telegram

import (
	"fmt"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func (ctrl *TelegramController) HandleDocument(msg *tgbotapi.Message) error {
	switch msg.Document.MimeType {
	case "application/x-bittorrent":
		return ctrl.DownloadTorrent(msg)
	default:
		return nil
	}
}

func (ctrl *TelegramController) DownloadTorrent(msg *tgbotapi.Message) error {
	file, err := ctrl.bot.GetFile(tgbotapi.FileConfig{FileID: msg.Document.FileID})
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", ctrl.downloadPath, msg.Document.FileName)

	err = ctrl.httpClient.DownloadFile(path, file.Link(ctrl.token))
	if err != nil {
		return err
	}

	err = ctrl.transClient.AddTorrent(path)
	if err != nil {
		return err
	}

	ctrl.SendMessage(msg.Chat.ID, fmt.Sprintf("Torrent '%s' added for downloading", msg.Document.FileName))

	return nil
}
