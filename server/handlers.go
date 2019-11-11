package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"
)

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

//get all auctions
func getAllAuctions(c *gin.Context) {
	allAuctions := data.GetAllAuctions()
	c.HTML(http.StatusOK, "auction_list/index.tmpl", allAuctions)
}

//getAuctionsById is handler function for getting particular auction page
func getAuctionsById(c *gin.Context) {
	id := c.Param("auction_id")
	auc := data.GetAuctionById(id)
	if auc == nil {
		c.JSON(200, fmt.Sprintf("Given id: %v not found", id))
		return
	}

	top5bids := data.GetTopFiveBids(id)
	if top5bids == nil {
		c.JSON(200, fmt.Sprintf("Given id: %v not found", id))
		return
	}
	isUserSignedIn := false
	if jwtToken, err := authMiddleware.ParseToken(c); err == nil {
		if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
			if userID, ok := claims[jwtIdentityKey].(string); ok == true {
				if user, _ := data.GetUserById(userID); user != nil {
					isUserSignedIn = true
				}
			} else {
				log.Printf("Could not convert claim to string")
			}
		} else {
			log.Printf("Could not extract claims into JWT")
		}
	}
	c.HTML(http.StatusOK, "auction/index.tmpl", gin.H{
		"auction":        auc,
		"bids":           top5bids,
		"isUserSignedIn": isUserSignedIn,
	})
}

//gets all bids from a auction
func getBidsAuctionsById(c *gin.Context) {
	id := c.Param("auction_id")
	top5bids := data.GetTopFiveBids(id)
	if top5bids == nil {
		c.JSON(200, fmt.Sprintf("Given id: %v not found", id))
		return
	}
	c.JSON(200, top5bids)
}

// register new user
func addNewUser(c *gin.Context) {
	var newuser models.User
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &newuser)

	//status:0-->success, status:1-->user exists
	//TODO: status:2-->userid not according to standard
	status := data.AddNewUser(&newuser)
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
		c.JSON(200, fmt.Sprintf("User's Bid Successfully added"))
	} else {
		c.JSON(400, fmt.Sprintf("User's Bid could not be added"))
	}
}

func declareResult(c *gin.Context, auctionID string) {
	var auc models.Auction
	data.GetAuctionByIdFromDB(&auc, auctionID)
	var bidlist models.BidList
	data.GetAllBidsFromDB(&bidlist, auctionID)
	sort.Slice(bidlist, func(i, j int) bool {
		return bidlist[i].Price > bidlist[j].Price
	})
	var winnerBid models.Result
	winnerBid.AuctionID = bidlist[0].AuctionID
	winnerBid.Price = bidlist[0].Price
	winnerBid.WinnerID = bidlist[0].UserID
	c.JSON(200, winnerBid)
	data.PushResultDB(&winnerBid)
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
