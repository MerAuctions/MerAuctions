package models

import (
	//"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Price of the item or the bid made on the item
type Price int64

//User represents the users making the bid
type User struct {
	UserID   string `json:"user_id, string"`
	UserName string
	Password string `json:"pwd, string"`
	Points   int    `json:"Points int"`
	Interest []string
}

//Bid represents a single bid in an Auction
type Bid struct {
	AuctionID string
	UserID    string
	Price     Price
	Count     int
	Time      int64
}

//Auction represents a single auction
type Auction struct {
	AuctionID primitive.ObjectID `bson:"_id, omitempty"`
	// AuctionID   primitive.ObjectID
	Title       string
	Image       []string // image encode in base64
	Tag         []string
	Description string
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

// Response stores response of /auction/create endpoint
type Response struct {
	Message string
	Auction Auction
}

//BidList to store List of bids
type BidList []Bid

//Tag to store image tag
type Tag struct {
	Name       string
	Confidence float64
}

//TagList to store List of image tags
type TagList []Tag

//Description to store image description
type Description string
