package logger

import (
	"TelegramSiteMonitorBot/web"
	"log"
	"os"
	"sync"
	"time"
)

func InitLogger(infoRespCh <-chan web.Response, mutex *sync.Mutex) {
	for {
		response := <-infoRespCh
		text := FormatResponse(response)
		go logToFile(text, mutex)
	}

}

func logToFile(str string, mutex *sync.Mutex) {
	fileName := time.Now().Format("2006-01-02") + ".log"

	mutex.Lock()
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE, 0666)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
		mutex.Unlock()
	}(file)

	if err != nil {
		log.Println(err)
		return
	}

	log.SetOutput(file)
	log.SetFlags(0)
	log.Println(str)
}
