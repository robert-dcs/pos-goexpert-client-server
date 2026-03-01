package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

/* =========================
   MOCKS
========================= */

type successRoundTripper struct{}

func (m successRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body := `{"USDBRL":{"bid":"5.00"}}`
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}, nil
}

type errorRoundTripper struct{}

func (m errorRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("api error")
}

type invalidJSONRoundTripper struct{}

func (m invalidJSONRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	body := `invalid-json`
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}, nil
}

/* =========================
   HELPERS
========================= */

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	createTable(db)
	return db
}

/* =========================
   TESTES
========================= */

func TestCotacaoHandlerSuccess(t *testing.T) {
	db := setupTestDB(t)

	dbTimeout = time.Second

	httpClient = &http.Client{
		Transport: successRoundTripper{},
		Timeout:   time.Second,
	}

	req := httptest.NewRequest("GET", "/cotacao", nil)
	w := httptest.NewRecorder()

	handler := cotacaoHandler(db)
	handler(w, req)

	res := w.Result()

	if res.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(res.Body).Decode(&body)

	if body["bid"] != "5.00" {
		t.Fatalf("expected bid 5.00")
	}
}

func TestCotacaoHandlerAPIError(t *testing.T) {
	db := setupTestDB(t)

	httpClient = &http.Client{
		Transport: errorRoundTripper{},
	}

	req := httptest.NewRequest("GET", "/cotacao", nil)
	w := httptest.NewRecorder()

	handler := cotacaoHandler(db)
	handler(w, req)

	if w.Result().StatusCode != 504 {
		t.Fatalf("expected 504 when API fails")
	}
}

func TestCotacaoHandlerInvalidJSON(t *testing.T) {
	db := setupTestDB(t)

	httpClient = &http.Client{
		Transport: invalidJSONRoundTripper{},
	}

	req := httptest.NewRequest("GET", "/cotacao", nil)
	w := httptest.NewRecorder()

	handler := cotacaoHandler(db)
	handler(w, req)

	if w.Result().StatusCode != 500 {
		t.Fatalf("expected 500 when JSON invalid")
	}
}

func TestCotacaoHandlerDBTimeout(t *testing.T) {
	db := setupTestDB(t)

	httpClient = &http.Client{
		Transport: successRoundTripper{},
	}

	// força timeout do banco
	dbTimeout = 1 * time.Nanosecond

	req := httptest.NewRequest("GET", "/cotacao", nil)
	w := httptest.NewRecorder()

	handler := cotacaoHandler(db)
	handler(w, req)

	// Mesmo com erro no banco, resposta ainda deve ser 200
	if w.Result().StatusCode != 200 {
		t.Fatalf("expected 200 even if DB fails")
	}
}
