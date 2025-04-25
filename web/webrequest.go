package web

import (
	"TelegramBot/config"
	"log"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	DateTime   string
	Url        string
}

func InfRequests(tomlConfig *config.TomlConfig, chLog, chErrors chan<- Response) {
	client := &http.Client{
		Timeout: time.Duration(tomlConfig.Settings.Timeout) * time.Second,
	}

	for {
		for _, url := range tomlConfig.Sites.Urls {
			go func(url string) {
				GetRequest(url, client, chLog, chErrors)
			}(url)
		}

		time.Sleep(time.Duration(tomlConfig.Settings.CheckInterval) * time.Second)
	}
}

func GetRequest(url string, client *http.Client, chLog, chErrors chan<- Response) {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Ошибка запроса к %s: %v\n", url, err)
		sendErrorResponse(url, chErrors)
		return
	}
	defer func() {
		if resp.Body != nil {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Ошибка закрытия ответа от %s: %v\n", url, err)
			}
		}
	}()

	response := Response{
		StatusCode: resp.StatusCode,
		DateTime:   time.Now().Format("2006-01-02 15:04:05"),
		Url:        url,
	}

	if response.StatusCode >= 400 {
		chErrors <- response
	}
	chLog <- response
}

func sendErrorResponse(url string, chErrors chan<- Response) {
	response := Response{
		StatusCode: 0,
		DateTime:   time.Now().Format("2006-01-02 15:04:05"),
		Url:        url,
	}
	chErrors <- response
}
