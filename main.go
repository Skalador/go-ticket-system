package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const PORT = ":8000"
const CONNECTIONSTRING = "mongodb+srv://kniederwanger:bat1OSclL7elzT0h@go-tickets.z7ats48.mongodb.net/"

type Tickets struct {
	Tickets []Ticket `json:"Tickets"`
}

type Ticket struct {
	Subject string `json:"Subject"`
	Description string `json:"Description"`
	ID string `json:"ID"`
	Severity string `json:"Severity"`
}

//create a slice of multiple tickets to display
var tickets = []Ticket{ // Store submitted tickets in a global slice
	{
		Subject: "Missing ID",
		Description: "IDs should be added to tickets",
		ID: "00001",
		Severity: "SEV4",
	},
	{
		Subject: "Missing Severity",
		Description: "Severity should be added to tickets",
		ID: "00002",
		Severity: "SEV4",
	},
	{
		Subject: "Missing Database interaction",
		Description: "Database integration should be added",
		ID: "00003",
		Severity: "SEV4",
	},
	{
		Subject: "Containerize Applications",
		Description: "Bring the entire structure in microservice architecture",
		ID: "00004",
		Severity: "SEV3",
	},
} 


func addTicket(req *http.Request) {
        // Handle form submission
        subject := req.FormValue("subject")
        description := req.FormValue("description")
		id := req.FormValue("id")
		severity := req.FormValue("severity")
        // Create a new ticket and add it to the list
        newTicket := Ticket{
            Subject:     subject,
            Description: description,
			ID: id,
			Severity: severity,
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

func insertOneTicket(client *mongo.Client,ticket Ticket) {
	collection :=client.Database("godb").Collection("tickets") //access collection
	result,err := collection.InsertOne(context.TODO(),tickets[0])
	log.Printf("Inserted document with _id: %v \n", result.InsertedID)
	if err != nil {
		panic(err)
	}
}

func insertManyTickets(client *mongo.Client,tickets []Ticket) {
	collection :=client.Database("godb").Collection("tickets") //access collection
	docs :=[]interface{}{tickets[0],tickets[1],tickets[2],tickets[3]}
	result,err := collection.InsertMany(context.TODO(),docs)
	log.Printf("Documents inserted: %v\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
    	log.Printf("Inserted document with _id: %v\n", id)
	}
	if err != nil {
		panic(err)
	}
}

func populateDB(client *mongo.Client, tickets []Ticket){
	//Read all available Databases
	filter :=bson.D{} //Access all
	dbs, err := client.ListDatabaseNames(context.TODO(),filter)
	godbExists := false
	if err != nil {
		panic(err)
	}
	for _, db := range dbs {
		log.Println("This database exists: ",db)
		if db == "godb" {
			log.Println("godb already exists, populating database not needed!")
			godbExists = true
		} 
	}

	if !godbExists {
		log.Println("Populating godb database!")
		insertManyTickets(client,tickets)
	}

}

func readAllTickets(client *mongo.Client, ticketsCache *Tickets) {
	collection :=client.Database("godb").Collection("tickets") //access collection
	filter :=bson.D{} //Access all
	ctx := context.TODO()

	cursor, err := collection.Find(ctx,filter)
	if err != nil {
		log.Fatal("Error finding data in collection")
		panic(err)
	}
	defer cursor.Close(ctx)
//	var data []bson.M
//
//	if err = cursor.All(ctx,&data); err != nil {
//		log.Fatal("Error parsing data from database")
//		panic(err)
//	}
//	log.Println(data)
	for cursor.Next(ctx) {
		var result Ticket
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding document:", err)
			continue //Continue to the next document on error
		}
		// Append the retrieved Ticket to the Tickets slice in ticketsCache
        ticketsCache.Tickets = append(ticketsCache.Tickets, result)
	}

//	for _,tt := range ticket {
//		log.Println(cursor.Decode(&tt))
//	}
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
	// Send a ping to the admin database and confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	  panic(err)
	}
	log.Println("Pinged your deployment, i.e. admin database. You successfully connected to MongoDB!")


	//Instantiate data and populate DB with sample data
	jsonFile, err := os.Open("data.json")
	if err != nil {
		log.Println(err)
	}
	log.Println("Sucessfully opened data.json!")
	defer jsonFile.Close()

	var ticketsCache Tickets
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &ticketsCache)
	populateDB(client,ticketsCache.Tickets)
	ticketsCache=Tickets{} //clear cache to fill it with data from DB

	//Read data from Database in cache
	readAllTickets(client,&ticketsCache)
	log.Println(ticketsCache.Tickets)

	//Todo: Remove ID field and handle it in the background


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //Fileserver to load css
	http.HandleFunc("/",rootHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/delete", deleteHandler)
    log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
