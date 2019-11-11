package server

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/gin-gonic/gin"
)

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

//get all auctions
func getAllAuctions(c *gin.Context) {
	var allAuctions models.AuctionList
	data.GetAllAuctionsFromDB(&allAuctions)
	c.JSON(200, allAuctions)
}

//get auction by id
func getAuctionsById(c *gin.Context) {
	id := c.Param("auction_id")
	var auc models.Auction
	data.GetAuctionByIdFromDB(&auc, id)
	c.JSON(200, auc)
}

//gets all bids from a auction
func getBidsAuctionsById(c *gin.Context) {
	id := c.Param("auction_id")
	var top5bids [5]models.Bid
	data.GetTopFiveBidsFromDB(&top5bids, id)
	c.JSON(200, top5bids)
}

// register new user
func addNewUser(c *gin.Context) {
	var newuser models.User
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &newuser)

	//status:0-->success, status:1-->user exists
	//TODO: status:2-->userid not according to standard
	status := data.AddNewUserToDB(&newuser)
	if status == 0 {
		c.JSON(200, fmt.Sprintf("User Successfully added"))
	} else {
		c.JSON(400, fmt.Sprintf("User Alredy exists"))
	}

}

//add bid by a registered user
func addBidAuctionIdByUserId(c *gin.Context) {
	var newbid models.Bid
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &newbid)
	auc_id := c.Param("auction_id")
	usr_id := c.Param("user_id")

	newbid.AuctionID = models.ID(auc_id)
	newbid.UserID = models.ID(usr_id)

	//TODO: check for price limits
	status := data.AddNewBid(&newbid)
	if status == 0 {
		c.JSON(200, fmt.Sprintf("User Successfully added"))
	} else {
		c.JSON(400, fmt.Sprintf("User Alredy exists"))
	}
}

func declareResult(auctionID string) models.Bid {
	var auc models.Auction
	data.GetAuctionByIdFromDB(&auc, auctionID)
	var bidlist models.BidList
	data.GetAllBidsFromDB(&bidlist, auctionID)
	sort.Slice(bidlist, func(i, j int) bool {
		return bidlist[i].Price > bidlist[j].Price
	})
	return bidlist[0]
}

//get results of an auction
func getResultByAuctionId(c *gin.Context) {
	id := c.Param("auction_id")
	var aucres models.Result
	var auc models.Auction
	data.GetAuctionByIdFromDB(&auc, id)

	if int64(auc.EndTime) <= time.Now().Unix() {
		//auction completed
		data.GetResult(&aucres, id)
		c.JSON(200, aucres)
	} else {
		c.String(400, fmt.Sprint("Auction Not completed yet"))
	}

}
