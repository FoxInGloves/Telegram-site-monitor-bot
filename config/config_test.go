package config

import (
	"os"
	"testing"
)

func TestGetTomlConfig_Nil_DefaultPathToConfig(t *testing.T) {
	GetTomlConfig()
	if PathToConfig != "config.toml" {
		t.Error("Неверный путь к конфигу")
	}
}

func TestGetTomlConfig_PathToTempConfig_ConfigPointer(t *testing.T) {
	mockConfig := `
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
	tempConfig, err := os.CreateTemp("", "testConfig.toml")
	if err != nil {
		t.Fatalf("ошибка создания временного файла: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("<UNK> <UNK> <UNK> <UNK>: %v", err)
		}
	}(tempConfig.Name())
	_, writeConfigErr := tempConfig.Write([]byte(mockConfig))
	if writeConfigErr != nil {
		t.Fatalf("Ошибка записи в файл: %v\n", writeConfigErr)
	}
	closeFileError := tempConfig.Close()
	if closeFileError != nil {
		t.Fatalf("<UNK> <UNK> <UNK> <UNK>: %v", closeFileError)
	}

	PathToConfig = tempConfig.Name()
	config, getConfigErr := GetTomlConfig()
	if getConfigErr != nil {
		t.Fatalf("Ошибка получения конфигурации: %v", getConfigErr)
	}

	if config.Telegram.BotToken != "YOUR_TELEGRAM_BOT_TOKEN" {
		t.Errorf("Неверный токен бота, ожидалось %v, получено %v", "YOUR_TELEGRAM_BOT_TOKEN", config.Telegram.BotToken)
	}
	if config.Telegram.ChatID != 123456789 {
		t.Errorf("Неверный chat_id, ожидалось %v, получено %v", 123456789, config.Telegram.ChatID)
	}

	if len(config.Sites.Urls) != 2 || config.Sites.Urls[0] != "https://example.com" || config.Sites.Urls[1] != "https://google.com" {
		t.Errorf("Неверные URL-ы сайтов")
	}

	if config.Settings.CheckInterval != 300 {
		t.Errorf("Неверный check_interval, ожидалось %v, получено %v", 300, config.Settings.CheckInterval)
	}

	if config.Settings.Timeout != 10 {
		t.Errorf("Неверный timeout, ожидалось %v, получено %v", 10, config.Settings.Timeout)
	}
}

func TestGetTomlConfig_InvalidPath_Error(t *testing.T) {
	_, getConfigErr := GetTomlConfig()
	if getConfigErr == nil {
		t.Error("Ожидалась ошибка")
	}
}

func TestUpdateConfig_NewConfig_NewConfig(t *testing.T) {
	mockConfig := `
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
	tempConfig, err := os.CreateTemp("", "testConfigUpdate.toml")
	if err != nil {
		t.Fatalf("ошибка создания временного файла: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("<UNK> <UNK> <UNK> <UNK>: %v", err)
		}
	}(tempConfig.Name())
	_, writeConfigErr := tempConfig.Write([]byte(mockConfig))
	if writeConfigErr != nil {
		t.Fatalf("Ошибка записи в файл: %v\n", writeConfigErr)
	}

	closeFileError := tempConfig.Close()
	if closeFileError != nil {
		t.Fatalf("<UNK> <UNK> <UNK> <UNK>: %v", closeFileError)
	}

	PathToConfig = tempConfig.Name()
	config := TomlConfig{}
	updateConfigError := config.UpdateConfig()
	if updateConfigError != nil {
		t.Fatalf("Ошибка получения конфигурации: %v", updateConfigError)
	}

	if config.Telegram.BotToken != "YOUR_TELEGRAM_BOT_TOKEN" {
		t.Errorf("Неверный токен бота, ожидалось %v, получено %v", "YOUR_TELEGRAM_BOT_TOKEN", config.Telegram.BotToken)
	}
	if config.Telegram.ChatID != 123456789 {
		t.Errorf("Неверный chat_id, ожидалось %v, получено %v", 123456789, config.Telegram.ChatID)
	}

	if len(config.Sites.Urls) != 2 || config.Sites.Urls[0] != "https://example.com" || config.Sites.Urls[1] != "https://google.com" {
		t.Errorf("Неверные URL-ы сайтов")
	}

	if config.Settings.CheckInterval != 300 {
		t.Errorf("Неверный check_interval, ожидалось %v, получено %v", 300, config.Settings.CheckInterval)
	}

	if config.Settings.Timeout != 10 {
		t.Errorf("Неверный timeout, ожидалось %v, получено %v", 10, config.Settings.Timeout)
	}
}

func TestUpdateConfig_Nil_Error(t *testing.T) {
	config := TomlConfig{}
	updateConfigError := config.UpdateConfig()
	if updateConfigError == nil {
		t.Error("Ожидалась ошибка обновления конфига")
	}
}
