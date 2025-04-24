package web

import (
	"TelegramBot/config"
	"io"
	"log"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	DateTime   string
	Url        string
}

func InfRequests(tomlConfig *config.TomlConfig, chLog chan<- Response, chErrors chan<- Response) {
	client := http.Client{Timeout: time.Duration(tomlConfig.Settings.Timeout) * time.Second}
	for {
		for i := 0; i < len(tomlConfig.Sites.Urls); i++ {
			go GetRequest(tomlConfig.Sites.Urls[i], &client, chLog, chErrors)
		}
		time.Sleep(time.Duration(tomlConfig.Settings.CheckInterval) * time.Second)
	}
}

func GetRequest(url string, client *http.Client, chLog chan<- Response, chErrors chan<- Response) {
	resp, err := client.Get(url)
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Println(err)
			}
		}(resp.Body)
	}

	if err != nil {
		log.Println(err)
		return
	}

	now := time.Now()
	dateTime := now.Format("2006-01-02 15:04:05")
	response := Response{StatusCode: resp.StatusCode, DateTime: dateTime, Url: url}

	if response.StatusCode > 400 {
		chLog <- response
		chErrors <- response
	} else {
		chLog <- response
	}
}
