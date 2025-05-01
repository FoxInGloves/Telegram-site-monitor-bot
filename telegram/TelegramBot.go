package telegram

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"log"
)

type TgBot struct {
	api     *tgbotapi.BotAPI
	chatID  int64
	updates chan Update
}

func NewBot(token string, chatID int) (*TgBot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	tgBot := &TgBot{
		api:     api,
		chatID:  int64(chatID),
		updates: make(chan Update, 100),
	}

	go tgBot.listen()

	return tgBot, nil
}

func (bot *TgBot) SendMessage(text string) error {
	message := tgbotapi.NewMessage(bot.chatID, text)
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
