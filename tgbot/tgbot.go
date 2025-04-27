package tgbot

import (
	"TelegramBot/config"
	"TelegramBot/logger"
	"TelegramBot/web"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
	"strings"
	"time"
)

func Init(botToken *string, chatId *int, channel <-chan web.Response) {
	bot, err := tgbotapi.NewBotAPI(*botToken)
	if err != nil {
		log.Panic("Ошибка инициализации бота:", err)
	}

	go handleCommands(bot, chatId)
	go handleResponses(bot, chatId, channel)
}

func handleResponses(bot *tgbotapi.BotAPI, chatId *int, channel <-chan web.Response) {
	for response := range channel {
		text := logger.FormatResponse(response)
		message := tgbotapi.NewMessage(int64(*chatId), text)

		if _, err := bot.Send(message); err != nil {
			log.Printf("error sending message: %v\n", err)
		}
	}
}

func handleCommands(bot *tgbotapi.BotAPI, chatId *int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic("Error receiving updates:", err)
	}

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "status":
			go handleStatusCommand(bot, chatId)
		default:
			go handleUnknownCommand(bot, chatId)
		}
	}
}

func handleStatusCommand(bot *tgbotapi.BotAPI, chatId *int) {
	tomlConfig, err := config.GetTomlConfig()
	if err != nil {
		log.Println(err)
	}
	client := &http.Client{Timeout: time.Duration(tomlConfig.Settings.Timeout) * time.Second}

	channel := make(chan web.Response, len(tomlConfig.Sites.Urls))
	defer close(channel)

	for _, url := range tomlConfig.Sites.Urls {
		go web.GetRequest(url, client, channel, nil)
	}

	var results []string
	for i := 0; i < len(tomlConfig.Sites.Urls); i++ {
		response := <-channel
		results = append(results, logger.FormatResponse(response))
	}

	text := strings.Join(results, "\n")
	message := tgbotapi.NewMessage(int64(*chatId), text)

	if _, err := bot.Send(message); err != nil {
		log.Printf("error sending message: %v\n", err)
	}
}

func handleUnknownCommand(bot *tgbotapi.BotAPI, chatId *int) {
	message := tgbotapi.NewMessage(int64(*chatId), "Неизвестная команда")
	if _, err := bot.Send(message); err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}
