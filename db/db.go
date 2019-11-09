package db

import (
	"context"
	"fmt"
	"log"

	//"github.com/MerAuctions/MerAuctions/models"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Connect to mongoDB of given URL
func ConnectDB(URL string) *mongo.Client {
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
