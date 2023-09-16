package db

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/Skalador/go-ticket-system/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitDB(CONNECTIONSTRING string) (*mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(CONNECTIONSTRING).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to the admin database and confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func findMaxID(ticketsCache *models.Tickets) int {
	if len(ticketsCache.Tickets) == 0 {
		return 0
	}

	max := ticketsCache.Tickets[0].ID
	for _, ticket := range ticketsCache.Tickets {
		if ticket.ID > max {
			max = ticket.ID
		}
	}

	return max
}

func AddTicketToCacheAndDB(req *http.Request, client *mongo.Client, ticketsCache *models.Tickets) {
	// Handle form submission
	subject := req.FormValue("subject")
	description := req.FormValue("description")
	id := findMaxID(ticketsCache) + 1
	severity := req.FormValue("severity")
	// Create a new ticket and add it to the list
	newTicket := models.Ticket{
		Subject:     subject,
		Description: description,
		ID:          id,
		Severity:    severity,
	}
	ticketsCache.Tickets = append(ticketsCache.Tickets, newTicket)
	insertOneTicket(client, newTicket)
	log.Println("Added ticket with id: ", id)
}

func obtainIDFromRequest(req *http.Request) int {
	idStr := req.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal("Error converting ID to integer:", err)
		panic(err)
	}
	return id
}

func DeleteTicketFromCacheAndDB(req *http.Request, client *mongo.Client, ticketsCache *models.Tickets) {
	id := obtainIDFromRequest(req)
	log.Println("Deleting ticket with id:", id)
	// Loop through the tickets and remove the one with the matching subject
	for i, ticket := range ticketsCache.Tickets {
		if ticket.ID == id {
			ticketsCache.Tickets = append(ticketsCache.Tickets[:i], ticketsCache.Tickets[i+1:]...)
			deleteOneTicket(client, ticket)
			break
		}
	}
}

func insertOneTicket(client *mongo.Client, ticket models.Ticket) {
	collection := client.Database("godb").Collection("tickets") //access collection
	ctx := context.TODO()
	result, err := collection.InsertOne(ctx, ticket)
	if err != nil {
		log.Fatal("Inserting one ticket failed! ", err.Error())
	}
	log.Printf("Inserted document with _id: %v \n", result.InsertedID)
}

func deleteOneTicket(client *mongo.Client, ticket models.Ticket) {
	collection := client.Database("godb").Collection("tickets") //access collection
	ctx := context.TODO()
	doc := bson.D{{"id", ticket.ID}}
	result, err := collection.DeleteOne(ctx, doc)
	if err != nil {
		log.Fatal("Deleting one ticket failed! ", err.Error())
	}
	log.Printf("Deleted %v document(s)\n", result.DeletedCount)
}

func insertManyTickets(client *mongo.Client, tickets []models.Ticket) {
	collection := client.Database("godb").Collection("tickets") //access collection
	docs := []interface{}{}
	// Iterate over the tickets and append each ticket to the docs slice
	for _, ticket := range tickets {
		docs = append(docs, ticket)
	}
	ctx := context.TODO()
	result, err := collection.InsertMany(ctx, docs)
	log.Printf("Documents inserted: %v\n", len(result.InsertedIDs))
	for _, id := range result.InsertedIDs {
		log.Printf("Inserted document with _id: %v\n", id)
	}
	if err != nil {
		log.Fatal("Inserting many tickets failed! ", err.Error())
	}
}

func PopulateDB(client *mongo.Client, tickets []models.Ticket) {
	//Read all available Databases
	filter := bson.D{} //Access all
	ctx := context.TODO()
	dbs, err := client.ListDatabaseNames(ctx, filter)
	godbExists := false
	if err != nil {
		panic(err)
	}
	for _, db := range dbs {
		log.Println("This database exists: ", db)
		if db == "godb" {
			log.Println("godb already exists, populating database not needed!")
			godbExists = true
		}
	}

	if !godbExists {
		log.Println("Populating godb database!")
		insertManyTickets(client, tickets)
	}

}

func ReadAllTickets(client *mongo.Client, ticketsCache *models.Tickets) {
	collection := client.Database("godb").Collection("tickets") //access collection
	filter := bson.D{}                                          //Access all
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error finding data in collection")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) { //iterate through the data
		var result models.Ticket
		if err := cursor.Decode(&result); err != nil {
			log.Println("Error decoding document:", err)
			continue //Continue to the next document on error
		}
		// Append the retrieved Ticket to the Tickets slice in ticketsCache
		ticketsCache.Tickets = append(ticketsCache.Tickets, result)
	}
}
