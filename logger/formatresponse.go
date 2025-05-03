package logger

import (
	"TelegramSiteMonitorBot/web"
	"fmt"
	"strconv"
)

func FormatResponse(response web.Response) string {
	var status string
	if response.StatusCode < 400 && response.StatusCode != 0 {
		status = "доступен"
	} else {
		status = "недоступен"
	}

	text := fmt.Sprintf("[%s] Сайт %s %s (Код: %s)", response.DateTime, response.Url, status,
		strconv.Itoa(response.StatusCode))
	return text
}
