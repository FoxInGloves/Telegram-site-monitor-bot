package main

import (
	"TelegramSiteMonitorBot/config"
	"TelegramSiteMonitorBot/logger"
	"TelegramSiteMonitorBot/telegram"
	"TelegramSiteMonitorBot/web"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	ParseFlags()

	tomlConfig, getConfigError := config.GetConfig()
	if getConfigError != nil {
		fmt.Println(getConfigError)
		return
	}

	logRespCh := make(chan web.Response, len(tomlConfig.Sites.URLs))
	errorsRespCh := make(chan web.Response, len(tomlConfig.Sites.URLs))
	defer close(logRespCh)
	defer close(errorsRespCh)

	var wg sync.WaitGroup

	wg.Add(3)
	mutex := &sync.Mutex{}
	go logger.InitLogger(logRespCh, mutex)

	telegramBot, err := telegram.NewBot(tomlConfig.Telegram.BotToken, tomlConfig.Telegram.ChatID)
	if err != nil {
		panic(err)
	}
	go telegram.RunTelegramBot(telegramBot, errorsRespCh)

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
	path := flag.String("path", "config.toml", "Path to config")
	p := flag.String("p", "", "Path to config (short version)")

	flag.Parse()

	if *p != "" {
		config.PathToConfig = *p
	} else {
		config.PathToConfig = *path
	}
}

func performRequests(tomlConfig *config.AppConfig, client *http.Client, logRespCh, errorsRespCh chan<- web.Response) {
	updateConfigError := tomlConfig.UpdateConfig()
	if updateConfigError != nil {
		log.Println(updateConfigError)
		return
	}

	for _, url := range tomlConfig.Sites.URLs {
		go func(url string) {
			web.GetRequest(url, client, logRespCh, errorsRespCh)
		}(url)
	}
}
