package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/MerAuctions/MerAuctions/server"
)

func InsertAuctionsToDB() {
	var auc []models.Auction
	file, err := ioutil.ReadFile("./auctions.json")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.InsertAuctions(&auc)
	if err != nil {
		return
	}
}

func RemoveAuctionsFromDB() {
	var auc []models.Auction
	file, err := ioutil.ReadFile("./auctions.json")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.DeleteAuctions(&auc)
	if err != nil {
		return
	}
}

func InsertBidsToDB() {
	var bids []models.Bid
	file, err := ioutil.ReadFile("./bids.json")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.InsertBids(&bids)
	if err != nil {
		return
	}
}

func RemoveBidsFromDB() {
	var bids []models.Bid
	file, err := ioutil.ReadFile("./bids.json")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.DeleteBids(&bids)
	if err != nil {
		return
	}
}

func main() {
	dbURL := "mongodb://localhost:27017"
	dbName := "testing-main"

	server.ConnectToDB(dbURL, dbName)
	InsertAuctionsToDB()
	InsertBidsToDB()
}
