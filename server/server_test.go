package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"

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

func InsertAuctionsToDB() {
	var auc []models.Auction
	file, err := ioutil.ReadFile("./seed-data/auctions.json")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.InsertAuctions(&auc)
	if err != nil {
		return
	}
}

func RemoveAuctionsFromDB() {
	var auc []models.Auction
	file, err := ioutil.ReadFile("./seed-data/auctions.json")
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal([]byte(file), &auc)
	err = data.DBclient.DeleteAuctions(&auc)
	if err != nil {
		return
	}
}

func InsertBidsToDB() {
	var bids []models.Bid
	file, err := ioutil.ReadFile("./seed-data/bids.json")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.InsertBids(&bids)
	if err != nil {
		return
	}
}

func RemoveBidsFromDB() {
	var bids []models.Bid
	file, err := ioutil.ReadFile("./seed-data/bids.json")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(file))
	json.Unmarshal([]byte(file), &bids)
	err = data.DBclient.DeleteBids(&bids)
	if err != nil {
		return
	}
}

var _ = Describe("Server", func() {
	var (
		router   *gin.Engine
		response *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		router = CreateRouter()
		dbURL := "mongodb://localhost:27017"
		dbName := "testing5"

		ConnectToDB(dbURL, dbName)
		InsertAuctionsToDB()
	})

	AfterEach(func() {
		RemoveAuctionsFromDB()
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

	Describe("The GET / endpoint", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/", nil)

		})

		// AfterEach(func() {
		// 	RemoveAuctionsFromDB()
		// })

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns all auctions", func() {
			var allauc models.AuctionList
			json.Unmarshal(response.Body.Bytes(), &allauc)
			actualauc, _ := json.Marshal(data.DBclient.GetAuctions())
			Expect(response.Body.Bytes()).To(Equal(actualauc))
		})
	})

	Describe("The GET auctions/:auction_id endpoint", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/auctions/1", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns auction 1", func() {
			var auc models.Auction
			json.Unmarshal(response.Body.Bytes(), &auc)
			actualauc, _ := data.DBclient.GetAuction("1")
			actual, _ := json.Marshal(actualauc)
			Expect(response.Body.Bytes()).To(Equal(actual))
		})
	})

	Describe("The GET auctions/:auction_id/bids endpoint", func() {

		BeforeEach(func() {
			InsertBidsToDB()
			response = performRequest(router, "GET", "/auctions/1/bids", nil)

			// newbid1 := models.Bid{"1", "1", "bowei", 2000000, 1573516800}
			// data.DBclient.InsertBid(&newbid1)
			// newbid2 := models.Bid{"2", "1", "vamshi", 2500000, 1573516810}
			// data.DBclient.InsertBid(&newbid2)
			// newbid3 := models.Bid{"3", "1", "harsh", 3000000, 1573516820}
			// data.DBclient.InsertBid(&newbid3)
			// newbid4 := models.Bid{"4", "1", "deepak", 4000000, 1573516830}
			// data.DBclient.InsertBid(&newbid4)
			// newbid5 := models.Bid{"5", "1", "vamshi", 3500000, 1573516840}
			// data.DBclient.InsertBid(&newbid5)
		})

		AfterEach(func() {
			RemoveBidsFromDB()
		})

		// It("Returns with Status 200", func() {
		// 	Expect(response.Code).To(Equal(200))
		// })

		It("Returns all running auctions", func() {
			var allbids []models.Bid
			json.Unmarshal(response.Body.Bytes(), &allbids)
			act, _ := data.DBclient.GetBids("1")
			var actual []models.Bid
			actual = *act
			sort.Slice(actual, func(i, j int) bool {
				return actual[i].Price < actual[j].Price
			})
			sort.Slice(allbids, func(i, j int) bool {
				return allbids[i].Price < allbids[j].Price
			})
			//actualbytes, _ := json.Marshal(&actual)
			Expect(allbids).To(Equal(actual))
		})
	})

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

	Describe("The GET /auctions/:auction_id/result", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/auctions/1/result", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns result of auction 1", func() {
			var res models.Result
			json.Unmarshal(response.Body.Bytes(), &res)
			// actual := data.DBclient.GetResult("1")
			// Expect(res).To(Equal(actual))
		})
	})
})
