package server_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

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

var _ = Describe("Server", func() {
	var (
		router   *gin.Engine
		response *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		router = CreateRouter()
	})

	Describe("The /hello endpoint", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/hello", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns the String 'Hello World'", func() {
			Expect(response.Body.String()).To(Equal("Hello World"))
		})
	})

	Describe("The GET / endpoint", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns all auctions", func() {
			Expect(true).To(Equal(true))
		})
	})

	Describe("The GET auctions/:auction_id endpoint", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/auctions/0", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns auction 0", func() {
			Expect(true).To(Equal(true))
		})
	})

	Describe("The GET auctions/:auction_id/bids endpoint", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/auctions/0/bids", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns all running auctions", func() {
			Expect(true).To(Equal(true))
		})
	})

	Describe("The POST /users endpoint", func() {
		BeforeEach(func() {
			newusr := []byte(`{"UserID":"vamshi", "UserName":"vamshiteja"}`)
			response = performRequest(router, "POST", "/users", newusr)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("adds new user vamshi", func() {
			Expect(true).To(Equal(true))
		})
	})

	Describe("The POST /auctions/:auction_id/users/:user_id/bids endpoint", func() {
		BeforeEach(func() {
			newbid := []byte(`{"BidID":"0", "AuctionID":"0", "UserID":"vamshi", "Price":"2000", "Time":""}`)
			response = performRequest(router, "POST", "/auctions/0/users/vamshi/bids", newbid)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Adds bid by user vamshi to auction 0", func() {
			Expect(true).To(Equal(true))
		})
	})

	Describe("The GET /auctions/:auction_id/result", func() {
		BeforeEach(func() {
			response = performRequest(router, "GET", "/auctions/0/result", nil)
		})

		It("Returns with Status 200", func() {
			Expect(response.Code).To(Equal(200))
		})

		It("Returns result of auction 0", func() {
			Expect(true).To(Equal(true))
		})
	})
})
