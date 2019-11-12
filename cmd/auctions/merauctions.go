package main

import (
	"flag"
	"fmt"

	"github.com/MerAuctions/MerAuctions/server"
)

func main() {

	dbURL := flag.String("mongodb-url", "mongodb://localhost:27017", "URL to connet to mongodb database")
	dbName := flag.String("database", "testing", "Database name in mongodb")
	flag.Parse()
	router := server.CreateRouter()
	fmt.Println(*dbName)
	server.ConnectToDB(*dbURL, *dbName)
	server.StartServer(router)
}
