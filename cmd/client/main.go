package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	serverURL     = "http://localhost:8080/cotacao"
	clientTimeout = 300 * time.Millisecond
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), clientTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL, nil)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro ao chamar server:", err)
		return
	}
	defer resp.Body.Close()

	var data CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
		return
	}

	content := "Dólar: " + data.Bid

	if err := os.WriteFile("cotacao.txt", []byte(content), 0644); err != nil {
		log.Println(err)
		return
	}

	log.Println("Cotação salva com sucesso:", content)
}
