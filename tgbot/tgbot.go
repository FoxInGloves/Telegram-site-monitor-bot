package tgbot

import (
	"TelegramBot/config"
	"TelegramBot/logger"
	"TelegramBot/web"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"time"
)

func Init(botToken *string, chatId *int, channel <-chan web.Response) {
	bot, err := tgbotapi.NewBotAPI(*botToken)
	if err != nil {
		log.Panic(err)
	}

	go InitCommands(bot, chatId)

	go func() {
		for {
			response := <-channel
			text := logger.FormatResponse(response)

			message := tgbotapi.NewMessage(int64(*chatId), text)
			_, err := bot.Send(message)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}()
}

func InitCommands(bot *tgbotapi.BotAPI, chatId *int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	bot.Debug = true

	for update := range updates {
		switch update.Message.Command() {
		case "status":
			go func() {
				tomlConfig := config.GetTomlConfig()
				urls := tomlConfig.Sites.Urls
				client := &http.Client{Timeout: time.Duration(tomlConfig.Settings.Timeout) * time.Second}

				channel := make(chan web.Response, 2)
				unusedChannel := make(chan web.Response, 2)
				for i := 0; i < len(urls); i++ {
					go web.GetRequest(urls[i], client, channel, unusedChannel)
				}

				text := ""
				for index := 0; index < len(urls); index++ {
					value := <-channel
					text += logger.FormatResponse(value) + "\n"
				}
				close(unusedChannel)
				close(channel)

				message := tgbotapi.NewMessage(int64(*chatId), text)
				_, err := bot.Send(message)
				if err != nil {
					log.Println(err.Error())
				}
			}()

		}
	}
}
