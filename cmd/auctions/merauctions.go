package main

import (
	"github.com/MerAuctions/MerAuctions/server"
)

func main() {
	router := server.CreateRouter()
	server.StartServer(router)
}
