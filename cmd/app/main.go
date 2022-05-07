package main

import (
	"log"

	"go-bot/internal/app/appconfig"
	"go-bot/internal/app/transmission"
	HTTP "go-bot/internal/app/transport/http"
	"go-bot/internal/app/transport/telegram"
	"go-bot/internal/app/transport/vk"
	"go-bot/internal/app/usecase"

	"github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	config, err := appconfig.ReadConfig()
	if err != nil {
		log.Fatalf("cannot read the app appconfig: %s", err)
	}

	teleBot, err := telegram.TelegramApiInit(config.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	transClient := transmission.TranmissionClient()
	httpClient := HTTP.NewHTTPClient()
	commands := usecase.NewCommandProvider()
	vkController := vk.NewVkController(api.NewVK(config.VKToken), commands)
	teleController := telegram.NewTelegramController(
		telegram.Settings{
			Token:        config.TelegramToken,
			DownloadPath: config.DownloadsPath,
			ChatIDs:      config.ChatIDs,
		},
		teleBot,
		commands,
		httpClient,
		transClient,
	)

	errChannel := make(chan error)

	go func(errChannel chan error) {
		errVK := vkController.EventLoop()
		if errVK != nil {
			errChannel <- errVK
		}
	}(errChannel)

	go func(errChannel chan error) {
		errTele := teleController.EventLoop()
		if errTele != nil {
			errChannel <- errTele
		}
	}(errChannel)

	log.Fatal(<-errChannel)
}
