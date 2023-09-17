# GO-Ticket-System

In this repostiory a `GO` web application is built which has the basic functionality of a ticketing system.

## Architecture
The ticket system utilizes the `net/http` package and a templating engine `html/template` for front- and backend interaction. The system uses a write-through cache model for the tickets with the `ticketsCache` variable. 

The backend database is `MongoDB`. MongoDB is hosted for free on Atlas. The Data can be visualized with `MongoDBCompass`.
On first startup the `MongoDB` will be checked for a database called `godb`. If this database does not exist, then the database will be created and populated with the data from `data.json`.

## Demo
https://github.com/Skalador/go-ticket-system/assets/117681263/87319c17-65a8-4ea1-a464-2d8f3b43c779




## Prerequisites
Installing the mongoDB driver:
```
go get go.mongodb.org/mongo-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/integration/mtest
go get github.com/stretchr/testify/assert
```

## Execute the code
An environment variable `MONGODB_CONNECTION_STRING` is used for the database connectivity, thus the connection string is not exposed in the code itself.

Expose the variable:
```
Windows: $env:MONGODB_CONNECTION_STRING = 'mongodb+srv://username:password@database/'
Linux: export MONGODB_CONNECTION_STRING="mongodb+srv://username:password@database/"
```

Run the code:
```
go run main.go
```

## Run tests

Running all tests:
```
go test ./...
```

Testing specific packages:
```
go test ./handlers
```

## Known issues
- There is no ticket cache timeout, thus a direct interaction with the database, e.g. deleting a document, will desync the programm.
- The environment variable approach is only intended for local testing
- Cross site scripting (XSS) is possible
