package main

import (
	"flag"

	"github.com/MerAuctions/MerAuctions/server"
)

func main() {

	dbURL := flag.String("mongodb-url", "mongodb://localhost:27017", "URL to connet to mongodb database")
	flag.Parse()
	router := server.CreateRouter()
	server.ConnectToDB(*dbURL)
	server.StartServer(router)
}
