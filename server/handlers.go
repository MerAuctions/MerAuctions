package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	// "io/ioutil"

	"github.com/MerAuctions/MerAuctions/api"
	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var maxBidsToRewards int = 5

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

//getAllAuctions is handler function to return list of all auctions
func getAllAuctions(c *gin.Context) {
	allAuctions := data.GetAllAuctions()
	// log.Println(allAuctions)
	c.HTML(http.StatusOK, "auction_list/index.tmpl", allAuctions)
}

//getAuctionsByID is handler function for getting particular auction page
func getAuctionsByID(c *gin.Context) {
	id := c.Param("auction_id")
	auc := data.GetAuctionById(id)
	if auc == nil {
		c.JSON(404, fmt.Sprintf("Given id: %v not found", id))
		return
	}
	fmt.Println("At the start: auction ", auc.AuctionID, " the end")

	top5bids := data.GetTopFiveBids(id)
	fmt.Println("top 5 bids ", top5bids)

	if top5bids == nil {
		c.JSON(404, fmt.Sprintf("Given id: %v not found", id))
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
		"auctionID":      auc.AuctionID.Hex(),
	})
}

//getAuctionsByID is handler function for getting particular auction page
func getCreateAuction(c *gin.Context) {
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
	c.HTML(http.StatusOK, "create_auction/index.tmpl", gin.H{
		"isUserSignedIn": isUserSignedIn,
	})
}

//getBidsAuctionsById is handler function to get all bids from a auction
func getBidsAuctionsById(c *gin.Context) {
	id := c.Param("auction_id")
	top5bids := data.GetTopFiveBids(id)
	if top5bids == nil {
		c.JSON(404, fmt.Sprintf("Given id: %v not found", id))
		return
	}
	c.JSON(200, top5bids)
}

//addNewUser registers a new user
func addNewUser(c *gin.Context) {
	var newuser models.User
	var responseSignup models.ResponseSignup
	var responseCode int

	c.ShouldBindJSON(&newuser)

	newUser, status := data.AddNewUser(&newuser)
	log.Println("User object: ", newUser)
	log.Println("Status: ", status)

	if newUser != nil {
		responseSignup.User = *newUser
	}

	if status == 0 {
		usr := models.User{
			UserID: newuser.UserID,
		}
		token, _, _ := authMiddleware.TokenGenerator(&usr)
		log.Println("cookie token ", token)
		c.SetCookie("token", token, 60*60, "/", "", false, false)
		responseCode = 200
		responseSignup.Message = "User signup successful"
	} else if status == 1 {
		responseCode = 500
		responseSignup.Message = "User already exists"
	} else if status == 2 {
		responseCode = 500
		responseSignup.Message = "UserID is empty"
	} else if status == 4 {
		responseCode = 500
		responseSignup.Message = "Password is empty"
	} else if status == 5 {
		responseCode = 500
		responseSignup.Message = "Error in creating new user"
	}

	c.JSON(responseCode, responseSignup)
}

// createAuction create a new auction for the user
func createAuction(c *gin.Context) {
	var newAuction models.Auction
	var response models.ResponseCreateAuction
	var responseCode int

	c.ShouldBindJSON(&newAuction)

	log.Println(newAuction)

	newAuction, status := data.AddNewAuction(&newAuction)
	response.Auction = newAuction

	if status == 0 {
		response.Message = "Auction Successfully created."
		responseCode = 200
	} else if status == 2 {
		response.Message = "Invalid Auction Title"
		responseCode = 500
	} else if status == 3 {
		response.Message = "Please upload auction image"
		responseCode = 500
	} else {
		response.Message = "Error in creating new auction"
		responseCode = 500
	}

	c.JSON(responseCode, response)
}

//addBidAuctionIdByUserId is handler function to add bid by a registered user
func addBidAuctionIdByUserId(c *gin.Context) {

	isUserSignedIn := false
	usr_id := ""
	if jwtToken, err := authMiddleware.ParseToken(c); err == nil {
		if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
			if userID, ok := claims[jwtIdentityKey].(string); ok == true {
				if user, _ := data.GetUserById(userID); user != nil {
					isUserSignedIn = true
					usr_id = userID
				}
			} else {
				log.Printf("Could not convert claim to string")
			}
		} else {
			log.Printf("Could not extract claims into JWT")
		}
	}

	if isUserSignedIn == false {
		log.Printf("User not logged in because of not signed in.")
		c.JSON(400, fmt.Sprintf("User is not logged in."))
		return
	}

	var newbid models.Bid
	price_map := gin.H{"price": ""}
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &price_map)

	auc_id := c.Param("auction_id")

	str_price, ok := price_map["price"].(string)
	if ok == false {
		log.Println("Invalid bid: error my converting interface to string")
		c.JSON(400, fmt.Sprintf("Invalid bid"))
		return
	}
	tmp_price, err := strconv.ParseFloat(str_price, 32)
	if err != nil {
		log.Println(err)
		c.JSON(400, fmt.Sprintf("Invalid bid"))
		return
	}

	newbid.AuctionID = auc_id
	newbid.UserID = usr_id
	newbid.Price = models.Price(tmp_price)

	//TODO: check for price limits
	status := data.AddNewBid(&newbid)
	if status == 0 {
		log.Printf("User's Bid Successfully added.")
		c.JSON(200, fmt.Sprintf("User's Bid Successfully added"))
	} else {
		log.Printf("User's Bid could not be added with status %d.", status)
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
		c.String(400, fmt.Sprint("Auction Not completed yet"))
	}
}

//this will populate the db
func addDataDB(c *gin.Context) {
	ok := data.PopulateDB()
	if ok == false {
		c.String(400, "Can't populate DB")
	} else {
		c.String(200, "DB is populated Successfully")
	}

}

//addRewardsToUsers is handler function to offer rewards when the auction ends
func addRewardsToUsers(c *gin.Context) {
	rewardPercentage := 0.005
	auctionID := c.Param("auction_id")
	auc := data.GetAuctionById(auctionID)

	bids := data.GetAllSortedBidsForAuction(auctionID)
	userBidFreq := make(map[string]int)

	for _, bid := range bids {
		freq, ok := userBidFreq[bid.UserID]
		if ok == false {
			userBidFreq[bid.UserID] = 1
		} else {
			userBidFreq[bid.UserID] = freq + 1
		}

		if freq <= maxBidsToRewards+1 {
			pointsForBidPrice := (rewardPercentage * float64(bid.Price))
			pointsForHighBid := float64(bid.Price-2*auc.BasePrice) / float64(2*auc.BasePrice)

			//TODO after auction creation done
			//pointsFromTime := float64(duration*60/(auc.EndTime - bid.Time))

			points := int(pointsForHighBid * pointsForBidPrice)
			if points <= 0 {
				continue
			} else {
				err := data.UpdateUser(bid.UserID, points)
				if err != nil {
					c.JSON(404, fmt.Sprint("User Not Found!"))
				}
			}
		}
	}

	c.JSON(200, fmt.Sprintf("Bidders are rewarded!"))

}

func addRewardsToUser(c *gin.Context) {
	rewardPercentage := 0.005
	auctionID := c.Param("auction_id")
	userID := c.Param("user_id")
	user, err := data.GetUserById(userID)

	if err != nil {
		log.Printf("Error in getting user details for '%s' from DB: %s", userID, err.Error())
	}
	auc := data.GetAuctionById(auctionID)
	bids := data.GetAllSortedBidsForAuction(auctionID)
	userBidFreq := 0

	for _, bid := range bids {
		if bid.UserID == userID {
			if userBidFreq < maxBidsToRewards {
				pointsForBidPrice := (rewardPercentage * float64(bid.Price))
				fmt.Println("auc.BasePrice:", auc.BasePrice)
				pointsForHighBid := float64(bid.Price-2*auc.BasePrice) / float64(2*auc.BasePrice)

				//TODO after auction creation done
				//pointsFromTime := float64(duration*60/(auc.EndTime - bid.Time))

				points := int(pointsForHighBid * pointsForBidPrice)
				fmt.Println("pointsForBidPrice:", pointsForBidPrice, "  pointsForHighBid:", pointsForHighBid, "  points:", points)
				if points <= 0 {
					continue
				} else {
					err := data.UpdateUser(bid.UserID, points)
					if err != nil {
						c.JSON(404, fmt.Sprint("User Not Found!"))
					}
				}
			}
			userBidFreq++
		}
	}
	user, err = data.GetUserById(userID)
	if err != nil {
		log.Printf("Error in getting user details for '%s' from DB: %s", userID, err.Error())
	}
	c.JSON(200, fmt.Sprintf("Congrats!, you have been awarded %v points.", user.Points))

}

// get picture user uploaded and save to /media/images
func uploadPicture(c *gin.Context) {
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	filename := filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, "media/images/"+filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
}

func getTagsfromImage(c *gin.Context) {
	imageName := c.Request.URL.Query().Get("imageName")
	tags := api.GetTagsForImage(imageName)
	c.JSON(http.StatusOK, tags)
}
