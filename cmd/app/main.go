package main

import (
	"log"

	"go-bot/internal/app/appconfig"
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

	httpClient := HTTP.NewHTTPClient()
	commands := usecase.NewCommandProvider()
	teleController := telegram.NewTelegramController(teleBot, commands, httpClient, config.TelegramToken, config.DownloadsPath)
	vkController := vk.NewVkController(api.NewVK(config.VKToken), commands)

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
