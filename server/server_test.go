package server_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/MerAuctions/MerAuctions/data"
	"github.com/MerAuctions/MerAuctions/models"
	. "github.com/MerAuctions/MerAuctions/server"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func InsertAuctionsToDB() *models.AuctionList {
	var auc models.AuctionList
	file, err := ioutil.ReadFile("./seed-data/auctions.json")
	if err != nil {
		log.Fatal("Error reading auctions.json : ", err.Error())
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.InsertAuctions(&auc)
	if err != nil {
		log.Fatal("Error populating auctions.json : ", err.Error())
	}
	return &auc
}

func InsertUsersToDB() *[]models.User {
	var users []models.User
	file, err := ioutil.ReadFile("./seed-data/users.json")
	if err != nil {
		log.Fatal("Error reading users.json : ", err.Error())
	}
	json.Unmarshal([]byte(file), &users)
	err = data.DBclient.InsertUsers(&users)
	if err != nil {
		log.Fatal("Error populating users.json : ", err.Error())
	}
	// log.Println(users)
	return &users
}

func GetAuctions() *models.AuctionList {
	var auc models.AuctionList
	file, err := ioutil.ReadFile("./seed-data/auctions.json")
	if err != nil {
		log.Fatal("Error reading auctions.json : ", err.Error())
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &auc)

	return &auc
}

func RemoveAuctionsFromDB() {
	var auc models.AuctionList
	file, err := ioutil.ReadFile("./seed-data/auctions.json")
	if err != nil {
		log.Fatal("Error reading auctions.json : ", err.Error())
	}
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.DeleteAuctions(&auc)
	if err != nil {
		log.Fatal("Error deleting auctions.json : ", err.Error())
	}
}

func InsertBidsToDB() *models.BidList {
	var bids models.BidList
	file, err := ioutil.ReadFile("./seed-data/bids.json")
	if err != nil {
		log.Fatal("Error reading bids.json : ", err.Error())
	}
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.InsertBids(&bids)
	if err != nil {
		log.Fatal("Error populating bids.json : ", err.Error())
	}
	return &bids
}

func InsertBidsForAuction(id string) *models.BidList {
	var bids models.BidList
	file, err := ioutil.ReadFile("./seed-data/bids.json")
	if err != nil {
		log.Fatal("Error reading bids.json : ", err.Error())
	}

	json.Unmarshal([]byte(file), &bids)

	// log.Println(bids)

	for i := range bids {
		bids[i].AuctionID = id
	}

	// log.Println(bids)

	err = data.DBclient.InsertBids(&bids)
	if err != nil {
		log.Fatal("Error populating bids.json : ", err.Error())
	}
	return &bids
}

func RemoveBidsFromDB() {
	var bids models.BidList
	file, err := ioutil.ReadFile("./seed-data/bids.json")
	if err != nil {
		log.Fatal("Error reading bids.json : ", err.Error())
	}
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.DeleteBids(&bids)
	if err != nil {
		log.Fatal("Error deleting bids.json : ", err.Error())
	}
}

func RemoveUsersFromDB() {
	err := data.DBclient.DeleteAllUsers()
	if err != nil {
		log.Fatal("Error deleting users: ", err.Error())
	}
}

var _ = Describe("Server", func() {
	var (
		router           *gin.Engine
		response         *httptest.ResponseRecorder
		insertedAuctions *models.AuctionList
	)

	BeforeEach(func() {
		router = CreateRouter()
		dbURL := "mongodb://localhost:27017"
		dbName := "testing"

		ConnectToDB(dbURL, dbName)
		RemoveAuctionsFromDB()
		RemoveBidsFromDB()
		RemoveUsersFromDB()
		insertedAuctions = InsertAuctionsToDB()
		InsertBidsToDB()
		InsertUsersToDB()
	})

	Describe("The /hello endpoint", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/hello", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns the String 'Hello World'", func() {
			Expect("Hello World").To(Equal("Hello World"))
		})
	})

	Describe("Getting all auctions from db", func() {
		var returnedAuctions *models.AuctionList
		BeforeEach(func() {
			returnedAuctions = data.GetAllAuctions()
		})

		It("Does not return nil", func() {
			Expect(returnedAuctions).To(Not(BeNil()))
		})

		It("Returns all 3 auctions", func() {
			Expect(returnedAuctions).To(Equal(insertedAuctions))
		})
	})

	// Describe("The GET auctions/:auction_id/rewards/:user_id endpoint", func() {
	// 	userID := "deepak"
	// 	rewardPercentage := 0.005
	// 	var bids []models.Bid
	// 	var auc *models.Auction
	// 	var user models.User
	// 	var responseAuction models.ResponseCreateAuction
	// 	BeforeEach(func() {
	// 		newAuction := (*GetAuctions())[0]
	// 		databytes, err := json.Marshal(newAuction)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		response = performRequest(router, "POST", "/auction/create", databytes)
	// 		json.Unmarshal(response.Body.Bytes(), &responseAuction)
	// 		newAuction = responseAuction.Auction
	// 		auctionID := newAuction.AuctionID
	// 		bytes, _ := auctionID.MarshalJSON()
	// 		// n := bytes.Index(bytes, []byte{0})
	// 		// id = string(id[:])
	// 		// id := fmt.Sprintf("%s", bytes[1:])
	// 		log.Println("Id: ", id)
	// 		InsertBidsForAuction(id)
	// 		bids = data.GetAllSortedBidsForAuction(id)
	// 		user = data.GetUserByID(userID)
	// 		response = performRequest(router, "GET", "/auctions/"+id+"/rewards/"+userID, nil)
	// 	})
	// 	It("Returns with Status 200", func() {
	// 		Expect(response.Code).To(Equal(200))
	// 	})
	// 	It("Updates User's shashank points in the DB for auction 5dca6431de52283587609581", func() {
	// 		userPoints := make(map[string]int)
	// 		for _, bid := range bids {
	// 			if bid.UserID == userID {
	// 				fmt.Println(user)
	// 				previousPoints := user.Points
	// 				pointsForBidPrice := (rewardPercentage * float64(bid.Price))
	// 				pointsForHighBid := float64(bid.Price-2*auc.BasePrice) / float64(2*auc.BasePrice)
	// 				points := int(pointsForHighBid * pointsForBidPrice)
	// 				if points <= 0 {
	// 					points = 0
	// 				}
	//
	// 				_, ok := userPoints[user.UserID]
	// 				if ok == true {
	// 					userPoints[user.UserID] += points
	// 				} else {
	// 					userPoints[user.UserID] = points + previousPoints
	// 				}
	// 			}
	// 		}
	//
	// 		user = data.GetUserByID(userID)
	// 		Expect(user.Points).To(Equal(userPoints[userID]))
	// 	})
	// })

	Describe("The POST auctions/create endpoint: Successfully created", func() {
		var newAuction models.Auction
		var responseAuction models.ResponseCreateAuction
		newAuction = (*GetAuctions())[0]
		BeforeEach(func() {
			dbURL := "mongodb://localhost:27017"
			dbName := "testing"
			ConnectToDB(dbURL, dbName)
			RemoveAuctionsFromDB()
			data, err := json.Marshal(newAuction)
			if err != nil {
				log.Fatal(err)
			}
			response = performRequest(router, "POST", "/auction/create", data)
		})
		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})
		It("Returns with Message Auction Successfully created.", func() {
			json.Unmarshal(response.Body.Bytes(), &responseAuction)
			Expect(responseAuction.Message).To(Equal("Auction Successfully created."))
		})
		It("Returns auction details", func() {
			responseAuction.Auction.AuctionID = newAuction.AuctionID
			// log.Println(responseAuction)
			// log.Println(responseAuction.Auction.CreatedBy)
			Expect(responseAuction.Auction).To(Equal(newAuction))
		})
	})

	Describe("The POST auctions/create endpoint: Title Error", func() {
		var newAuction models.Auction
		var responseAuction models.ResponseCreateAuction
		newAuction = (*GetAuctions())[0]
		BeforeEach(func() {
			dbURL := "mongodb://localhost:27017"
			dbName := "testing"
			ConnectToDB(dbURL, dbName)
			RemoveAuctionsFromDB()
			newAuction.Title = ""
			data, err := json.Marshal(newAuction)
			if err != nil {
				log.Fatal(err)
			}
			response = performRequest(router, "POST", "/auction/create", data)
		})
		It("Returns with Status 500", func() {
			Expect(response.Code).To(Equal(500))
		})
		It("Returns with Message Invalid Auction Title", func() {
			json.Unmarshal(response.Body.Bytes(), &responseAuction)
			Expect(responseAuction.Message).To(Equal("Invalid Auction Title"))
		})
		It("Returns auction details", func() {
			Expect(responseAuction.Auction).To(Equal(newAuction))
		})
	})

	Describe("The POST auctions/create endpoint: Image Error", func() {
		var newAuction models.Auction
		var responseAuction models.ResponseCreateAuction
		newAuction = (*GetAuctions())[0]
		BeforeEach(func() {
			dbURL := "mongodb://localhost:27017"
			dbName := "testing"
			ConnectToDB(dbURL, dbName)
			RemoveAuctionsFromDB()
			newAuction.Image = []string{}
			data, err := json.Marshal(newAuction)
			if err != nil {
				log.Fatal(err)
			}
			response = performRequest(router, "POST", "/auction/create", data)
		})
		It("Returns with Status 500", func() {
			Expect(response.Code).To(Equal(500))
		})
		It("Returns with Message Please upload auction image", func() {
			json.Unmarshal(response.Body.Bytes(), &responseAuction)
			Expect(responseAuction.Message).To(Equal("Please upload auction image"))
		})
		It("Returns auction details", func() {
			Expect(responseAuction.Auction).To(Equal(newAuction))
		})
	})

	Describe("The POST /user/signup endpoint", func() {
		var newUser []byte
		var responseSignup models.ResponseSignup
		BeforeEach(func() {
			newUser = []byte(`{"user_id":"jd", "user_name":"john_doe", "pwd":"pwd_john_doe"}`)
			response = performRequest(router, "POST", "/user/signup", newUser)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns with Message User signup successful", func() {
			json.Unmarshal(response.Body.Bytes(), &responseSignup)
			Expect(responseSignup.Message).To(Equal("User signup successful"))
		})

		It("adds new user jd", func() {
			var user models.User
			json.Unmarshal(newUser, &user)
			log.Println("User: ", user)
			log.Println("Response: ", responseSignup)
			Expect(responseSignup.User).To(Equal(user))
		})
	})

	Describe("The POST /user/signup endpoint: User already exist", func() {
		var newUser []byte
		var responseSignup models.ResponseSignup
		BeforeEach(func() {
			newUser = []byte(`{"user_id":"jd", "user_name":"john_doe", "pwd":"pwd_john_doe"}`)
			response = performRequest(router, "POST", "/user/signup", newUser)
			response = performRequest(router, "POST", "/user/signup", newUser)
		})

		It("Returns with Status 500", func() {
			Expect(response.Code).To(Equal(500))
		})

		It("Returns with Message User already exists", func() {
			json.Unmarshal(response.Body.Bytes(), &responseSignup)
			Expect(responseSignup.Message).To(Equal("User already exists"))
		})
	})

	Describe("The POST /user/signup endpoint: UserID is empty", func() {
		var newUser []byte
		var responseSignup models.ResponseSignup
		BeforeEach(func() {
			newUser = []byte(`{"user_id":"", "user_name":"john_doe", "pwd":"pwd_john_doe"}`)
			response = performRequest(router, "POST", "/user/signup", newUser)
		})

		It("Returns with Status 500", func() {
			Expect(response.Code).To(Equal(500))
		})

		It("Returns with Message UserID is empty", func() {
			json.Unmarshal(response.Body.Bytes(), &responseSignup)
			Expect(responseSignup.Message).To(Equal("UserID is empty"))
		})
	})

	Describe("The POST /user/signup endpoint: Password is empty", func() {
		var newUser []byte
		var responseSignup models.ResponseSignup
		BeforeEach(func() {
			newUser = []byte(`{"user_id":"jd", "user_name":"john_doe", "pwd":""}`)
			response = performRequest(router, "POST", "/user/signup", newUser)
		})

		It("Returns with Status 500", func() {
			Expect(response.Code).To(Equal(500))
		})

		It("Returns with Message UserID is empty", func() {
			json.Unmarshal(response.Body.Bytes(), &responseSignup)
			Expect(responseSignup.Message).To(Equal("Password is empty"))
		})
	})

	Describe("The /user/:user_id endpoint", func() {
		var new_user []byte
		var responseSignup models.ResponseSignup
		var newUser, responseUser models.User
		BeforeEach(func() {
			new_user = []byte(`{"user_id":"jd", "user_name":"john_doe", "pwd":"pwd_john_doe"}`)
			response = performRequest(router, "POST", "/user/signup", new_user)
			json.Unmarshal(response.Body.Bytes(), &responseSignup)
			newUser = responseSignup.User
			response = performRequest(router, "GET", "/user/jd", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns the user details with input id", func() {
			json.Unmarshal(response.Body.Bytes(), &responseUser)
			Expect(responseUser).To(Equal(newUser))
		})
	})

	Describe("The GET /auction/create/getTagsfromImage", func() {
		BeforeEach(func() {

			response = performRequest(router, "GET", "/auction/create/getTagsfromImage?imageName=test1.png", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})
	})

	Describe("The GET /auction/create/getDescriptionfromImage", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/auction/create/getDescriptionfromImage?imageName=test1.png", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})
	})

	// Describe("The POST /auctions/:auction_id/users/:user_id/bids endpoint", func() {
	// 	BeforeEach(func() {
	// 		newbid := []byte(`{"BidID":"0", "AuctionID":"1", "UserID":"vamshi", "Price":"2000", "Count": 11, "Time":""}`)
	// 		response = performRequest(router, "POST", "/auctions/0/users/vamshi/bids", newbid)
	// 	})

	// 	It("Returns with Status 200", func() {
	// 		Expect(response.Code).To(Equal(200))
	// 	})

	// 	It("Adds bid by user vamshi to auction 0", func() {
	// 		// newbid := models.Bid{"0", "1", "vamshi", models.Price(2000), 121232323}
	// 		// var bids *[]models.Bid
	// 		// bids, _ = data.DBclient.GetBids("1")
	// 		// bids = append(*bids, *newbid)
	// 		Expect(true).To(Equal(true))
	// 	})
	// })

	// Describe("The GET /auctions/:auction_id/result", func() {
	// 	BeforeEach(func() {
	// 		response = performRequest(router, "GET", "/auctions/1/result", nil)
	// 	})

	// 	It("Returns with Status 200", func() {
	// 		Expect(response.Code).To(Equal(200))
	// 	})

	// 	It("Returns result of auction 1", func() {
	// 		var res models.Result
	// 		json.Unmarshal(response.Body.Bytes(), &res)
	// 		// actual := data.DBclient.GetResult("1")
	// 		// Expect(res).To(Equal(actual))
	// 	})
})
