package telegram

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
	"net/http"
)

type TgBot struct {
	api     *tgbotapi.BotAPI
	chatId  int64
	updates chan Update
}

func NewBot(token string, chatId int) (*TgBot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return newBotInstance(api, chatId), nil
}

func NewBotWithProxy(token string, chatId int, client *http.Client) (*TgBot, error) {
	api, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		return nil, err
	}

	return newBotInstance(api, chatId), nil
}

func newBotInstance(api *tgbotapi.BotAPI, chatId int) *TgBot {
	tgBot := &TgBot{
		api:     api,
		chatId:  int64(chatId),
		updates: make(chan Update, 100),
	}

	go tgBot.listen()

	return tgBot
}

func (bot *TgBot) SendMessage(text string) error {
	message := tgbotapi.NewMessage(bot.chatId, text)
	_, err := bot.api.Send(message)
	return err
}

func (bot *TgBot) Updates() <-chan Update {
	return bot.updates
}

func (bot *TgBot) listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.api.GetUpdatesChan(u)
	if err != nil {
		log.Println("error receiving updates:", err)
		return
	}

	for update := range updates {
		if update.Message != nil {
			bot.updates <- Update{
				Text:    update.Message.Text,
				Command: update.Message.Command(),
			}
		}
	}
}
