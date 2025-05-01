package telegram

import (
	"TelegramSiteMonitorBot/config"
	"TelegramSiteMonitorBot/logger"
	"TelegramSiteMonitorBot/web"
	"log"
	"net/http"
	"strings"
	"time"
)

func RunTelegramBot(bot Bot, channel <-chan web.Response) {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	client := &http.Client{Timeout: time.Duration(cfg.Settings.Timeout) * time.Second}

	go handleCommands(bot, cfg, client)
	go handleWebResponses(bot, channel)
}

func handleWebResponses(bot Bot, channel <-chan web.Response) {
	for response := range channel {
		text := logger.FormatResponse(response)
		if err := bot.SendMessage(text); err != nil {
			log.Printf("error sending message: %v\n", err)
		}
	}
}

func handleCommands(bot Bot, cfg *config.AppConfig, client *http.Client) {
	for update := range bot.Updates() {
		if update.Command == "" {
			continue
		}

		switch strings.ToLower(update.Command) {
		case "status":
			go handleStatusCommand(bot, cfg.Sites.URLs, client)
		default:
			go handleUnknownCommand(bot)
		}
	}
}

func handleStatusCommand(bot Bot, URLs []string, client *http.Client) {

	responsesCh := make(chan web.Response, len(URLs))
	defer close(responsesCh)

	for _, url := range URLs {
		go web.GetRequest(url, client, responsesCh, nil)
	}

	var results []string
	for i := 0; i < len(URLs); i++ {
		response := <-responsesCh
		results = append(results, logger.FormatResponse(response))
	}

	text := strings.Join(results, "\n")

	if err := bot.SendMessage(text); err != nil {
		log.Printf("error sending message: %v", err)
	}
}

func handleUnknownCommand(bot Bot) {
	if err := bot.SendMessage("Неизвестная команда"); err != nil {
		log.Printf("error sending message: %v", err)
	}
}
