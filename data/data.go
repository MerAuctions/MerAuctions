package data

// package main
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/MerAuctions/MerAuctions/db"
	"github.com/MerAuctions/MerAuctions/models"
)

var DBclient *db.DBClient

func GetAllAuctions() *models.AuctionList {
	return DBclient.GetAuctions()
}

func GetAuctionsByAuctionTag(tag string) *[]models.Auction {
	auc, err := DBclient.GetAuctionsByTag(tag)
	if err != nil {
		log.Println("err:", err)
		return nil
	}
	return auc
}

//TODO return error
func GetAuctionById(id string) *models.Auction {
	auc, err := DBclient.GetAuction(id)
	if err != nil {
		log.Println("err:", err)
		return nil
	}
	return auc
}

func GetAuctionByAuctionID(id primitive.ObjectID) *models.Auction {
	auc, err := DBclient.GetAuctionByID(id)
	if err != nil {
		log.Println("err:", err)
		return nil
	}
	return auc
}

func GetAuctionByUserId(id string) *models.AuctionList {
	bids, err := DBclient.GetBidsbyUser(id)
	if err != nil {
		log.Println("err:", err)
		return nil
	}
	log.Println(bids)
	var aucs models.AuctionList
	for _, bid := range *bids {
		auc, _ := DBclient.GetAuction(bid.AuctionID)
		aucs = append(aucs, *auc)
	}
	log.Println(aucs)
	return &aucs
}

func GetTopFiveBids(auctionID string) *[]models.Bid {
	tmp_bids, err := DBclient.GetBids(auctionID)
	log.Println("tmp_bids:", tmp_bids)
	if err != nil {
		log.Println("err:", err)
		return nil //TODO: also give error
	}
	bids := *tmp_bids
	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Time > bids[j].Time
	})
	if len(bids) < 5 {
		return &bids
	}
	result := bids[:5]
	log.Println("result:", result)
	return &result
}

// This function returns all the bids for an auction by price sorted
func GetAllSortedBidsForAuction(auctionID string) []models.Bid {
	tmp_bids, err := DBclient.GetBids(auctionID)
	log.Println(tmp_bids)
	if err != nil {
		log.Println("Error in finding bids ", err)
		return nil //TODO: also give error
	}
	bids := *tmp_bids
	sort.SliceStable(bids, func(i, j int) bool {
		return bids[i].Time > bids[j].Time
	})
	// log.Println("Bids: ", bids)
	return bids
}

// This function returns all user who bid for an auction
func GetAllUsersForAuction(auctionID string) []models.User {
	users, err := DBclient.GetUsers(auctionID)
	if err != nil {
		log.Println("No bidders for auction ", auctionID)
	}

	return *users
}

func AddNewUser(usr *models.User) (*models.User, int) {
	var user *models.User
	var err error

	if usr.UserID == "" {
		log.Println("UserID is empty")
		return user, 2
	} else if usr.Password == "" {
		log.Println("Password is empty")
		return user, 4
	}

	user, err = DBclient.Getuser(string(usr.UserID))

	if err == nil {
		//user already existst, so exists
		log.Println("User already exists")
		return user, 1
	}

	//User doesn't exit and needed to be inserted in the db
	err = DBclient.InsertUser(usr)
	if err != nil {
		//unable to insert user
		log.Println("Error in creating new user (AddNewUser) ", err)
		return user, 5 //TODO: discuss which status code to give
	}

	user, err = DBclient.Getuser(string(usr.UserID))
	if err != nil {
		log.Println("Error in creating new user (Getuser): ", err)
		return user, 5
	}

	log.Println("User signup successful")
	return user, 0
}

//This function returns User by UserID
func GetUserByID(userID string) models.User {
	temp_user, err := DBclient.Getuser(userID)
	if err != nil {
		log.Println(userID, " User not Found!")
	}
	user := *temp_user
	return user
}

//This function updates an User details by ID
func UpdateUser(userID string, points int) error {
	return DBclient.UpdateUser(userID, points)
}

func AddNewAuction(auction *models.Auction) (models.Auction, int) {
	if auction.Title == "" {
		log.Println("Invalid Auction Title")
		return *auction, 2
	} else if len(auction.Image) == 0 {
		log.Println("Please upload auction image")
		return *auction, 3
	}

	id, err := DBclient.InsertAuction(auction)
	if err != nil {
		log.Println("Error in creating new auction")
		return *auction, 1
	}

	auction, err = DBclient.GetAuctionByID(id)
	log.Println("Auction created successfully")
	return *auction, 0
}

func AddNewBid(bid *models.Bid) int {
	//check if the given user and the given auction is present in db
	_, err := DBclient.Getuser(string(bid.UserID))
	if err != nil {
		//user doesn't exit
		return 1
	}

	auc, err := DBclient.GetAuction(string(bid.AuctionID))
	if err != nil {
		//auction doesn't exit
		return 1
	}

	currentTime := int64(time.Now().Unix())
	log.Printf("Current unix time: %v    Time at which auction ends: %v", currentTime, auc.EndTime)
	if currentTime > auc.EndTime {
		return 2 //TODO: discuss which status code to give
	}
	bid.Time = currentTime
	err = DBclient.InsertBid(bid)
	if err != nil {
		//unable to insert bid
		return 2 //TODO: discuss which status code to give
	}

	return 0
}

func GetResult(auctionID string) *models.Result {
	auc, err := DBclient.GetAuction(auctionID)
	if err != nil {
		//auction doesn't exit
		return nil
	}

	currentTime := int64(time.Now().Unix())
	// fmt.Printf("Current unix time: %v    Time at which auction ends: %v", currentTime, auc.EndTime)
	if currentTime < auc.EndTime {
		return nil
	}

	tmp_bids, err := DBclient.GetBids(auctionID)
	if err != nil {
		//cannot obtain all the bids
		return nil
	}

	bids := *tmp_bids
	if len(bids) == 0 {
		// No bid made for this auction
		return nil
	}

	sort.SliceStable(bids, func(i, j int) bool {
		if bids[i].Price == bids[j].Price {
			return bids[i].Time < bids[j].Time
		}
		return bids[i].Price > bids[j].Price
	})
	winningBid := bids[0]
	result := models.Result{
		AuctionID: winningBid.AuctionID,
		WinnerID:  winningBid.UserID,
		Price:     winningBid.Price,
	}

	return &result
}

func GetUserById(id string) (*models.User, error) {
	usr, err := DBclient.Getuser(id)
	if err != nil {
		log.Println("Error fetching user details: ", err)
		return nil, err
	}
	return usr, nil
}

//this will populate the db
func PopulateDB() bool {
	var auc models.AuctionList
	file, err := ioutil.ReadFile("./server/seed-data/auctions.json")
	if err != nil {
		log.Println("Error reading auctions.json : ", err.Error())
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &auc)

	// deleting all the data it exists before
	err = DBclient.DeleteAllCollections()
	if err != nil {
		log.Println("Error in deleting pre-existing data : ", err.Error())
	}

	//setting the time for different aucitons
	auc[0].EndTime = int64(time.Now().Add(time.Hour * 2).Unix())
	auc[1].EndTime = int64(time.Now().Add(time.Hour * 2).Unix())
	auc[2].EndTime = int64(time.Now().Add(time.Minute * 2).Unix())
	err = DBclient.InsertAuctions(&auc)
	if err != nil {
		log.Println("Error populating auctions.json : ", err.Error())
	}
	return true

}
