package data

import (
	"context"
	"fmt"
	"log"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	"github.com/MerAuctions/MerAuctions/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	client *mongo.Client
}

//Connect to mongoDB of given URL
func ConnectDB(URL string) *DBClient {
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
	return &DBClient{
		client: client,
	}
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
