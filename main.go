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
	parseFlags()

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

	wg.Add(1)
	mutex := &sync.Mutex{}
	go logger.InitLogger(logRespCh, mutex)

	var telegramBot telegram.Bot
	var initBotErr, initBotWithProxyErr error

	telegramBot, initBotErr = tryInitTelegramBot(tomlConfig)
	if initBotErr != nil {
		log.Println(initBotErr)
		fmt.Println("Failed to connect to the bot")

		telegramBot, initBotWithProxyErr = tryInitTelegramBotWithProxy(tomlConfig)
		if initBotWithProxyErr != nil {
			log.Println(initBotWithProxyErr)
			fmt.Println("Failed to connect to bot with proxy")
			return
		}
	}

	wg.Add(2)
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

func parseFlags() {
	path := flag.String("path", "config.toml", "Path to config")
	p := flag.String("p", "", "Path to config (short version)")

	flag.Parse()

	if *p != "" {
		config.PathToConfig = *p
	} else {
		config.PathToConfig = *path
	}
}

func tryInitTelegramBot(cfg *config.AppConfig) (telegram.Bot, error) {
	telegramBot, err := telegram.NewBot(cfg.Telegram.BotToken, cfg.Telegram.ChatID)
	if err != nil {
		return nil, err
	}

	return telegramBot, nil
}

func tryInitTelegramBotWithProxy(cfg *config.AppConfig) (telegram.Bot, error) {
	proxyAddress := cfg.Settings.Proxy

	transport, err := web.GetWebTransport(proxyAddress)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	telegramBot, tgErr := telegram.NewBotWithProxy(cfg.Telegram.BotToken, cfg.Telegram.ChatID, client)
	if tgErr != nil {
		return nil, tgErr
	}

	return telegramBot, nil
}

func performRequests(tomlConfig *config.AppConfig, client *http.Client, logRespCh, errorsRespCh chan<- web.Response) {
	updateConfigError := tomlConfig.UpdateConfig()
	if updateConfigError != nil {
		log.Println(updateConfigError)
		return
	}

	for _, siteUrl := range tomlConfig.Sites.URLs {
		go func(url string) {
			web.GetRequest(url, client, logRespCh, errorsRespCh)
		}(siteUrl)
	}
}
