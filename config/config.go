package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
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
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	return &tomlConfig
}

func (tomlConfig *TomlConfig) UpdateConfig() string {
	for {
		time.Sleep(time.Duration(float32(tomlConfig.Settings.CheckInterval)/1.73) * time.Second)
		if _, err := toml.DecodeFile(PathToConfig, &tomlConfig); err != nil {
			fmt.Println("Error:", err)
			return "Конфиг поврежден! Проверьте конфиг и перезапустите программу."
		}
	}
}
