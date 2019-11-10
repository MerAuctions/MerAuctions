package db

import (
	"context"
	"fmt"
	"log"

	"github.com/MerAuctions/MerAuctions/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
  URL string
	client *mongo.Client
}

//Connect to mongoDB of given URL
func (c *DBClient)ConnectDB(url string) *DBClient {
	// Set client options
	clientOptions := options.Client().ApplyURI(url)

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
    URL: url,
    client: client,
  }
}

//Insert a user in db
func (c *DBClient)InsertUser(usr *models.User) error{
  collection := c.client.Database("production").Collection("users")
  insertResult, err := collection.InsertOne(context.TODO(), usr)
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println("Inserted user: ", insertResult.InsertedID)
  return err
}

//Insert a bid in db
func (c *DBClient)InsertBid(bid *models.Bid) error{
  collection := c.client.Database("production").Collection("bids")
  insertResult, err := collection.InsertOne(context.TODO(), bid)
  if err != nil {
      log.Fatal(err)
  }
  fmt.Println("Inserted bid: ", insertResult.InsertedID)
  return err
}

//Get list of an Auction with ID
func (c *DBClient)GetAuction(id string) *models.Auction{
  var auction models.Auction
  collection := c.client.Database("production").Collection("auctions")
  filter := bson.D{{"AuctionID", id}}
  err := collection.FindOne(context.TODO(), filter).Decode(&auction)
  if err!=nil {
    log.Fatal(err)
  }

  return &auction
}

//Get list of all the Auctions
func (c *DBClient)GetAuctions() *models.AuctionList{
  var auctions models.AuctionList
  collection := c.client.Database("production").Collection("auctions")
  cur, err := collection.Find(context.Background(), bson.D{{}})

  if err!=nil {
    log.Fatal(err)
  }
  defer cur.Close(context.Background())
  for cur.Next(context.Background()) {
    var elem models.Auction
    err := cur.Decode(&elem)
    auctions = append(auctions, elem)
    if err != nil {
      log.Fatal(err)
    }
  }
  if err := cur.Err(); err != nil {
    log.Fatal(err)
  }

  return &auctions
}

//get the list of all the bids
func (c *DBClient)GetBids() *[]models.Bid{
  var bids []models.Bid
  collection := c.client.Database("production").Collection("bid")
  cur, err := collection.Find(context.Background(), bson.D{{}})

  if err!=nil {
    log.Fatal(err)
  }
  defer cur.Close(context.Background())
  for cur.Next(context.Background()) {
    var elem models.Bid
    err := cur.Decode(&elem)
    bids = append(bids, elem)
    if err != nil {
      log.Fatal(err)
    }
  }
  if err := cur.Err(); err != nil {
    log.Fatal(err)
  }

  return &bids
}
