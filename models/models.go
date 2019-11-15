package models

import (
	//"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Price of the item or the bid made on the item
type Price int

//User represents the users making the bid
type User struct {
	UserID   string `json:"user_id, string"`
	UserName string
	Password string `json:"pwd, string"`
	Interest []string
	Points   int
}

//Bid represents a single bid in an Auction
type Bid struct {
	AuctionID string
	UserID    string
	Price     Price
	Time      int64
	Count     int
}

//Auction represents a single auction
type Auction struct {
	AuctionID primitive.ObjectID `bson:"_id, omitempty"`
	// AuctionID   primitive.ObjectID
	Title       string
	Image       []string // image encode in base64
	Description string
	Tag         []string
	BasePrice   Price
	EndTime     int64
}

//AuctionList is a list of auctions
type AuctionList []Auction

// Result to store the result of an auction
type Result struct {
	AuctionID string
	WinnerID  string
	Price     Price
}

//BidList to store List of bids
type BidList []Bid
