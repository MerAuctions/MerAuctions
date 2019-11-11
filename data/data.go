package data

import (
	// "context"
	"fmt"
	// "log"
	"sort"

	"github.com/MerAuctions/MerAuctions/models"
	"github.com/MerAuctions/MerAuctions/db"
)

var DBclient db.DBClient

func GetAllAuctions()*models.AuctionList {
	return DBclient.GetAuctions()
}

func GetAuctionById(id string)*models.Auction {
	return DBclient.GetAuction(id)
}

func GetTopFiveBids(auctionID string)*[]models.Bid {
	bids := *DBclient.GetBids(auctionID)
	fmt.Println(bids)
	sort.SliceStable(bids, func(i, j int) bool{
		return bids[i].Time > bids[j].Time
	})
	if len(bids) < 5{
		return &bids
	}
	result := bids[:5]
	return &result
}

func AddNewUser(usr *models.User) error {

	return nil
}

func AddNewBid(bid *models.Bid) error {

	return nil
}

func GetResult(auctionID string) *models.Result{
	return nil
}
