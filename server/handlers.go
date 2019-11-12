package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/gin-gonic/gin"

	"github.com/dgrijalva/jwt-go"
)

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

//getAllAuctions is handler function to return list of all auctions
func getAllAuctions(c *gin.Context) {
	allAuctions := data.GetAllAuctions()
	c.HTML(http.StatusOK, "auction_list/index.tmpl", allAuctions)
}

//getAuctionsByID is handler function for getting particular auction page
func getAuctionsByID(c *gin.Context) {
	id := c.Param("auction_id")
	auc := data.GetAuctionById(id)
	if auc == nil {
		c.JSON(200, fmt.Sprintf("Given id: %v not found", id))
		return
	}
	fmt.Println("At the start: auction ",auc.AuctionID," the end")

	top5bids := data.GetTopFiveBids(id)
	fmt.Println("top 5 bids ",top5bids)

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

//getBidsAuctionsById is handler function to get all bids from a auction
func getBidsAuctionsById(c *gin.Context) {
	id := c.Param("auction_id")
	top5bids := data.GetTopFiveBids(id)
	if top5bids == nil {
		c.JSON(200, fmt.Sprintf("Given id: %v not found", id))
		return
	}
	c.JSON(200, top5bids)
}

//addNewUser registers a new user
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

//addBidAuctionIdByUserId is handler function to add bid by a registered user
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

//getResultByAuctionId is handler function to check result by ID
func getResultByAuctionId(c *gin.Context) {
	id := c.Param("auction_id")
	auc := data.GetAuctionById(id)

	if int64(auc.EndTime) <= time.Now().Unix() {
		//auction completed
		aucres := data.GetResult(id)
		c.JSON(200, aucres)
	} else {
		c.String(200, fmt.Sprint("Auction Not completed yet"))
	}
}
