package telegram

import (
	"TelegramSiteMonitorBot/config"
	"TelegramSiteMonitorBot/web"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// Mock bot
type MockBot struct {
	Messages   []string
	ErrToThrow error
	Command    string
	Done       func()
}

func (m *MockBot) SendMessage(text string) error {
	m.Messages = append(m.Messages, text)
	if m.Done != nil {
		m.Done()
	}
	return m.ErrToThrow
}

func (m *MockBot) Updates() <-chan Update {
	ch := make(chan Update)
	go func() {
		ch <- Update{Command: m.Command}
		close(ch)
	}()
	return ch
}

// Helper function for tests
func newTestConfig(URLs []string) *config.AppConfig {
	cfg := config.AppConfig{}
	val := reflect.ValueOf(&cfg.Sites).Elem()
	urlsField := val.FieldByName("URLs")
	if urlsField.CanSet() {
		urlsField.Set(reflect.ValueOf(URLs))
	} else {
		panic("URLs field is not editable")
	}
	return &cfg
}

// Tests
func TestHandleWebResponses_MessageInChan(t *testing.T) {
	responsesCh := make(chan web.Response, 1)
	defer close(responsesCh)
	var wg sync.WaitGroup
	wg.Add(1)
	bot := &MockBot{
		Done: func() {
			wg.Done()
		},
	}

	go handleWebResponses(bot, responsesCh)
	responsesCh <- web.Response{StatusCode: 200}
	wg.Wait()

	if len(bot.Messages) != 1 {
		t.Errorf("expected 1 message, received %d", len(bot.Messages))
	}
}

func TestHandleWebResponses_Error(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)
	responsesCh := make(chan web.Response, 1)
	defer close(responsesCh)
	var wg sync.WaitGroup
	bot := &MockBot{
		ErrToThrow: fmt.Errorf("error"),
		Done: func() {
			wg.Done()
		},
	}
	wg.Add(1)

	go handleWebResponses(bot, responsesCh)
	responsesCh <- web.Response{StatusCode: 200}

	if strings.Contains(logBuf.String(), "error") {
		t.Errorf("expected error in log")
	}
}

func TestHandleCommands_Status_Response(t *testing.T) {
	var wg sync.WaitGroup
	bot := &MockBot{
		Done: func() {
			wg.Done()
			return
		},
		Command: "status",
	}
	cfg := newTestConfig([]string{"https://example.com", "https://google.com"})
	client := http.Client{Timeout: 5 * time.Second}

	wg.Add(1)
	handleCommands(bot, cfg, &client)
	wg.Wait()

	if len(bot.Messages) != 1 && !strings.Contains(bot.Messages[0], "https://example.com") {
		t.Error("expected status message")
	}
}

func TestHandleCommands_Nil_Nothing(t *testing.T) {
	bot := &MockBot{
		Done: func() {
			return
		},
		Command: "",
	}
	cfg := newTestConfig([]string{"https://example.com", "https://google.com"})
	client := http.Client{Timeout: 5 * time.Second}

	handleCommands(bot, cfg, &client)

	if len(bot.Messages) != 0 {
		t.Error("expected 0 messages")
	}
}

func TestHandleCommands_UnknownCommand_Response(t *testing.T) {
	var wg sync.WaitGroup
	bot := &MockBot{
		Done: func() {
			wg.Done()
			return
		},
		Command: "greeting",
	}
	cfg := newTestConfig([]string{"https://example.com", "https://google.com"})
	client := http.Client{Timeout: 5 * time.Second}

	wg.Add(1)
	handleCommands(bot, cfg, &client)
	wg.Wait()

	if len(bot.Messages) != 1 && bot.Messages[0] != "Неизвестная команда" {
		t.Error("expected unknown command message")
	}
}

func TestHandleCommands_Status(t *testing.T) {
	updates := make(chan Update, 1)
	updates <- Update{Command: "status"}
	close(updates)
}

func TestHandleStatusCommand_MessageInCh(t *testing.T) {
	bot := &MockBot{
		Done: func() {
			return
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	URLs := []string{server.URL}
	client := http.Client{Timeout: 5 * time.Second}

	handleStatusCommand(bot, URLs, &client)

	if len(bot.Messages) != 1 {
		t.Errorf("expected 1 message, received %d", len(bot.Messages))
	}
}

func TestHandleStatusCommand_Error(t *testing.T) {
	bot := &MockBot{
		ErrToThrow: fmt.Errorf("error"),
		Done: func() {
			return
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	URLs := []string{server.URL}
	client := http.Client{Timeout: 5 * time.Second}

	handleStatusCommand(bot, URLs, &client)

	if len(bot.Messages) != 1 {
		t.Errorf("expected 1 message, received %d", len(bot.Messages))
	}
}

func TestHandleUnknownCommand_Message(t *testing.T) {
	bot := &MockBot{
		Done: func() {
			return
		},
	}

	handleUnknownCommand(bot)

	if len(bot.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(bot.Messages))
	}
	if bot.Messages[0] != "Неизвестная команда" {
		t.Errorf("unexpected message: %s", bot.Messages[0])
	}
}

func TestHandleUnknownCommand_Error(t *testing.T) {
	var logBuf bytes.Buffer
	log.SetOutput(&logBuf)
	defer log.SetOutput(os.Stderr)
	bot := &MockBot{}
	bot.ErrToThrow = fmt.Errorf("error")

	handleUnknownCommand(bot)

	if logBuf.String() == "" {
		t.Error("expected error in log")
	}
}
