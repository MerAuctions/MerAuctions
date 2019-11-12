package db

// package main

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
	URL    string
	DBname string
	client *mongo.Client
}

//Connect to mongoDB of given URL
func ConnectDB(url string, db string) *DBClient {
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
		URL:    url,
		DBname: db,
		client: client,
	}
}

//Insert a user in db
func (c *DBClient) InsertUser(usr *models.User) error {
	collection := c.client.Database(c.DBname).Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), usr)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	fmt.Println("Inserted user: ", insertResult.InsertedID)
	return err
}

//Insert a bid in db
func (c *DBClient) InsertBid(bid *models.Bid) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	insertResult, err := collection.InsertOne(context.TODO(), bid)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	fmt.Println("Inserted bid: ", insertResult.InsertedID)
	return err
}

func (c *DBClient) InsertBids(bids *[]models.Bid) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	for _, bid := range *bids {
		insertResult, err := collection.InsertOne(context.TODO(), bid)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		fmt.Println("Inserted bid: ", insertResult.InsertedID)
	}
	return nil
}

func (c *DBClient) DeleteBid(bid *models.Bid) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	_, err := collection.DeleteOne(context.TODO(), bid)
	if err != nil {
		//log.Fatal(err)
		return err
	}
	fmt.Println("Deleted bid")
	return err
}

func (c *DBClient) DeleteBids(bids *[]models.Bid) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	for _, bid := range *bids {
		_, err := collection.DeleteOne(context.TODO(), bid)
		if err != nil {
			//log.Fatal(err)
			return err
		}
	}
	fmt.Println("Deleted all bids")
	return nil
}

//Insert a auction in db
func (c *DBClient) InsertAuction(auction *models.Auction) error {
	collection := c.client.Database(c.DBname).Collection("auctions")
	insertResult, err := collection.InsertOne(context.TODO(), auction)
	if err != nil {
		// log.Fatal(err)
		return err
	}
	fmt.Println("Inserted auction: ", insertResult.InsertedID)
	return err
}

func (c *DBClient) DeleteAuction(auction *models.Auction) error {
	collection := c.client.Database(c.DBname).Collection("auctions")
	_, err := collection.DeleteOne(context.TODO(), auction)
	if err != nil {
		// log.Fatal(err)
		return err
	}
	fmt.Println("Inserted auction: ", auction)
	return err
}

func (c *DBClient) InsertAuctions(auctions *[]models.Auction) error {
	collection := c.client.Database(c.DBname).Collection("auctions")
	for _, auc := range *auctions {
		insertResult, err := collection.InsertOne(context.TODO(), auc)
		if err != nil {
			// log.Fatal(err)
			return err
		}
		fmt.Println("Inserted auction: ", insertResult.InsertedID)
	}
	return nil
}

func (c *DBClient) DeleteAuctions(auctions *[]models.Auction) error {
	collection := c.client.Database(c.DBname).Collection("auctions")
	for _, auc := range *auctions {
		_, err := collection.DeleteOne(context.TODO(), auc)
		if err != nil {
			// log.Fatal(err)
			return err
		}
	}
	fmt.Println("Deleted all auctions")
	return nil
}

//Get an Auction with ID
func (c *DBClient) GetAuction(id string) (*models.Auction, error) {
	var auction models.Auction
	collection := c.client.Database(c.DBname).Collection("auctions")
	filter := bson.D{{"auctionid", id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&auction)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	return &auction, nil
}

//Get list of all the Auctions
func (c *DBClient) GetAuctions() *models.AuctionList {
	var auctions models.AuctionList
	collection := c.client.Database(c.DBname).Collection("auctions")
	cur, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		//log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var elem models.Auction
		err := cur.Decode(&elem)
		auctions = append(auctions, elem)
		if err != nil {
			//log.Fatal(err)
		}
	}
	if err := cur.Err(); err != nil {
		//log.Fatal(err)
	}

	return &auctions
}

//get the list of all the bids
func (c *DBClient) GetBids(AuctionId string) (*[]models.Bid, error) {
	var bids []models.Bid
	collection := c.client.Database(c.DBname).Collection("bids")
	filter := bson.D{{"auctionid", AuctionId}}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var elem models.Bid
		err := cur.Decode(&elem)
		bids = append(bids, elem)
		if err != nil {
			// log.Fatal(err)
			return nil, err
		}
	}
	if err := cur.Err(); err != nil {
		// log.Fatal(err)
		return nil, err
	}
	fmt.Println("here", bids)
	return &bids, nil
}

//get an user by id
func (c *DBClient) Getuser(id string) (*models.User, error) {
	var user models.User
	collection := c.client.Database(c.DBname).Collection("users")
	filter := bson.D{{"userid", id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("usr", user)
	return &user, nil
}

//Following is for testing the db locally
// func main(){
//   dbclient := ConnectDB("mongodb://localhost:27017","test6")
//   fmt.Println(dbclient.URL)
//   usr := models.User{"1","deepak"}
//   dbclient.InsertUser(&usr)
//   bid1 := models.Bid{"1","2","3",323,9301233}
//   bid2 := models.Bid{"2","2","1",99.823,9238478}
//   bid3 := models.Bid{"3","2","5",101.23,2984792}
//   bid4 := models.Bid{"4","3","6",834.823,398374}
//   bid5 := models.Bid{"5","6","7",8934.823,2349879}
//   dbclient.InsertBid(&bid1)
//   dbclient.InsertBid(&bid2)
//   dbclient.InsertBid(&bid3)
//   dbclient.InsertBid(&bid4)
//   dbclient.InsertBid(&bid5)
//
// 	auc := models.Auction{"1","thisisahashofanimage","thisisadesc.",298347289}
// 	fmt.Println("inserting an auciton")
// 	fmt.Println(dbclient.InsertAuction(&auc))
//   fmt.Println("getting all the bids")
//   fmt.Println(dbclient.GetBids("6"))
//   fmt.Println("Getting all the auctions")
//   fmt.Println(dbclient.GetAuctions())
// 	fmt.Println("Getting the user")
//   fmt.Println(dbclient.Getuser("1"))
//
//
// }

// type DBClient struct {
// 	dbName string
// 	client *mongo.Client
// }
//
// //Connect to mongoDB of given URL
// func ConnectDB(URL string, dbName string) *DBClient {
// 	// Set client options
// 	clientOptions := options.Client().ApplyURI(URL)
//
// 	// Connect to MongoDB
// 	client, err := mongo.Connect(context.TODO(), clientOptions)
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Check the connection
// 	err = client.Ping(context.TODO(), nil)
//
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	fmt.Println("Connected to MongoDB!")
// 	return &DBClient{
// 		client: client,
// 		dbName: dbName,
// 	}
// }
//
// func (c *DBClient) TestAddData() error {
// 	collection := c.client.Database(c.dbName).Collection("numbers")
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
// 	//inserting data
// 	res, err := collection.InsertOne(ctx, models.User{
// 		UserID:   "test1",
// 		UserName: "harsh",
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("Added document with ID: %s", res.InsertedID)
// 	return nil
// }
//
// func (c *DBClient) TestGetData() (*models.User, error) {
// 	collection := c.client.Database(c.dbName).Collection("numbers")
// 	var result models.User
// 	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
// 	//querying data
// 	err := collection.FindOne(ctx, bson.M{"username": "harsh"}).Decode(&result)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &result, nil
// }
