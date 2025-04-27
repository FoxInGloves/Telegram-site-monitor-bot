package web

import (
	"log"
	"net/http"
	"time"
)

type Response struct {
	StatusCode int
	DateTime   string
	Url        string
}

func GetRequest(url string, client *http.Client, chLog, chErrors chan<- Response) {
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("error requesting %s: %v\n", url, err)
		sendErrorResponse(url, chErrors)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response from %s: %v\n", url, err)
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
