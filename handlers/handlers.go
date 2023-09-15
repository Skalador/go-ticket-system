package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/Skalador/go-ticket-system/db"
	"github.com/Skalador/go-ticket-system/models"
	"go.mongodb.org/mongo-driver/mongo"
)

//Wrapper function for rootHandler
func TicketRootHandler(ticketsCache *models.Tickets) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request)  {
		t,err := template.ParseFiles("index.html")
		log.Println("Received request:", req)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Send response:", w)
		t.Execute(w,ticketsCache.Tickets)
	}
}


//wrapper function for delete handler
func TicketDeleteHandler(client *mongo.Client,ticketsCache *models.Tickets) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request)  {
		log.Println("Received request:", req)
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		db.DeleteTicketFromCacheAndDB(req, client, ticketsCache)

		// Redirect back to the main page
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}


//wrapper function for submit handler
func TicketSubmitHandler(client *mongo.Client,ticketsCache *models.Tickets) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request)  {
		log.Println("Received request:", req)
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		db.AddTicketToCacheAndDB(req, client, ticketsCache)

		// Redirect back to the main page
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}
