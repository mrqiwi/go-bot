package appconfig

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	VKToken          string  `envconfig:"vktoken"`
	TelegramToken    string  `envconfig:"teletoken"`
	ChatIDs          []int64 `envconfig:"chat_ids"`
	PCAddress        string  `envconfig:"pc_address"`
	UserName         string  `envconfig:"user_name"`
	UserPassword     string  `envconfig:"user_password"`
	MacAddress       string  `envconfig:"mac_address"`
	BroadcastAddress string  `envconfig:"broadcast_address"`
}

func ReadConfig() (Config, error) {
	var c Config

	err := envconfig.Process("go-bot", &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
