package config

import (
	"os"
	"testing"
)

var mockConfig = `
[telegram]  
bot_token = "YOUR_TELEGRAM_BOT_TOKEN"  
chat_id = 123456789  

[sites]  
urls = [  
  "https://example.com",  
  "https://google.com",  
]  

[settings]  
check_interval = 300
timeout = 10
`

func TestGetTomlConfig_PathToTempConfig_ConfigPointer(t *testing.T) {
	tempConfig, err := os.CreateTemp("", "testConfig.toml")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		removeErr := os.Remove(name)
		if removeErr != nil {
			t.Fatal(removeErr.Error())
		}
	}(tempConfig.Name())
	_, writeConfigErr := tempConfig.Write([]byte(mockConfig))
	if writeConfigErr != nil {
		t.Fatalf("error writing to file: %v\n", writeConfigErr)
	}
	closeFileError := tempConfig.Close()
	if closeFileError != nil {
		t.Fatalf("close file error %v", closeFileError)
	}

	PathToConfig = tempConfig.Name()
	config, getConfigErr := GetConfig()
	if getConfigErr != nil {
		t.Fatalf("error getting configuration: %v", getConfigErr)
	}

	if config.Telegram.BotToken != "YOUR_TELEGRAM_BOT_TOKEN" {
		t.Errorf("invalid bot token, expected %v, received %v", "YOUR_TELEGRAM_BOT_TOKEN", config.Telegram.BotToken)
	}
	if config.Telegram.ChatID != 123456789 {
		t.Errorf("invalid chat_id, expected %v, received %v", 123456789, config.Telegram.ChatID)
	}

	if len(config.Sites.URLs) != 2 || config.Sites.URLs[0] != "https://example.com" || config.Sites.URLs[1] != "https://google.com" {
		t.Errorf("invalid site URLs")
	}

	if config.Settings.CheckInterval != 300 {
		t.Errorf("invalid check_interval, expected %v, received %v", 300, config.Settings.CheckInterval)
	}

	if config.Settings.Timeout != 10 {
		t.Errorf("invalid timeout, expected %v, received %v", 10, config.Settings.Timeout)
	}
}

func TestGetTomlConfig_InvalidPath_Error(t *testing.T) {
	_, getConfigErr := GetConfig()
	if getConfigErr == nil {
		t.Error("error expected")
	}
}

func TestUpdateConfig_NewConfig_NewConfig(t *testing.T) {
	tempConfig, err := os.CreateTemp("", "testConfigUpdate.toml")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("file deletion error: %v", err)
		}
	}(tempConfig.Name())
	_, writeConfigErr := tempConfig.Write([]byte(mockConfig))
	if writeConfigErr != nil {
		t.Fatalf("error writing to file: %v", writeConfigErr)
	}

	closeFileError := tempConfig.Close()
	if closeFileError != nil {
		t.Fatalf("error closing file: %v", closeFileError)
	}

	PathToConfig = tempConfig.Name()
	config := AppConfig{}
	updateConfigError := config.UpdateConfig()
	if updateConfigError != nil {
		t.Fatalf("error getting configuration: %v", updateConfigError)
	}

	if config.Telegram.BotToken != "YOUR_TELEGRAM_BOT_TOKEN" {
		t.Errorf("invalid bot token, expected %v, received %v", "YOUR_TELEGRAM_BOT_TOKEN", config.Telegram.BotToken)
	}
	if config.Telegram.ChatID != 123456789 {
		t.Errorf("invalid chat_id, expected %v, received %v", 123456789, config.Telegram.ChatID)
	}

	if len(config.Sites.URLs) != 2 || config.Sites.URLs[0] != "https://example.com" || config.Sites.URLs[1] != "https://google.com" {
		t.Errorf("invalid site URLs")
	}

	if config.Settings.CheckInterval != 300 {
		t.Errorf("invalid check_interval, expected %v, received %v", 300, config.Settings.CheckInterval)
	}

	if config.Settings.Timeout != 10 {
		t.Errorf("invalid timeout, expected %v, received %v", 10, config.Settings.Timeout)
	}
}

func TestUpdateConfig_Nil_Error(t *testing.T) {
	config := AppConfig{}
	updateConfigError := config.UpdateConfig()
	if updateConfigError == nil {
		t.Error("config update error expected")
	}
}
