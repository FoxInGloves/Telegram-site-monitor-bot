package main

import (
	"TelegramBot/config"
	"TelegramBot/logger"
	"TelegramBot/tgbot"
	"TelegramBot/web"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	ParseFlags()

	var wg sync.WaitGroup

	tomlConfig, getConfigError := config.GetTomlConfig()
	if getConfigError != nil {
		fmt.Println(getConfigError)
		return
	}

	logRespCh := make(chan web.Response, len(tomlConfig.Sites.Urls))
	errorsRespCh := make(chan web.Response, len(tomlConfig.Sites.Urls))
	defer close(logRespCh)
	defer close(errorsRespCh)

	wg.Add(3)
	go tgbot.Init(&tomlConfig.Telegram.BotToken, &tomlConfig.Telegram.ChatID, errorsRespCh)

	mutex := &sync.Mutex{}
	go logger.InitLogger(logRespCh, mutex)

	go func() {
		client := &http.Client{
			Timeout: time.Duration(tomlConfig.Settings.Timeout) * time.Second,
		}
		for {
			performRequests(tomlConfig, client, logRespCh, errorsRespCh)
			time.Sleep(time.Duration(tomlConfig.Settings.CheckInterval) * time.Second)
		}
	}()

	wg.Wait()
}

func ParseFlags() {
	path := flag.String("path", "config.toml", "Путь к конфигу")
	p := flag.String("p", "", "Путь к конфигу (короткая версия)")

	flag.Parse()

	if *p != "" {
		config.PathToConfig = *p
	} else {
		config.PathToConfig = *path
	}
}

func performRequests(tomlConfig *config.TomlConfig, client *http.Client, logRespCh, errorsRespCh chan<- web.Response) {
	updateConfigError := tomlConfig.UpdateConfig()
	if updateConfigError != nil {
		log.Println(updateConfigError)
		return
	}

	for _, url := range tomlConfig.Sites.Urls {
		go func(url string) {
			web.GetRequest(url, client, logRespCh, errorsRespCh)
		}(url)
	}
}
