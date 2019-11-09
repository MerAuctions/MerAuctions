package data

import (
	"github.com/MerAuctions/MerAuctions/models"
)

func GetAllAuctionsFromDB(*models.AuctionList) {

}

func GetAuctionByIdFromDB(auc *models.Auction, id string) {

}

func GetTopFiveBidsFromDB(top5bids *[5]models.Bid, id string) {

}

func AddNewUserToDB(usr *models.User) int {

	return 0
}

func AddNewBid(bid *models.Bid) int {

	return 0
}

func GetResult(res *models.Result, id string) {

}
