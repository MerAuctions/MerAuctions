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

// func GetAuctions() *models.AuctionList {
//
// }

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
		insertedAuctions = InsertAuctionsToDB()
		InsertBidsToDB()
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

	Describe("The GET auctions/:auction_id/bids endpoint", func() {

		auctionID := "5dca6431de52283587609581"
		BeforeEach(func() {
			InsertBidsToDB()
			response = performRequest(router, "GET", "/auctions/"+auctionID+"/bids", nil)
		})
		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})
		It("Returns top 5 bids of running auction 5dca6431de52283587609581", func() {
			var receivedBidsList *[]models.Bid
			json.Unmarshal(response.Body.Bytes(), &receivedBidsList)
			returnedBids := data.GetTopFiveBids(auctionID)
			Expect(receivedBidsList).To(Equal(returnedBids))
		})
	})

	// Describe("The POST auctions/create endpoint", func() {
	// 	var newAuction models.Auction
	// 	getAuction
	// 	BeforeEach(func() {
	// 		var data []byte
	// 		response = performRequest(router, "POST", "/auctions/create", data)
	// 	})
	// 	It("Returns with Status 200", func() {
	// 		Expect(response.Code).To(Equal(200))
	// 	})
	// 	It("Returns top 5 bids of running auction 5dca6431de52283587609581", func() {
	// 		var auction models.Auction
	// 		json.Unmarshal(response.Body.Bytes(), &auction)
	// 		Expect(auction).To(Equal(newAuction))
	// 	})
	// })

	// Describe("The POST /users endpoint", func() {
	// 	BeforeEach(func() {
	// 		newusr := []byte(`{"UserID":"vamshi", "UserName":"vamshiteja"}`)
	// 		response = performRequest(router, "POST", "/users", newusr)
	// 	})

	// 	It("Returns with Status 200", func() {
	// 		Expect(response.Code).To(Equal(200))
	// 	})

	// 	It("adds new user vamshi", func() {
	// 		newusr := []byte(`{"UserID":"vamshi", "UserName":"vamshiteja"}`)
	// 		usr, _ := data.DBclient.Getuser("vamshi")
	// 		usrbyte, _ := json.Marshal(usr)
	// 		Expect(usrbyte).To(Equal(newusr))
	// 	})
	// })

	// Describe("The POST /auctions/:auction_id/users/:user_id/bids endpoint", func() {
	// 	BeforeEach(func() {
	// 		newbid := []byte(`{"BidID":"0", "AuctionID":"1", "UserID":"vamshi", "Price":"2000", "Time":""}`)
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
