package db

// package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MerAuctions/MerAuctions/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		log.Println(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Println(err)
	}

	fmt.Println("Connected to MongoDB!")
	return &DBClient{
		URL:    url,
		DBname: db,
		client: client,
	}
}

//Insert a user in db
func (c *DBClient) DeleteAllCollections() error {
	if err := c.client.Database(c.DBname).Collection("users").Drop(context.TODO()); err != nil {
		if err != nil {
			log.Println("Could not drop users collection in database")
			return err
		}
	}
	log.Println("Dropped users collection from database")

	if err := c.client.Database(c.DBname).Collection("bids").Drop(context.TODO()); err != nil {
		if err != nil {
			log.Println("Could not drop bids collection in database")
			return err
		}
	}
	log.Println("Dropped bids collection from database")

	if err := c.client.Database(c.DBname).Collection("auctions").Drop(context.TODO()); err != nil {
		if err != nil {
			log.Println("Could not drop auctions collection in database")
			return err
		}
	}
	log.Println("Dropped auctions collection from database")
	return nil
}

//Insert a user in db
func (c *DBClient) InsertUser(usr *models.User) error {
	collection := c.client.Database(c.DBname).Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), usr)
	if err != nil {
		log.Println("Error in creating new user(Insert User) : ", err)
		return err
	}
	log.Println("Inserted user: ", insertResult.InsertedID)
	return nil
}

func (c *DBClient) InsertUsers(users *[]models.User) error {
	collection := c.client.Database(c.DBname).Collection("users")
	for _, usr := range *users {
		insertResult, err := collection.InsertOne(context.TODO(), usr)
		if err != nil {
			log.Println("Error in inserting users: ", err)
			return err
		}
		log.Println("Inserted user: ", insertResult.InsertedID)
	}
	return nil
}

func (c *DBClient) DeleteUsers(users *[]models.User) error {
	collection := c.client.Database(c.DBname).Collection("users")
	for _, user := range *users {
		_, err := collection.DeleteOne(context.TODO(), user)
		if err != nil {
			return err
		}
	}
	log.Println("Deleted all users")
	return nil
}

//Insert a bid in db
func (c *DBClient) InsertBid(bid *models.Bid) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	insertResult, err := collection.InsertOne(context.TODO(), bid)
	if err != nil {
		log.Println("err:", err)
		return err
	}
	log.Println("Inserted bid: ", insertResult.InsertedID)
	return err
}

func (c *DBClient) InsertBids(bids *models.BidList) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	for _, bid := range *bids {
		insertResult, err := collection.InsertOne(context.TODO(), bid)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Inserted bid: ", insertResult.InsertedID)
	}
	return nil
}

func (c *DBClient) DeleteBid(bid *models.Bid) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	_, err := collection.DeleteOne(context.TODO(), bid)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Deleted bid")
	return nil
}

func (c *DBClient) DeleteBids(bids *models.BidList) error {
	collection := c.client.Database(c.DBname).Collection("bids")
	for _, bid := range *bids {
		_, err := collection.DeleteOne(context.TODO(), bid)
		if err != nil {
			return err
		}
	}
	log.Println("Deleted all bids")
	return nil
}

//Insert a auction in db
func (c *DBClient) InsertAuction(auction *models.Auction) (primitive.ObjectID, error) {
	collection := c.client.Database(c.DBname).Collection("auctions")
	auction.AuctionID = primitive.NewObjectIDFromTimestamp(time.Now())
	insertResult, err := collection.InsertOne(context.TODO(), auction)
	id := insertResult.InsertedID.(primitive.ObjectID)

	if err != nil {
		log.Println(err)
		return id, err
	}
	log.Println("Inserted auction: ", insertResult.InsertedID)
	return id, err
}

func (c *DBClient) DeleteAuction(auction *models.Auction) error {
	collection := c.client.Database(c.DBname).Collection("auctions")
	_, err := collection.DeleteOne(context.TODO(), auction)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Inserted auction: ", auction)
	return nil
}

func (c *DBClient) InsertAuctions(auctions *models.AuctionList) error {
	collection := c.client.Database(c.DBname).Collection("auctions")
	for _, auc := range *auctions {
		insertResult, err := collection.InsertOne(context.TODO(), auc)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Inserted auction: ", insertResult.InsertedID)
	}
	return nil
}

func (c *DBClient) DeleteAuctions(auctions *models.AuctionList) error {
	ctx, _ := context.WithTimeout(context.Background(), 1000*time.Second)
	collection := c.client.Database(c.DBname).Collection("auctions")
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return err
		// for _, auc := range *auctions {
		// 	_, err := collection.DeleteOne(context.TODO(), auc)
		// 	if err != nil {
		// 		// log.Fatal(err)
		// 		return err
		// 	}
	}
	log.Println("Deleted all auctions")
	return nil
}

//Get an Auction by Tag
func (c *DBClient) GetAuctionsByTag(tag string) (*[]models.Auction, error) {
	var auctions []models.Auction
	collection := c.client.Database(c.DBname).Collection("auctions")
	log.Println(tag)

	filter := bson.D{{"tag", tag}}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Println("Error in fetching auction by Tag ", err)
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var elem models.Auction
		err := cur.Decode(&elem)
		auctions = append(auctions, elem)
		if err != nil {
			// log.Fatal(err)
			return nil, err
		}
	}
	if err := cur.Err(); err != nil {
		// log.Fatal(err)
		return nil, err
	}
	return &auctions, nil

}

//Get an Auction with ID
func (c *DBClient) GetAuction(id string) (*models.Auction, error) {
	var auction models.Auction
	collection := c.client.Database(c.DBname).Collection("auctions")
	log.Println(id)
	docID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		log.Println("Error in fetching auction(GetAuction) ", err)
		return nil, err
	}
	filter := bson.D{{"_id", docID}}
	err = collection.FindOne(context.TODO(), filter).Decode(&auction)
	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	return &auction, nil
}

func (c *DBClient) GetAuctionByID(id primitive.ObjectID) (*models.Auction, error) {
	var auction models.Auction
	collection := c.client.Database(c.DBname).Collection("auctions")
	filter := bson.D{{"_id", id}}
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
		log.Println("Error in fetching all auctions", err.Error())
		return nil
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var elem models.Auction
		err := cur.Decode(&elem)
		auctions = append(auctions, elem)
		if err != nil {
			log.Println("Error in decoding auctions", err.Error())
			return nil
		}
	}
	if err := cur.Err(); err != nil {
		log.Println("Error in reading auctionst list", err.Error())
		return nil
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
	return &bids, nil
}

//get the list of all the bids by user id
func (c *DBClient) GetBidsbyUser(UserId string) (*[]models.Bid, error) {
	var bids []models.Bid
	collection := c.client.Database(c.DBname).Collection("bids")
	filter := bson.D{{"userid", UserId}}
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
	return &bids, nil
}

//get an user by id
func (c *DBClient) Getuser(id string) (*models.User, error) {
	var user models.User
	collection := c.client.Database(c.DBname).Collection("users")
	filter := bson.D{{"userid", id}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//Update User in db by UserID
func (c *DBClient) UpdateUser(userID string, points int) error {
	collection := c.client.Database(c.DBname).Collection("users")

	filter := bson.D{{"userid", userID}}
	update := bson.D{
		{"$inc", bson.D{
			{"points", points},
		}},
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return err
}

func (c *DBClient) DeleteAllUsers() error {
	ctx, _ := context.WithTimeout(context.Background(), 1000*time.Second)
	collection := c.client.Database(c.DBname).Collection("users")
	_, err := collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Deleted all users")
	return nil
}

//get the list of all the users
func (c *DBClient) GetUsers(AuctionId string) (*[]models.User, error) {
	var users []models.User
	collection := c.client.Database(c.DBname).Collection("users")
	filter := bson.D{{"auctionid", AuctionId}}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var elem models.User
		err := cur.Decode(&elem)
		users = append(users, elem)
		if err != nil {
			return nil, err
		}
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return &users, nil
}

func ObjectIDToString(id primitive.ObjectID) string {
	return id.Hex()
}

// // Following is for testing the db locally
// func main(){
//   dbclient := ConnectDB("mongodb://localhost:27017","test7")
//   fmt.Println(dbclient.URL)
//   usr := models.User{"1","deepak"}
//   dbclient.InsertUser(&usr)
//   bid1 := models.Bid{"5dc937cc88d9a2eaff817723","3",323,9301233}
//   bid2 := models.Bid{"5dc937cc88d9a2eaff817723","1",99.823,9238478}
//   bid3 := models.Bid{"5dc937cc88d9a2eaff817723","5",101.23,2984792}
//   bid4 := models.Bid{"5dc937cc88d9a2eaff817723","6",834.823,398374}
//   bid5 := models.Bid{"5dc937cc88d9a2eaff817723","7",8934.823,2349879}
//   dbclient.InsertBid(&bid1)
//   dbclient.InsertBid(&bid2)
//   dbclient.InsertBid(&bid3)
//   dbclient.InsertBid(&bid4)
//   dbclient.InsertBid(&bid5)
//
// 	// auc := models.Auction{"thisisahashofanimage","thisisadesc.",298347289}
// 	// fmt.Println("inserting an auciton")
// 	// fmt.Println(dbclient.InsertAuction(&auc))
//   fmt.Println("getting all the bids")
//   fmt.Println(dbclient.GetBids("5dc937cc88d9a2eaff817723"))
//   fmt.Println("Getting all the auctions")
//   fmt.Println(dbclient.GetAuctions())
// 	fmt.Println("Getting the user")
//   fmt.Println(dbclient.Getuser("1"))
// }
