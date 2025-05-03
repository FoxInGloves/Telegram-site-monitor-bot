package main

import (
	"TelegramSiteMonitorBot/config"
	"TelegramSiteMonitorBot/telegram"
	"TelegramSiteMonitorBot/web"
	"bytes"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Mock
type MockBot struct {
	Token  string
	ChatID int
}

func (m *MockBot) SendMessage(message string) error {
	if message == "" {
		return errors.New("пустое сообщение")
	}
	return nil
}

func (m *MockBot) Updates() <-chan telegram.Update {
	return nil
}

func MockNewBot(token string, chatID int) (telegram.Bot, error) {
	if token == "" || chatID == 0 {
		return nil, errors.New("incorrect parameters")
	}
	return &MockBot{Token: token, ChatID: chatID}, nil
}

type mockRoundTripper struct {
	StatusCode   int
	ErrorToThrow error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.ErrorToThrow != nil {
		return nil, m.ErrorToThrow
	}
	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(strings.NewReader("ok")),
		Header:     make(http.Header),
	}, nil
}

var mockConfig = `
[telegram]  
bot_token = "YOUR_TELEGRAM_BOT_TOKEN"  
chat_id = 123456789  

[sites]  
urls = [  
  "https://example.com",  
  "https://google.com",  
]  

[settings]  
check_interval = 300
timeout = 10
`

func TestParseFlags(t *testing.T) {
	tests := []struct {
		args     []string
		expected string
	}{
		{[]string{"cmd", "-path=test_config.toml"}, "test_config.toml"},
		{[]string{"cmd", "-p=short_config.toml"}, "short_config.toml"},
		{[]string{"cmd"}, "config.toml"},
	}

	for _, tt := range tests {

		os.Args = tt.args
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		parseFlags()

		if config.PathToConfig != tt.expected {
			t.Errorf("expected %v, received %v", tt.expected, config.PathToConfig)
		}
	}
}

func TestPerformRequests_SuccessResponse(t *testing.T) {
	tempConfig, err := os.CreateTemp("", "testConfig.toml")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		removeErr := os.Remove(name)
		if removeErr != nil {
			t.Fatal(removeErr.Error())
		}
	}(tempConfig.Name())
	_, writeConfigErr := tempConfig.Write([]byte(mockConfig))
	if writeConfigErr != nil {
		t.Fatalf("error writing to file: %v\n", writeConfigErr)
	}
	closeFileError := tempConfig.Close()
	if closeFileError != nil {
		t.Fatalf("close file error %v", closeFileError)
	}
	config.PathToConfig = tempConfig.Name()
	cfg := &config.AppConfig{}
	client := &http.Client{
		Transport: &mockRoundTripper{
			StatusCode: 200,
		},
	}
	logCh := make(chan web.Response, 2)
	errCh := make(chan web.Response, 2)

	performRequests(cfg, client, logCh, errCh)
	time.Sleep(100 * time.Millisecond)

	if len(logCh) != 2 {
		t.Errorf("Expected 2 log responses, got %d", len(logCh))
	}
	if len(errCh) != 0 {
		t.Errorf("Expected 0 error responses, got %d", len(errCh))
	}
}

// Tests
func TestPerformRequests_ErrorResponse(t *testing.T) {
	tempConfig, err := os.CreateTemp("", "testConfig.toml")
	if err != nil {
		t.Fatalf("error creating temporary file: %v", err)
	}
	defer func(name string) {
		removeErr := os.Remove(name)
		if removeErr != nil {
			t.Fatal(removeErr.Error())
		}
	}(tempConfig.Name())
	_, writeConfigErr := tempConfig.Write([]byte(mockConfig))
	if writeConfigErr != nil {
		t.Fatalf("error writing to file: %v\n", writeConfigErr)
	}
	closeFileError := tempConfig.Close()
	if closeFileError != nil {
		t.Fatalf("close file error %v", closeFileError)
	}
	config.PathToConfig = tempConfig.Name()
	cfg := &config.AppConfig{}
	client := &http.Client{
		Transport: &mockRoundTripper{
			StatusCode: 500,
		},
	}
	logCh := make(chan web.Response, 2)
	errCh := make(chan web.Response, 2)

	performRequests(cfg, client, logCh, errCh)
	time.Sleep(100 * time.Millisecond)

	if len(logCh) != 2 {
		t.Errorf("Expected 2 log responses, got %d", len(logCh))
	}
	if len(errCh) != 2 {
		t.Errorf("Expected 2 error responses, got %d", len(errCh))
	}
}

func TestPerformRequests_UpdateError(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)
	config.PathToConfig = "tempConfig.Name()"
	cfg := &config.AppConfig{}
	client := &http.Client{
		Transport: &mockRoundTripper{
			StatusCode: 500,
		},
	}
	logCh := make(chan web.Response, 2)
	errCh := make(chan web.Response, 2)

	performRequests(cfg, client, logCh, errCh)
	time.Sleep(100 * time.Millisecond)

	if logBuf.String() == "" {
		t.Errorf("Expected error in log")
	}
}
