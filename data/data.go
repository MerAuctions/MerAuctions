package data

import (
	"github.com/MerAuctions/MerAuctions/models"
)

func GetAllAuctionsFromDB(*models.AuctionList) {

}

func GetAuctionByIdFromDB(auc *models.Auction, id string) {

}

func GetAllBidsFromDB(allbids *models.BidList, id string) {

}

func GetTopFiveBidsFromDB(top5bids *[5]models.Bid, id string) {

}

func AddNewUserToDB(usr *models.User) int {

	return 0
}

func AddNewBid(bid *models.Bid) int {

	return 0
}

func PushResultDB(*models.Result) bool {

	return true
}

func GetResult(res *models.Result, id string) {

}

func GetUserById(id string) (*models.User, error) {
	return &models.User{
		UserID: "harsh",
		Password: "temp",
	}, nil
}
