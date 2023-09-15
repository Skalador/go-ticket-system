package models

type Tickets struct {
	Tickets []Ticket `json:"Tickets"`
}

type Ticket struct {
	Subject     string `json:"Subject"`
	Description string `json:"Description"`
	ID          int    `json:"ID"`
	Severity    string `json:"Severity"`
}