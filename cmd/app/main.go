package main

import (
	"log"
	"os"

	"go-bot/internal/app/transport/telegram"
	"go-bot/internal/app/transport/vk"
	"go-bot/internal/app/usecase"

	"github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	token := os.Getenv("TELETOKEN")
	if token == "" {
		log.Fatalf("Telegram token is empty")
	}

	teleBot, err := telegram.TelegramApiInit(token)
	if err != nil {
		log.Fatal(err)
	}

	token = os.Getenv("VKTOKEN")
	if token == "" {
		log.Fatalf("VK token is empty")
	}

	commands := usecase.NewCommandProvider()
	teleController := telegram.NewTelegramController(teleBot, commands)
	vkController := vk.NewVkController(api.NewVK(token), commands)

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
