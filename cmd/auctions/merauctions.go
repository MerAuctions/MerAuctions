package main

import (
	"flag"

	"github.com/MerAuctions/MerAuctions/server"
)

func main() {

	dbURL := flag.String("mongodb-url", "mongodb://localhost:27017", "URL to connet to mongodb database")
	dbName := flag.String("database", "testing", "Database name in mongodb")
	flag.Parse()
	router := server.CreateRouter()
	server.ConnectToDB(*dbURL, *dbName)
	server.StartServer(router)
}
