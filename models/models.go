package models

import (
	//"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Price of the item or the bid made on the item
type Price int64

//User represents the users making the bid
type User struct {
	UserID   string   `json:"user_id, string"`
	UserName string   `json:"user_name, string"`
	Password string   `json:"pwd, string"`
	Points   int      `json:"points, int"`
	Interest []string `json:"interest"`
}

//Bid represents a single bid in an Auction
type Bid struct {
	AuctionID string `json:"auction_id, string"`
	UserID    string `json:"user_id, string"`
	Price     Price  `json:"price"`
	Count     int    `json:"count, int"`
	Time      int64  `json:"time"`
}

//Auction represents a single auction
type Auction struct {
	AuctionID primitive.ObjectID `bson:"_id, omitempty"`
	// AuctionID   primitive.ObjectID
	Title       string   `json:"title"`
	Image       []string `json:"image"` // image encode in base64
	Tag         []string `json:"tag"`
	Description string   `json:"description"`
	BasePrice   Price    `json:"price"`
	EndTime     int64    `json:"time"`
	CreatedBy   string   `json:"created_by, string"`
}

//AuctionList is a list of auctions
type AuctionList []Auction

// Result to store the result of an auction
type Result struct {
	AuctionID string
	WinnerID  string
	Price     Price
}

// ResponseCreateAuction stores response of /auction/create endpoint
type ResponseCreateAuction struct {
	Message string
	Auction Auction
}

// ResponseSignup stores response of /user/signup endpoint
type ResponseSignup struct {
	Message string
	User    User
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
