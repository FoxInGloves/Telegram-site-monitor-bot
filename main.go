package main

import (
	"TelegramBot/config"
	"TelegramBot/logger"
	"TelegramBot/tgbot"
	"TelegramBot/web"
	"flag"
	"fmt"
	"sync"
)

func main() {
	ParseFlags()

	var wg sync.WaitGroup

	tomlConfig := config.GetTomlConfig()
	wg.Add(1)
	go tomlConfig.UpdateConfig()

	logRespCh := make(chan web.Response, len(tomlConfig.Sites.Urls))
	errorsRespCh := make(chan web.Response, len(tomlConfig.Sites.Urls))
	defer close(logRespCh)
	defer close(errorsRespCh)

	wg.Add(1)
	go tgbot.Init(&tomlConfig.Telegram.BotToken, &tomlConfig.Telegram.ChatId, errorsRespCh)

	mutex := &sync.Mutex{}
	wg.Add(1)
	go logger.InitLogger(logRespCh, mutex)

	wg.Add(1)
	go web.InfRequests(tomlConfig, logRespCh, errorsRespCh)

	wg.Wait()
}

func ParseFlags() {
	path := flag.String("path", "config.toml", "Путь к конфигу")
	p := flag.String("p", "config.toml", "Путь к конфигу")
	flag.Parse()
	if *path != "config.toml" && *p != "config.toml" {
		fmt.Println("Путь к конфигу взят из флага -path")
		config.PathToConfig = *path
	} else if *p != "config.toml" {
		config.PathToConfig = *p
	} else {
		config.PathToConfig = *path
	}
}
