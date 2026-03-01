package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

/* ======================
   SUCESSO
====================== */

func TestClientSuccess(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"bid":"5.50"}`))
	}))
	defer mockServer.Close()

	serverURL = mockServer.URL
	httpClient = http.DefaultClient

	run()

	data, err := os.ReadFile("cotacao.txt")
	if err != nil {
		t.Fatal("file not created")
	}

	if string(data) != "Dólar: 5.50" {
		t.Fatal("unexpected file content")
	}

	os.Remove("cotacao.txt")
}

/* ======================
   ERRO HTTP
====================== */

type errorRoundTripper struct{}

func (e errorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("connection error")
}

func TestClientHTTPError(t *testing.T) {
	httpClient = &http.Client{
		Transport: errorRoundTripper{},
		Timeout:   time.Second,
	}

	run()

	_, err := os.Stat("cotacao.txt")
	if !os.IsNotExist(err) {
		t.Fatal("file should not be created on HTTP error")
	}
}

/* ======================
   JSON INVÁLIDO
====================== */

func TestClientInvalidJSON(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`invalid-json`))
	}))
	defer mockServer.Close()

	serverURL = mockServer.URL
	httpClient = http.DefaultClient

	run()

	_, err := os.Stat("cotacao.txt")
	if !os.IsNotExist(err) {
		t.Fatal("file should not be created on invalid JSON")
	}
}
