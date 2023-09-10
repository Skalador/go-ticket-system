package main

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PORT = ":8000"
const CONNECTIONSTRING = "mongodb+srv://kniederwanger:bat1OSclL7elzT0h@test-cluster.bwnvnol.mongodb.net/?retryWrites=true&w=majority"
//var client *mongo.Client

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

	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(CONNECTIONSTRING).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	
	if err != nil {
	  panic(err)
	}
	defer func() {
	  if err = client.Disconnect(context.TODO()); err != nil {
		panic(err)
	  }
	}()
	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	  panic(err)
	}
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

	//Todo: Create Database with GO
	//Create Collection with GO
	//Upload data.json with GO
	//Read database with GO
	

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //Fileserver to load css
	http.HandleFunc("/",rootHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/delete", deleteHandler)
    log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
