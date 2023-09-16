package models

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Tickets struct {
	Tickets []Ticket `json:"Tickets"`
}

type Ticket struct {
	Subject     string `json:"Subject"`
	Description string `json:"Description"`
	ID          int    `json:"ID"`
	Severity    string `json:"Severity"`
}

func InitTicketsCache(ticketsCache *Tickets) {
	jsonFile, err := os.Open("../data.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ticketsCache)
}
