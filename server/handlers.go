package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/dgrijalva/jwt-go"
)

func hello(c *gin.Context) {
	c.String(200, "Hello World")
}

//gets loginurl for google sign in
func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

//random state for the session
func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

//google login handler
func googleloginHandler(c *gin.Context) {
	state = randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}

//google auth page
func googleauthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := conf.Client(oauth2.NoContext, tok)
	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	fmt.Println(email)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	var userdetails models.GUser
	json.Unmarshal(data, &userdetails)
	log.Println("Email body: ", userdetails.Email)
	c.Status(http.StatusOK)
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
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &newbid)
	auc_id := c.Param("auction_id")

	newbid.AuctionID = auc_id
	newbid.UserID = usr_id

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
