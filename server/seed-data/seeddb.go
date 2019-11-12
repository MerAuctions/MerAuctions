package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/MerAuctions/MerAuctions/server"
)

func InsertAuctionsToDB() *models.AuctionList {
	var auc models.AuctionList
	file, err := ioutil.ReadFile("./auctions.json")
	if err != nil {
		log.Fatal("Error reading auctions.json : ", err.Error())
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.InsertAuctions(&auc)
	if err != nil {
		log.Fatal("Error populating auctions.json : ", err.Error())
	}
	return &auc
}

func RemoveAuctionsFromDB() {
	var auc models.AuctionList
	file, err := ioutil.ReadFile("./auctions.json")
	if err != nil {
		log.Fatal("Error reading auctions.json : ", err.Error())
	}
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.DeleteAuctions(&auc)
	if err != nil {
		log.Fatal("Error deleting auctions.json : ", err.Error())
	}
}

func InsertBidsToDB() *models.BidList {
	var bids models.BidList
	file, err := ioutil.ReadFile("./bids.json")
	if err != nil {
		log.Fatal("Error reading bids.json : ", err.Error())
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.InsertBids(&bids)
	if err != nil {
		log.Fatal("Error populating bids.json : ", err.Error())
	}
	return &bids
}

func RemoveBidsFromDB() {
	var bids models.BidList
	file, err := ioutil.ReadFile("./bids.json")
	if err != nil {
		log.Fatal("Error reading bids.json : ", err.Error())
	}
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.DeleteBids(&bids)
	if err != nil {
		log.Fatal("Error deleting bids.json : ", err.Error())
	}
}

func main() {
	dbURL := "mongodb://localhost:27017"
	dbName := "testing"

	server.ConnectToDB(dbURL, dbName)
	InsertAuctionsToDB()
	InsertBidsToDB()
}
