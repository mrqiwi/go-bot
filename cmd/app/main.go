package main

import (
	"log"

	"go-bot/internal/app/appconfig"
	"go-bot/internal/app/logging"
	"go-bot/internal/app/transmission"
	HTTP "go-bot/internal/app/transport/http"
	"go-bot/internal/app/transport/telegram"
	"go-bot/internal/app/transport/vk"
	"go-bot/internal/app/usecase"

	"github.com/SevereCloud/vksdk/v2/api"
)

func main() {
	logger, loggerCleanup, err := logging.InitLogger()
	if err != nil {
		log.Fatalf("init logger failed: %s", err)
	}

	config, err := appconfig.ReadConfig()
	if err != nil {
		logger.Fatalf("cannot read the app appconfig: %s", err)
	}

	teleBot, err := telegram.TelegramApiInit(config.TelegramToken, logger)
	if err != nil {
		logger.Fatalf("cannot init telegram api: %s", err)
	}

	transClient := transmission.TranmissionClient()
	httpClient := HTTP.NewHTTPClient()
	commands := usecase.NewCommandProvider()
	vkController := vk.NewVkController(logger, api.NewVK(config.VKToken), commands)
	teleController := telegram.NewTelegramController(
		telegram.Settings{
			Token:        config.TelegramToken,
			DownloadPath: config.DownloadsPath,
			ChatIDs:      config.ChatIDs,
		},
		logger,
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

	log.Print(<-errChannel)
	loggerCleanup()
}
