package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

const (
	apiURL     = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	apiTimeout = 200 * time.Millisecond
	dbTimeout  = 10 * time.Millisecond
	serverPort = ":8080"
)

type USDBRL struct {
	Bid string `json:"bid"`
}

type APIResponse struct {
	USDBRL USDBRL `json:"USDBRL"`
}

func main() {
	db, err := sql.Open("sqlite", "cotacoes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.SetMaxOpenConns(1)
	createTable(db)

	http.HandleFunc("/cotacao", cotacaoHandler(db))

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(serverPort, nil))
}

func createTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		data TIMESTAMP
	)`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

func cotacaoHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctxAPI, cancelAPI := context.WithTimeout(r.Context(), apiTimeout)
		defer cancelAPI()

		req, err := http.NewRequestWithContext(ctxAPI, http.MethodGet, apiURL, nil)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("Erro na API:", err)
			http.Error(w, err.Error(), 504)
			return
		}
		defer resp.Body.Close()

		var data APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		bid := data.USDBRL.Bid

		ctxDB, cancelDB := context.WithTimeout(r.Context(), dbTimeout)
		defer cancelDB()

		_, err = db.ExecContext(ctxDB,
			"INSERT INTO cotacoes (bid, data) VALUES (?, ?)",
			bid, time.Now(),
		)
		if err != nil {
			log.Println("Erro no banco:", err)
		}

		json.NewEncoder(w).Encode(map[string]string{
			"bid": bid,
		})
	}
}
