package data

import (
	"context"
	"fmt"
	"log"
	"time"
)

<<<<<<< HEAD

func GetAllAuctionsFromDB(*models.AuctionList) {

}

func GetAuctionByIdFromDB(auc *models.Auction, id string) {
=======
type DBClient struct {
	dbName string
	client *mongo.Client
}

//Connect to mongoDB of given URL
func ConnectDB(URL string, dbName string) *DBClient {
	// Set client options
	clientOptions := options.Client().ApplyURI(URL)
>>>>>>> a009c765f83a58230ad3ae2351990f867c264efc

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

<<<<<<< HEAD
=======
	fmt.Println("Connected to MongoDB!")
	return &DBClient{
		client: client,
		dbName: dbName,
	}
>>>>>>> a009c765f83a58230ad3ae2351990f867c264efc
}

func (c *DBClient) TestAddData() error {
	collection := c.client.Database(c.dbName).Collection("numbers")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//inserting data
	res, err := collection.InsertOne(ctx, models.User{
		UserID:   "test1",
		UserName: "harsh",
	})
	if err != nil {
		return err
	}
	log.Printf("Added document with ID: %s", res.InsertedID)
	return nil
}

func (c *DBClient) TestGetData() (*models.User, error) {
	collection := c.client.Database(c.dbName).Collection("numbers")
	var result models.User
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	//querying data
	err := collection.FindOne(ctx, bson.M{"username": "harsh"}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
