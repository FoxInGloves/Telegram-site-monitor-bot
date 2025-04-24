package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"os"
	"time"
)

type TomlConfig struct {
	Telegram TelegramConfig
	Sites    SitesConfig
	Settings SettingsConfig
}

type TelegramConfig struct {
	BotToken string `toml:"bot_token"`
	ChatId   int    `toml:"chat_id"`
}

type SitesConfig struct {
	Urls []string `toml:"urls"`
}

type SettingsConfig struct {
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
