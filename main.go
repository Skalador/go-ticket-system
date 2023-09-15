package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Skalador/go-ticket-system/db"
	"github.com/Skalador/go-ticket-system/handlers"
	"github.com/Skalador/go-ticket-system/models"
)

const PORT = ":8000"

func main() {
	//create mongoDB client
	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	client, err := db.InitDB(connectionString)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	//Instantiate data and populate DB with sample data from data.json
	jsonFile, err := os.Open("data.json")
	if err != nil {
		log.Println(err)
	}
	log.Println("Sucessfully opened data.json!")
	defer jsonFile.Close()

	var ticketsCache models.Tickets
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ticketsCache)

	db.PopulateDB(client,ticketsCache.Tickets)
	ticketsCache=models.Tickets{} //clear cache to fill it with data from DB

	//Read data from Database in cache
	db.ReadAllTickets(client,&ticketsCache)

	//Start the web application
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //Fileserver to load css
	//Use wrapper functions for all handlers, to transfer data without global variables
	http.HandleFunc("/",handlers.TicketRootHandler(&ticketsCache))
	http.HandleFunc("/submit", handlers.TicketSubmitHandler(client,&ticketsCache))
	http.HandleFunc("/delete", handlers.TicketDeleteHandler(client,&ticketsCache))
    log.Printf("Listing for requests at http://localhost:%v/ \n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
