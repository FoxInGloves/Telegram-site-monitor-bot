package web

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestGetRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	chLog := make(chan Response, 1)
	chErrors := make(chan Response, 1)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	GetRequest(server.URL, client, chLog, chErrors)

	select {
	case resp := <-chLog:
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected а responce code %d, received %d", http.StatusOK, resp.StatusCode)
		}
	default:
		t.Error("expected a response in chLog")
	}

	select {
	case <-chErrors:
		t.Error("did not expect anything in chErrors")
	default:
		// ок
	}
}

func TestGetRequest_BadURL_MsgInErrCh(t *testing.T) {
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	badURL := "http://localhost:12345"

	chLog := make(chan Response, 1)
	chErrors := make(chan Response, 1)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	GetRequest(badURL, client, chLog, chErrors)

	select {
	case resp := <-chErrors:
		if resp.StatusCode != 0 {
			t.Errorf("expected status code 0 for error, got %d", resp.StatusCode)
		}
	default:
		t.Error("expected an error response in chErrors")
	}

	select {
	case <-chLog:
		t.Error("did not expect anything in chLog")
	default:
		// ок
	}
}

func TestPostRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	chLog := make(chan Response, 1)
	chErrors := make(chan Response, 1)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	GetRequest(server.URL, client, chLog, chErrors)

	select {
	case resp := <-chErrors:
		if resp.StatusCode != 500 {
			t.Errorf("expected status code 500 for error, got %d", resp.StatusCode)
		}
	default:
		t.Error("expected an error response in chErrors")
	}
}
