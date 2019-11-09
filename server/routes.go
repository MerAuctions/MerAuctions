package server

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine) {

	router.GET("/hello", hello)

	router.GET("/", getAllAuctions)
	router.GET("/auctions/:id", getAuctionsById)
	router.GET("/auctions/:id/bids", getBidsAuctionsById)
	router.POST("/users", addNewUser)
	router.POST("/auctions/:id/users/:id/bids", addBidAuctionIdByUserId)
	router.GET("/auctions/:id/result", getResultByAuctionId)

}
