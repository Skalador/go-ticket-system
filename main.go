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



func addTicketToCacheAndDB(req *http.Request, client *mongo.Client, ticketsCache *Tickets) {
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
		ticketsCache.Tickets = append(ticketsCache.Tickets, newTicket)
		insertOneTicket(client,newTicket)
		log.Println("Added ticket with id: ", id)
}

func deleteTicketFromCacheAndDB (req *http.Request, client *mongo.Client, ticketsCache *Tickets) {
	id := req.FormValue("id")
	log.Println("Deleting ticket with id:", id)
	// Loop through the tickets and remove the one with the matching subject
	for i, ticket := range ticketsCache.Tickets {
		if ticket.ID == id {
			ticketsCache.Tickets = append(ticketsCache.Tickets[:i], ticketsCache.Tickets[i+1:]...)
			deleteOneTicket(client,ticket)
			break
		}
	}
}

//Wrapper function for rootHandler
func ticketRootHandler(ticketsCache *Tickets) func(http.ResponseWriter, *http.Request) {
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
func ticketDeleteHandler(client *mongo.Client,ticketsCache *Tickets) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request)  {
		log.Println("Received request:", req)
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		deleteTicketFromCacheAndDB(req,client,ticketsCache)

		// Redirect back to the main page
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}


//wrapper function for submit handler
func ticketSubmitHandler(client *mongo.Client,ticketsCache *Tickets) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request)  {
		log.Println("Received request:", req)
		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		addTicketToCacheAndDB(req, client, ticketsCache)

		// Redirect back to the main page
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}
}


func insertOneTicket(client *mongo.Client,ticket Ticket) {
	collection :=client.Database("godb").Collection("tickets") //access collection
	ctx := context.TODO()
	result,err := collection.InsertOne(ctx,ticket)
	log.Printf("Inserted document with _id: %v \n", result.InsertedID)
	if err != nil {
		panic(err)
	}
}

func deleteOneTicket(client *mongo.Client,ticket Ticket) {
	collection :=client.Database("godb").Collection("tickets") //access collection
	ctx := context.TODO()
	doc:=bson.D{{"id", ticket.ID}}
	result, err := collection.DeleteOne(ctx, doc)
	if err != nil {
		panic(err)
	}
	log.Printf("Deleted %v document(s)\n", result.DeletedCount)
}

func insertManyTickets(client *mongo.Client,tickets []Ticket) {
	collection :=client.Database("godb").Collection("tickets") //access collection
	docs :=[]interface{}{} 
	// Iterate over the tickets and append each ticket to the docs slice
	for _, ticket := range tickets {
		docs = append(docs, ticket)
	}
	ctx := context.TODO()
	result,err := collection.InsertMany(ctx,docs)
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
	ctx := context.TODO()
	dbs, err := client.ListDatabaseNames(ctx,filter)
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

	for cursor.Next(ctx) { //iterate through the data
		var result Ticket
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding document:", err)
			continue //Continue to the next document on error
		}
		// Append the retrieved Ticket to the Tickets slice in ticketsCache
        ticketsCache.Tickets = append(ticketsCache.Tickets, result)
	}
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

	//Todo: Remove ID field and handle it in the background


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) //Fileserver to load css
	//Use wrapper functions for all handlers, to transfer data without global variables
	http.HandleFunc("/",ticketRootHandler(&ticketsCache))
	http.HandleFunc("/submit", ticketSubmitHandler(client,&ticketsCache))
	http.HandleFunc("/delete", ticketDeleteHandler(client,&ticketsCache))
    log.Println("Listing for requests at http://localhost:8000/")
	log.Fatal(http.ListenAndServe(PORT, nil))
}
