package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"time"
)

type TomlConfig struct {
	Telegram telegramConfig
	Sites    sitesConfig
	Settings settingsConfig
}

type telegramConfig struct {
	BotToken string `toml:"bot_token"`
	ChatId   int    `toml:"chat_id"`
}

type sitesConfig struct {
	Urls []string `toml:"urls"`
}

type settingsConfig struct {
	CheckInterval int `toml:"check_interval"`
	Timeout       int `toml:"timeout"`
}

var PathToConfig string

func GetTomlConfig() *TomlConfig {
	if PathToConfig == "" {
		PathToConfig = "config.toml"
	}

	var tomlConfig TomlConfig
	if _, err := toml.DecodeFile(PathToConfig, &tomlConfig); err != nil {
		log.Println("Конфиг поврежден! Проверьте конфиг и перезапустите программу.")
		log.Println("Error:", err)
		os.Exit(1)
	}
	return &tomlConfig
}

func (tomlConfig *TomlConfig) UpdateConfig() string {
	for {
		time.Sleep(time.Duration(tomlConfig.Settings.CheckInterval) * time.Second)
		if _, err := toml.DecodeFile(PathToConfig, &tomlConfig); err != nil {
			log.Println("Конфиг поврежден! Проверьте конфиг и перезапустите программу.")
			log.Println("Error:", err)
			os.Exit(1)
		}
	}
}
