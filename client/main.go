package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro ao criar requisição")
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Erro ao realizar requisição")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erro ao ler body da requisição")
	}
	var cotacao Cotacao
	json.Unmarshal(body, &cotacao)

	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal("Erro ao criar arquivo cotacao.txt")
	}
	_, err = f.Write([]byte("Dólar: " + cotacao.Bid))
	if err != nil {
		log.Fatal("Erro ao escrever no arquivo cotacao.txt")
	}
	defer f.Close()
	log.Println("cotacao.txt criado com sucesso!")
}
