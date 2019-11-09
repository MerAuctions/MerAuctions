package data

import (
	"context"
	"fmt"
	"log"

	//"go.mongodb.org/mongo-driver/bson"
	"github.com/MerAuctions/MerAuctions/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Connect to mongoDB of given URL
func connectDB(URL string) *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI(URL)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return client
}

func getData(Database string, Collection string) {
	//collection := client.Database(Database).Collection(Collection)
}

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
