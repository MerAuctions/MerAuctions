package models

import (
	//"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Price of the item or the bid made on the item
type Price float64

//User represents the users making the bid
type User struct {
	UserID   string `json:"user_id, string"`
	UserName string
	Password string `json:"pwd, string"`
	Interest []string
}

//Bid represents a single bid in an Auction
type Bid struct {
	AuctionID string
	UserID    string
	Price     Price
	Time      int64
}

//Auction represents a single auction
type Auction struct {
	AuctionID   primitive.ObjectID `bson:"_id, omitempty"`
	Title       string
	Image       string // image encode in base64
	Description string
	EndTime     int64
}

//AuctionList is a list of auctions
type AuctionList []Auction

type Result struct {
	AuctionID string
	WinnerID  string
	Price     Price
}

//List of bids
type BidList []Bid
