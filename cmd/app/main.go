package main

import (
	"log"
	"os"

	HTTP "go-bot/internal/app/transport/http"
	"go-bot/internal/app/transport/telegram"
	"go-bot/internal/app/transport/vk"
	"go-bot/internal/app/usecase"

	"github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	downloadPath := os.Getenv("DOWNLOADS")
	if downloadPath == "" {
		log.Fatal("Download path is not set")
	}

	teleToken := os.Getenv("TELETOKEN")
	if teleToken == "" {
		log.Fatal("Telegram token is empty")
	}

	teleBot, err := telegram.TelegramApiInit(teleToken)
	if err != nil {
		log.Fatal(err)
	}

	vkToken := os.Getenv("VKTOKEN")
	if vkToken == "" {
		log.Fatal("VK token is empty")
	}

	httpClient := HTTP.NewHTTPClient()
	commands := usecase.NewCommandProvider()
	teleController := telegram.NewTelegramController(teleBot, commands, httpClient, teleToken, downloadPath)
	vkController := vk.NewVkController(api.NewVK(vkToken), commands)

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
