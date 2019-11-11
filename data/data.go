package data

import (
	// "context"
	"fmt"
	// "log"
	"sort"

	"github.com/MerAuctions/MerAuctions/db"
	"github.com/MerAuctions/MerAuctions/models"
)

var DBclient *db.DBClient

func GetAllAuctions() *models.AuctionList {
	return DBclient.GetAuctions()
}

//TODO return error
func GetAuctionById(id string) *models.Auction {
	auc, err := DBclient.GetAuction(id)
	if err != nil {
		return nil
	}
	return auc
}

func GetTopFiveBids(auctionID string) *[]models.Bid {
	temp_bids, err := DBclient.GetBids(auctionID)
	if err != nil {
		return nil
	}
	bids := *temp_bids
	fmt.Println(bids)
	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Time > bids[j].Time
	})
	if len(bids) < 5 {
		return &bids
	}
	result := bids[:5]
	return &result
}

func AddNewUser(usr *models.User) int {

	return 0
}

func AddNewBid(bid *models.Bid) int {

	return 0
}

func GetResult(auctionID string) *models.Result {
	return nil
}

func GetUserById(id string) (*models.User, error) {
	return &models.User{
		UserID:   "harsh",
		Password: "temp",
	}, nil
}
