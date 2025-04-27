package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Telegram telegramConfig
	Sites    sitesConfig
	Settings settingsConfig
}

type telegramConfig struct {
	BotToken string `toml:"bot_token"`
	ChatID   int    `toml:"chat_id"`
}

type sitesConfig struct {
	Urls []string `toml:"urls"`
}

type settingsConfig struct {
	CheckInterval int `toml:"check_interval"`
	Timeout       int `toml:"timeout"`
}

var PathToConfig string

func GetTomlConfig() (*TomlConfig, error) {
	if PathToConfig == "" {
		PathToConfig = "config.toml"
	}

	var tomlConfig TomlConfig
	if _, err := toml.DecodeFile(PathToConfig, &tomlConfig); err != nil {
		configError := fmt.Errorf("Config is corrupted! Check the config and restart the program.\nError: %w", err)
		return nil, configError
	}
	return &tomlConfig, nil
}

func (tomlConfig *TomlConfig) UpdateConfig() error {
	if _, err := toml.DecodeFile(PathToConfig, &tomlConfig); err != nil {
		return fmt.Errorf("Config is corrupted! Check the config and restart the program.\nError: %w", err)
	}
	return nil
}
