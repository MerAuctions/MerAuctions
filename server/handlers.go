package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	// "io/ioutil"

	"github.com/MerAuctions/MerAuctions/api"
	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var maxBidsToRewards int = 5

func hello(c *gin.Context) {
	c.JSON(200, "hello")
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
	if auc.Tag[len(auc.Tag)-1] == "" {
		auc.Tag = auc.Tag[:len(auc.Tag)-1]
	}
	// fmt.Println("At the start: auction ", auc.AuctionID, " the end")

	top5bids := data.GetTopFiveBids(id)
	// fmt.Println("top 5 bids ", top5bids)

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
	c.HTML(http.StatusOK, "create_auction/index.tmpl", gin.H{
		"isUserSignedIn": isUserSignedIn,
		"user_id":        usr_id,
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
	// log.Println("User object: ", newUser)
	// log.Println("Status: ", status)

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

func getAuctionsByTag(c *gin.Context) {
	tagAuctions := data.GetAuctionsByAuctionTag(c.Param("tag"))

	c.HTML(http.StatusOK, "auction_list/index.tmpl", tagAuctions)
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

	log.Printf("Bid user check success.")

	var newbid models.Bid
	price_map := gin.H{"price": ""}
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &price_map)

	auc_id := c.Param("auction_id")
	auc := data.GetAuctionById(auc_id)

	str_price, ok := price_map["price"].(string)
	if ok == false {
		log.Println("Invalid bid: error my converting interface to string")
		c.JSON(400, fmt.Sprintf("Invalid bid"))
		return
	}
	tmp_price, err := strconv.ParseFloat(str_price, 32)
	if err != nil {
		log.Println("err:", err)
		c.JSON(400, fmt.Sprintf("Invalid bid"))
		return
	}

	currentBid := models.Price(tmp_price)
	bids := data.GetAllSortedBidsForAuction(auc_id)
	var highestBid models.Price = 0
	for _, bid := range bids {
		log.Println(bid)
		if highestBid < bid.Price {
			highestBid = bid.Price
		}
	}

	if len(bids) == 0 && currentBid < auc.BasePrice {
		c.JSON(500, fmt.Sprintf("You can only place bids higher than the base price!"))
	} else {

		if (len(bids) == 0 && currentBid >= auc.BasePrice) || (currentBid > auc.BasePrice && currentBid > highestBid) {
			newbid.AuctionID = auc_id
			newbid.UserID = usr_id
			newbid.Price = models.Price(currentBid)

			log.Println(newbid)

			//TODO: check for price limits
			status := data.AddNewBid(&newbid)
			if status == 0 {
				log.Printf("User's Bid Successfully added.")
				c.JSON(200, fmt.Sprintf("User's Bid Successfully added"))
			} else {
				log.Printf("User's Bid could not be added with status %d.", status)
				c.JSON(400, fmt.Sprintf("User's Bid could not be added"))
			}
			if err != nil {
				c.JSON(404, fmt.Sprint("Auction Not Found!"))
			}
		} else {
			c.JSON(500, fmt.Sprintf("You can only bid above the current highest bid!"))
		}

		log.Println("Done bidding.")
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

func getUserByUserID(c *gin.Context) {
	id := c.Param("user_id")
	user, err := data.GetUserById(id)

	if err != nil {
		c.JSON(500, user)
	}

	c.JSON(200, user)
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
				// fmt.Println("auc.BasePrice:", auc.BasePrice)
				pointsForHighBid := float64(bid.Price-2*auc.BasePrice) / float64(2*auc.BasePrice)

				//TODO after auction creation done
				//pointsFromTime := float64(duration*60/(auc.EndTime - bid.Time))

				points := int(pointsForHighBid * pointsForBidPrice)
				// fmt.Println("pointsForBidPrice:", pointsForBidPrice, "  pointsForHighBid:", pointsForHighBid, "  points:", points)
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

// This function gets personalised Auctions based on User Interests
func getPersonalisedAuctions(c *gin.Context) {
	userID := c.Param("user_id")
	user := data.GetUserByID(userID)
	interests := user.Interest
	auctions := *data.GetAllAuctions()
	similarityMap := make(map[primitive.ObjectID]int)

	for _, auc := range auctions {
		count := 0
		tagMap := make(map[string]int)
		for _, tag := range auc.Tag {
			tagMap[strings.ToLower(tag)] = 0
		}
		for _, interest := range interests {
			_, ok := tagMap[strings.ToLower(interest)]
			if ok == true {
				count++
			}
		}
		similarityMap[auc.AuctionID] = count

	}

	var sortedAuctions entries
	for k, v := range similarityMap {
		sortedAuctions = append(sortedAuctions, entry{val: v, key: k})
	}

	sort.Sort(sort.Reverse(sortedAuctions))

	var personalisedAuctions []primitive.ObjectID

	for _, auc := range sortedAuctions {
		fmt.Println("auction:", auc.key, " similar:", auc.val)
		personalisedAuctions = append(personalisedAuctions, auc.key)
	}

	var aucs []models.Auction

	for _, auc := range personalisedAuctions {
		tmp_auc := data.GetAuctionByAuctionID(auc)
		aucs = append(aucs, *tmp_auc)
	}

	c.JSON(200, aucs)
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
	log.Println("tags : ", tags)
	c.JSON(http.StatusOK, tags)
}

func getDescriptionfromImage(c *gin.Context) {
	imageName := c.Request.URL.Query().Get("imageName")
	description := api.GetDescriptionForImage(imageName)
	log.Println("description : ", description)
	c.JSON(http.StatusOK, description)
}

type entry struct {
	val int
	key primitive.ObjectID
}

type entries []entry

func (s entries) Len() int           { return len(s) }
func (s entries) Less(i, j int) bool { return s[i].val < s[j].val }
func (s entries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func getUserAuctions(c *gin.Context) {
	listAuctions := data.GetAuctionByUserId(c.Param("user_id"))
	log.Println(c.Param("user_id"))
	log.Println(listAuctions)
	c.HTML(http.StatusOK, "auction_list/index.tmpl", listAuctions)
}
