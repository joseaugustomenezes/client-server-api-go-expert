package main

import (
	"context"
	"encoding/json"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Cotacao Cotacao `json:"USDBRL"`
}

type Cotacao struct {
	ID  int    `gorm:"primaryKey" json:"-"`
	Bid string `json:"bid"`
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*200)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("Erro ao criar requisição")
		http.Error(w, "Erro ao criar requisição", http.StatusInternalServerError)
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Println("Erro ao realizar requisição")
		http.Error(w, "Erro ao realizar requisição", http.StatusInternalServerError)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Erro ao ler o body da requisição")
		http.Error(w, "Erro ao ler o body da requisição", http.StatusInternalServerError)
	}
	var data Data
	json.Unmarshal(body, &data)

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Println("Erro ao abrir conexão com o banco de dados")
		http.Error(w, "Erro ao abrir conexão com o banco de dados", http.StatusInternalServerError)
	}
	err = db.AutoMigrate(&Cotacao{})
	if err != nil {
		log.Println("Erro ao criar table Cotacao")
		http.Error(w, "Erro ao criar table Cotacao", http.StatusInternalServerError)
	}

	bdCtx, bdCancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer bdCancel()
	err = db.WithContext(bdCtx).Create(&data.Cotacao).Error
	if err != nil {
		log.Println("Erro ao adicionar cotação ao bd")
		http.Error(w, "Erro ao adicionar cotação ao bd", http.StatusInternalServerError)
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseBody, _ := json.Marshal(data.Cotacao)
	w.Write(responseBody)
	log.Println("Requisição realizada com sucesso!")
}
