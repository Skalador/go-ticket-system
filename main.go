package main

import (
	"log"
	"net/http"
	"text/template"
)

const PORT = ":8000"

type Ticket struct {
	Subject string
	Description string
//	ID string
//	Severity int
}

//create a slice of multiple tickets to display
var tickets = []Ticket{ // Store submitted tickets in a global slice
	{
		Subject: "Missing ID",
		Description: "IDs should be added to tickets",
	},
	{
		Subject: "Missing Severity",
		Description: "Severity should be added to tickets",
	},
	{
		Subject: "Missing Database interaction",
		Description: "Database integration should be added",
	},
	{
		Subject: "Containerize Applications",
		Description: "Bring the entire structure in microservice architecture",
	},
} 


func addTicket(req *http.Request) {
        // Handle form submission
        subject := req.FormValue("subject")
        description := req.FormValue("description")
        // Create a new ticket and add it to the list
        newTicket := Ticket{
            Subject:     subject,
            Description: description,
        }
        tickets = append(tickets, newTicket)
		log.Println("Added ticket with subject: ", subject)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	t,err := template.ParseFiles("index.html")
	log.Println("Received request:", req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Send response:", w)
	t.Execute(w,tickets)
}

func deleteHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request:", req)
    if req.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    subject := req.FormValue("subject")
	log.Println("Deleting ticket with subject:", subject)
    // Loop through the tickets and remove the one with the matching subject
    for i, ticket := range tickets {
        if ticket.Subject == subject {
            tickets = append(tickets[:i], tickets[i+1:]...)
            break
        }
    }

    // Redirect back to the main page
    http.Redirect(w, req, "/", http.StatusSeeOther)
}

func submitHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Received request:", req)
	if req.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

	addTicket(req)

	// Redirect back to the main page
    http.Redirect(w, req, "/", http.StatusSeeOther)
}

func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //Fileserver to load css
	http.HandleFunc("/",rootHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/delete", deleteHandler)
    log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
