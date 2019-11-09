package data

import (
	"context"
	"fmt"
	"log"
	"time"
)


func GetAllAuctionsFromDB(*models.AuctionList) {

}

func GetAuctionByIdFromDB(auc *models.Auction, id string) {

}

func GetTopFiveBidsFromDB(top5bids *[5]models.Bid, id string) {

}

func AddNewUserToDB(usr *models.User) int {

	return 0
}

func AddNewBid(bid *models.Bid) int {

	return 0
}

func GetResult(res *models.Result, id string) {

}

func (c *DBClient) TestData() (*models.User, error) {
	collection := c.client.Database("testing").Collection("numbers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//inserting data
	res, err := collection.InsertOne(ctx, models.User{
		UserID:   "test1",
		UserName: "harsh",
	})
	if err != nil {
		return nil, err
	}
	log.Printf("Added document with ID: %s", res.InsertedID)

	var result models.User
	ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
	//querying data
	err = collection.FindOne(ctx, bson.M{"username": "harsh"}).Decode(&result)

	if err != nil {
		return nil, err
	}
	return &result, nil
}
