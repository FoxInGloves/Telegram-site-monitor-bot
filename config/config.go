package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type AppConfig struct {
	Telegram telegramConfig
	Sites    sitesConfig
	Settings settingsConfig
}

type telegramConfig struct {
	BotToken string `toml:"bot_token"`
	ChatID   int    `toml:"chat_id"`
}

type sitesConfig struct {
	URLs []string `toml:"urls"`
}

type settingsConfig struct {
	CheckInterval int `toml:"check_interval"`
	Timeout       int `toml:"timeout"`
}

var PathToConfig = "config.toml"

func GetConfig() (*AppConfig, error) {
	var tomlConfig AppConfig
	if _, err := toml.DecodeFile(PathToConfig, &tomlConfig); err != nil {
		configError := fmt.Errorf("Config is corrupted! Check the config and restart the program.\nError: %w", err)
		return nil, configError
	}
	return &tomlConfig, nil
}

func (tomlConfig *AppConfig) UpdateConfig() error {
	if _, err := toml.DecodeFile(PathToConfig, &tomlConfig); err != nil {
		return fmt.Errorf("Config is corrupted! Check the config and restart the program.\nError: %w", err)
	}
	return nil
}
