# GO-Ticket-System

In this repostiory a `GO` web application is built which has the basic functionality of a ticketing system.

## Architecture
The ticket system utilizes the `net/http` package and a templating engine `html/template` for front- and backend interaction. The system uses a write-through cache model for the tickets with the `ticketsCache` variable. 

The backend database is `MongoDB`. MongoDB is hosted for free on Atlas. The Data can be visualized with `MongoDBCompass`.
On first startup the `MongoDB` will be checked for a database called `godb`. If this database does not exist, then the database will be created and populated with the data from `data.json`.

## Demo
![](https://github.com/Skalador/go-ticket-system/demo.mp4)

## Prerequisites
Installing the mongoDB driver:
```
go get go.mongodb.org/mongo-driver/mongo
```

## Execute the code
An environment variable `MONGODB_CONNECTION_STRING` is used for the database connectivity, thus the connection string is not exposed in the code itself.

Expose the variable:
```
Windows: $env:MONGODB_CONNECTION_STRING = 'mongodb+srv://username:password@database/'
Linux: export MONGODB_CONNECTION_STRING="mongodb+srv://username:password@database/"
```

Run the code:
``
go run main.go
```

## Known issues
- There is no ticket cache timeout, thus a direct interaction with the database, e.g. deleting a document, will desync the programm.
- The environment variable approach is only intended for local testing
- Cross site scripting (XSS) is possible
