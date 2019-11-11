package models

import (
  //"time"
)

//Price of the item or the bid made on the item
type Price float32

//ID represents the id of an auciton, user or a bid
type ID string

//User represents the users making the bid
type User struct{
  UserID ID
  UserName string
}

//Bid represents a single bid in an Auction
type Bid struct{
  AuctionID ID
  UserID ID
  Price Price
  Time int64
}

//Auction represents a single auction
type Auction struct{
  Image string      // image encode in base64
  Description string
  EndTime int64
}

//AuctionList is a list of auctions
type AuctionList []Auction

type Result struct{
  AuctionID ID
  WinnerID ID
  Price Price
}

//List of most recent bids
type LatestBids []Bid
